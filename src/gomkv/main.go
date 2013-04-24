package main

import (
	"exec"
	"flag"
	"fmt"
	"gomkv/config"
	"gomkv/handbrake"
	"os"
	"os/user"
	"path/filepath"
	toml "github.com/stvp/go-toml-config"
)

var files []string

var defaults = config.GomkvConfig{}
var debug = false

const (
	DEBUG_LEVEL_BASIC = 1
	DEBUG_LEVEL_RAGEL = 2
	DEBUG_LEVEL_EXEC  = 3
)

func init() {
	var err error
	var debuglvl = 0
	mobile := false

	if myuser, err := user.Current(); err != nil {
		panic(err)
	} else {
		gomkvpath := myuser.HomeDir + "/.gomkvrc"
		if _, err := os.Open(gomkvpath); err == nil {
			if err := toml.Parse(myuser.HomeDir + "/.gomkvrc"); err != nil {
				panic(err)
			}
			defaults.DestDir = *toml.String("dest-dir", "")
			defaults.Languages = *toml.String("languages", "")
			defaults.Profile = *toml.String("profile", "")
		}
	}

	flag.StringVar(&defaults.Profile, "profile", defaults.Profile, "Default Encoding Profile. Defaults to 'High Profile'")
	flag.StringVar(&defaults.Prefix, "prefix", config.DEFAULT_PREFIX, "Default Prefix for output filename(s)")
	flag.BoolVar(&defaults.Episodic, "series", false, "Videos are episodes of a series")
	flag.IntVar(&defaults.EpisodeOffset, "episode", 1, "Episode starting offset.")
	flag.IntVar(&defaults.SeasonOffset, "season", 1, "Season starting offset.")
	flag.BoolVar(&defaults.AacOnly, "aac", false, "Encode audio using aac, instead of copying")
	flag.BoolVar(&mobile, "mobile", false, "Use mobile friendly settings")
	flag.BoolVar(&defaults.EnableSubs, "subs", true, "Copy subtitles")
	flag.IntVar(&debuglvl, "debug", 0, "Debug level 1..3")
	flag.StringVar(&defaults.SrcDir, "source-dir", "", "directory containing video files. Defaults to current working directory.")
	flag.StringVar(&defaults.DestDir, "dest-dir", defaults.DestDir, "directory you want video files to be created")
	flag.StringVar(&defaults.Languages, "languages", defaults.Languages, "list of languages and order to copy, comma separated e.g. English,Japanese")
	flag.StringVar(&defaults.DefaultSub, "subtitle-default", "", "Enable subtitles by default for the language matching this value. e.g. -subtitle-default=English")
	flag.Parse()

	workingDir := ""
	if defaults.SrcDir == "" {
		workingDir, err = filepath.Abs(".")
		if err != nil {
			panic(err)
		}
	} else {
		workingDir, err = validateDirectory(defaults.SrcDir)
		if err != nil {
			panic(err)
		}
	}
	defaults.SrcDir = workingDir

	if defaults.DestDir, err = validateDirectory(defaults.DestDir); err != nil {
		panic(err)
	}
	if debug {
		fmt.Printf("srcdir: %s destdir: %s", defaults.SrcDir, defaults.DestDir)
	}

	switch {
	case debuglvl == DEBUG_LEVEL_BASIC:
		debug = true
	case debuglvl == DEBUG_LEVEL_RAGEL:
		debug = true
		handbrake.DebugEnabled = true
	case debuglvl == DEBUG_LEVEL_EXEC:
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

func processOne(file string) {
	std, err := exec.Command("HandBrakeCLI", "-t0", "-i", file)
	if err != nil {
		panic(err)
	}
	meta := handbrake.ParseOutput(std.Err)
	if debug {
		fmt.Fprintln(os.Stderr, meta)
	}
	result, err := handbrake.FormatCLIOutput(meta, &defaults)
	fmt.Println(result)
}

func main() {
	if err := os.Chdir(defaults.SrcDir); err != nil {
		panic(err)
	}

	if debug {
		fmt.Fprintln(os.Stderr, "Working Directory: "+defaults.SrcDir)
	}
	mkv, err := filepath.Glob(defaults.SrcDir + "/*.mkv")
	if err != nil {
		panic(err)
	}
	m4v, err := filepath.Glob(defaults.SrcDir + "/*.m4v")
	if err != nil {
		panic(err)
	}
	files = append(mkv, m4v...)
	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No mkv/m4v files found in %s. Exiting.\n", defaults.SrcDir)
		os.Exit(1)
	}
	for _, file := range files {
		processOne(file)
	}
}
