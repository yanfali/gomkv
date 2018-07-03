package config

import (
	"strings"
)

// exports
const (
	DefaultProfile  = "High Profile"
	DefaultPrefix   = ""
	DefaultFileGlob = "m[k4]v"
)

// GomkvSession struct
type GomkvSession struct {
	Episode int
	Chapter int
}

// GomkvConfig configuration of program
type GomkvConfig struct {
	Profile        string
	Prefix         string
	EpisodeOffset  int
	SeasonOffset   int
	Episodic       bool
	AacOnly        bool
	M4v            bool
	EnableSubs     bool
	SrcDir         string
	DestDir        string
	Languages      string
	DefaultSub     string
	SplitFileEvery int
	DisableAAC     bool
	Goroutines     int
	FileGlob       string
}

// LanguageOrderMap describes the preference order they should be encoded in
func (g *GomkvConfig) LanguageOrderMap() map[string]int {
	langOrder := map[string]int{}
	for i, language := range strings.Split(g.Languages, ",") {
		langOrder[language] = i
	}
	return langOrder
}

// Mobile configures handbrake to be mobile friendly
func (g *GomkvConfig) Mobile() *GomkvConfig {
	g.Profile = "Universal"
	g.AacOnly = true
	g.M4v = true
	g.EnableSubs = false
	return g
}
