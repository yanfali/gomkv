package handbrake

import (
	"math"
)

// ChapterMeta data
type ChapterMeta struct {
	Index    int
	Duration string
}

// AudioMeta data
type AudioMeta struct {
	Language  string
	Codec     string
	Channels  string
	Frequency int
	Bps       int
	Index     int
}

// AudioMetas is collection of Audio metadata
type AudioMetas []*AudioMeta

// Len Sort function support
func (a AudioMetas) Len() int { return len(a) }

// Swap Sort function support
func (a AudioMetas) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// ByLanguage Audio metadata ordered by language preference
type ByLanguage struct {
	AudioMetas
	Order map[string]int
}

// GetOrder of language based on the Language key
func (a ByLanguage) GetOrder(key string) int {
	val, ok := a.Order[key]
	if ok {
		return val
	}
	return int(math.MaxInt32) // TODO depends on Go Version and platform
}

// Less Sort support
func (a ByLanguage) Less(i, j int) bool {
	lhs := a.AudioMetas[i].Language
	rhs := a.AudioMetas[j].Language
	return a.GetOrder(lhs) < a.GetOrder(rhs)
}

// SubtitleMeta data
type SubtitleMeta struct {
	Language string
	Type     string
	Format   string
}

// Meta data for configuration of HandbrakeCLI
type Meta struct {
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
