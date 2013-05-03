package config

import (
	"strings"
)

const (
	DEFAULT_PROFILE = "High Profile"
	DEFAULT_PREFIX  = ""
)

type GomkvSession struct {
	Episode int
	Chapter int
}

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
}

func (g *GomkvConfig) LanguageOrderMap() map[string]int {
	langOrder := map[string]int{}
	for i, language := range strings.Split(g.Languages, ",") {
		langOrder[language] = i
	}
	return langOrder
}

func (g *GomkvConfig) Mobile() *GomkvConfig {
	g.Profile = "Universal"
	g.AacOnly = true
	g.M4v = true
	g.EnableSubs = false
	return g
}
