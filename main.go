package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/yanfali/gomkv/config"
	"github.com/yanfali/gomkv/exec"
	"github.com/yanfali/gomkv/handbrake"
)

var (
	files    []string
	defaults = config.GomkvConfig{}
	debuglvl int
)

const (
	debugLevelBasic = iota
	debugLevelRagel = iota
	debugLevelExec  = iota
)

func init() {
	log.SetPrefix("[gomkv] ")
	var (
		err error
	)
	debuglvl = 0
	mobile := false
	flag.StringVar(&defaults.Profile, "profile", config.DefaultProfile, "Encoding Profile. Defaults to 'High Profile'")
	flag.StringVar(&defaults.Prefix, "prefix", config.DefaultPrefix, "Prefix for output filename(s)")
	flag.BoolVar(&defaults.Episodic, "series", false, "Videos are episodes in a series")
	flag.IntVar(&defaults.EpisodeOffset, "episode", 1, "Episode starting offset.")
	flag.IntVar(&defaults.SeasonOffset, "season", 1, "Season starting offset.")
	flag.BoolVar(&defaults.AacOnly, "aac", false, "Encode audio using aac, instead of copying verbatim")
	flag.BoolVar(&mobile, "mobile", false, "Use mobile friendly settings")
	flag.BoolVar(&defaults.EnableSubs, "subs", true, "Copy subtitles")
	flag.IntVar(&debuglvl, "debug", 0, "Debug level 1..3. 1) basic, 2) parser 3) execution of commands")
	flag.StringVar(&defaults.SrcDir, "source-dir", "", "directory containing video files. Defaults to current working directory.")
	flag.StringVar(&defaults.DestDir, "dest-dir", "", "directory you want new video files to be created")
	flag.StringVar(&defaults.Languages, "languages", "", "list of languages and order to copy, comma separated e.g. English,Japanese")
	flag.StringVar(&defaults.DefaultSub, "subtitle-default", "", "Enable subtitles by default for the language matching this value. e.g. -subtitle-default=English")
	flag.IntVar(&defaults.SplitFileEvery, "split-chapters", 0, "Create one file for every N chapters. Only works with --series. e.g. -split-chapters 5")
	flag.BoolVar(&defaults.DisableAAC, "disable-aac", false, "Disable Automatic AAC Audio Generation For Non-Mobile")
	flag.IntVar(&defaults.Goroutines, "goroutines", 2, "Max number of go routines to use for parsing. Controls instances of HandbrakeCLI used to parse video files")
	flag.StringVar(&defaults.FileGlob, "fileglob", config.DefaultFileGlob, "default fileglob pattern")
	flag.Parse()

	var workingDir string
	if defaults.SrcDir == "" {
		workingDir, err = filepath.Abs(".")
		if err != nil {
			log.Fatalf("Error resolving cwd directory: %v\n", err)
		}
	} else {
		workingDir, err = validateDirectory(defaults.SrcDir)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	}
	defaults.SrcDir = workingDir

	if defaults.DestDir, err = validateDirectory(defaults.DestDir); err != nil {
		log.Fatalf("%v\n", err)
	}
	if debuglvl > 0 {
		log.Printf("srcdir: %q destdir: %q", defaults.SrcDir, defaults.DestDir)
	}

	switch debuglvl {
	case debugLevelBasic:
	case debugLevelRagel:
		handbrake.DebugEnabled = true
	case debugLevelExec:
		handbrake.DebugEnabled = true
		exec.Debug = true
	}
	if debuglvl > 0 {
		fmt.Fprintf(os.Stderr, "Enabling debug level %d\n", debuglvl)
	}
	if mobile {
		defaults.Mobile()
	}
	if defaults.SplitFileEvery > 0 {
		fmt.Fprintf(os.Stderr, "%sDisabling Goroutines Because of Chapter Splitting\n", log.Prefix())
		defaults.Goroutines = 1
	}
}

func processOne(session *config.GomkvSession, file string) error {
	std, err := exec.Command("HandBrakeCLI", "-t0", "-i", file)
	if err != nil {
		log.Println(err)
		return err
	}

	// This is the amdahl blocker. For unrelated video files this is
	// embarrassingly parallel. For cases where each file is an episode
	// this is parallelizable. For cases where each set of chapters
	// is part of a set of chapters we must wait until we have processed
	// the meta data.
	meta := handbrake.ParseOutput(std.Err)
	if debuglvl > 0 {
		log.Printf("+%v", meta)
	}
	if defaults.Episodic && defaults.SplitFileEvery > 0 {
		session.Chapter = 1
	}

	// copy session object rather than pass through so it can increment
	// episode number independently of original
	thisSession := *session
	results, err := handbrake.FormatCLIOutput(meta, &defaults, &thisSession)
	if err != nil {
		return err
	}

	// incremement the episode
	if defaults.Episodic {
		// handle chapters within a file
		if defaults.SplitFileEvery > 0 {
			session.Episode += int(math.Max(float64(len(meta.Chapter)/defaults.SplitFileEvery), 1))
		}
	}
	for _, result := range results {
		fmt.Println(result)
	}
	return nil
}

func main() {
	if err := os.Chdir(defaults.SrcDir); err != nil {
		log.Fatalln(err)
	}

	if debuglvl > 0 {
		log.Printf("Working Directory: %q\n", defaults.SrcDir)
	}
	files, err := filepath.Glob(defaults.SrcDir + "/*[.]" + defaults.FileGlob)
	if err != nil {
		log.Fatalln(err)
	}
	if len(files) == 0 {
		log.Printf("No %s files found in %q. Exiting.\n", defaults.FileGlob, defaults.SrcDir)
		os.Exit(1)
	}
	session := &config.GomkvSession{Episode: defaults.EpisodeOffset}

	var wg sync.WaitGroup
	wg.Add(len(files)) // wait for output from all files
	// limit the number of goroutines that can run concurrently
	semaphore := make(chan int, defaults.Goroutines)
	for _, file := range files {
		go func(file string, session config.GomkvSession) {
			semaphore <- 1 // mark a goroutine in use
			if err := processOne(&session, file); err != nil {
				log.Printf("%q %v", file, err)
			}
			wg.Done()   // notify wait group a file has been processed
			<-semaphore // go routine returned
		}(file, *session)
		if defaults.Episodic && defaults.SplitFileEvery == 0 {
			session.Episode++
		}
	}
	wg.Wait() // wait for all files to have been processed
}
