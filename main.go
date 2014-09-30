package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/yanfali/gomkv/config"
	"github.com/yanfali/gomkv/exec"
	"github.com/yanfali/gomkv/handbrake"
)

var (
	files    []string
	defaults = config.GomkvConfig{}
	debug    = false
)

const (
	DEBUG_LEVEL_BASIC = 1
	DEBUG_LEVEL_RAGEL = 2
	DEBUG_LEVEL_EXEC  = 3
)

func init() {
	log.SetPrefix("[gomkv] ")
	var (
		err      error
		debuglvl = 0
	)
	mobile := false
	flag.StringVar(&defaults.Profile, "profile", config.DEFAULT_PROFILE, "Default Encoding Profile. Defaults to 'High Profile'")
	flag.StringVar(&defaults.Prefix, "prefix", config.DEFAULT_PREFIX, "Default Prefix for output filename(s)")
	flag.BoolVar(&defaults.Episodic, "series", false, "Videos are episodes of a series")
	flag.IntVar(&defaults.EpisodeOffset, "episode", 1, "Episode starting offset.")
	flag.IntVar(&defaults.SeasonOffset, "season", 1, "Season starting offset.")
	flag.BoolVar(&defaults.AacOnly, "aac", false, "Encode audio using aac, instead of copying")
	flag.BoolVar(&mobile, "mobile", false, "Use mobile friendly settings")
	flag.BoolVar(&defaults.EnableSubs, "subs", true, "Copy subtitles")
	flag.IntVar(&debuglvl, "debug", 0, "Debug level 1..3")
	flag.StringVar(&defaults.SrcDir, "source-dir", "", "directory containing video files. Defaults to current working directory.")
	flag.StringVar(&defaults.DestDir, "dest-dir", "", "directory you want video files to be created")
	flag.StringVar(&defaults.Languages, "languages", "", "list of languages and order to copy, comma separated e.g. English,Japanese")
	flag.StringVar(&defaults.DefaultSub, "subtitle-default", "", "Enable subtitles by default for the language matching this value. e.g. -subtitle-default=English")
	flag.IntVar(&defaults.SplitFileEvery, "split-chapters", 0, "Create one file for every N chapters. Only works with --series. e.g. -split-chapters 5")
	flag.BoolVar(&defaults.DisableAAC, "disable-aac", false, "Disable Automatic AAC Audio Generation For Non-Mobile")
	flag.Parse()

	workingDir := ""
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
	if debug {
		log.Printf("srcdir: %q destdir: %q", defaults.SrcDir, defaults.DestDir)
	}

	switch debuglvl {
	case DEBUG_LEVEL_BASIC:
		debug = true
	case DEBUG_LEVEL_RAGEL:
		debug = true
		handbrake.DebugEnabled = true
	case DEBUG_LEVEL_EXEC:
		debug = true
		handbrake.DebugEnabled = true
		exec.Debug = true
	}
	if debuglvl > 0 {
		fmt.Fprintf(os.Stderr, "Enabling debug level %d\n", debuglvl)
	}
	if mobile {
		defaults.Mobile()
	}
}

func processOne(session *config.GomkvSession, file string) error {
	std, err := exec.Command("HandBrakeCLI", "-t0", "-i", file)
	if err != nil {
		log.Println(err)
		return err
	}
	meta := handbrake.ParseOutput(std.Err)
	if debug {
		log.Println(os.Stderr, meta)
	}
	if defaults.Episodic && defaults.SplitFileEvery > 0 {
		session.Chapter = 1
	}

	results, err := handbrake.FormatCLIOutput(meta, &defaults, session)
	if err != nil {
		return err
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

	if debug {
		log.Printf("Working Directory: %q\n" + defaults.SrcDir)
	}
	mkv, err := filepath.Glob(defaults.SrcDir + "/*.mkv")
	if err != nil {
		log.Fatalln(err)
	}
	m4v, err := filepath.Glob(defaults.SrcDir + "/*.m4v")
	if err != nil {
		log.Fatalln(err)
	}
	files = append(mkv, m4v...)
	if len(files) == 0 {
		log.Printf("No mkv/m4v files found in %q. Exiting.\n", defaults.SrcDir)
		os.Exit(1)
	}
	session := &config.GomkvSession{Episode: defaults.EpisodeOffset}
	for _, file := range files {
		if err := processOne(session, file); err != nil {
			log.Printf("%q %v", file, err)
		}
	}
}
