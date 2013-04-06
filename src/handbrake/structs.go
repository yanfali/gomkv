package handbrake

import ()

type AudioMeta struct {
	Language  string
	Codec     string
	Channels  float32
	Frequency int
	Bps       int
}

type SubtitleMeta struct {
	Language string
	Type     string
}

type HandBrakeMeta struct {
	Title       string
	Duration    float64
	Height      int
	Width       int
	Pixelaspect string
	Aspect      float32
	Fps         float32
	Autocrop    string
	Audio       []AudioMeta
	Subtitle    []SubtitleMeta
}
