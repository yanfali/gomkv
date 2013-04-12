package handbrake

import ()

type AudioMeta struct {
	Language  string
	Codec     string
	Channels  string
	Frequency int
	Bps       int
}

type SubtitleMeta struct {
	Language string
	Type     string
	Format   string
}

type HandBrakeMeta struct {
	Title       string
	Duration    float64
	Height      int
	Width       int
	Pixelaspect string
	Aspect      string
	Fps         string
	Autocrop    string
	Audio       []AudioMeta
	Subtitle    []SubtitleMeta
}
