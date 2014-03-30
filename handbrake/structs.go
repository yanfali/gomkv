package handbrake

import (
	"math"
)

type ChapterMeta struct {
	Index    int
	Duration string
}

type AudioMeta struct {
	Language  string
	Codec     string
	Channels  string
	Frequency int
	Bps       int
	Index     int
}

type AudioMetas []*AudioMeta

func (a AudioMetas) Len() int      { return len(a) }
func (a AudioMetas) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type ByLanguage struct {
	AudioMetas
	Order map[string]int
}

func (a ByLanguage) GetOrder(key string) int {
	val, ok := a.Order[key]
	if ok {
		return val
	}
	return int(math.MaxInt32) // TODO depends on Go Version and platform
}

func (a ByLanguage) Less(i, j int) bool {
	lhs := a.AudioMetas[i].Language
	rhs := a.AudioMetas[j].Language
	return a.GetOrder(lhs) < a.GetOrder(rhs)
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
	Audio       AudioMetas
	Subtitle    []SubtitleMeta
	Chapter     []ChapterMeta
}
