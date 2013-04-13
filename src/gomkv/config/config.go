package config

const (
	DEFAULT_PROFILE = "High Profile"
	DEFAULT_PREFIX  = ""
)

type GomkvConfig struct {
	Profile       string
	Prefix        string
	EpisodeOffset int
	SeasonOffset  int
	Episodic      bool
	AacOnly       bool
	M4v           bool
	EnableSubs    bool
}

func (g *GomkvConfig) Mobile() *GomkvConfig {
	g.Profile = "Universal"
	g.AacOnly = true
	g.M4v = true
	g.EnableSubs = false
	return g
}
