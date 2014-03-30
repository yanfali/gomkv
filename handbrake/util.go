package handbrake

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Section int

var (
	DebugEnabled bool = false
)

const (
	NONE Section = iota
	CHAPTER
	AUDIO
	SUBTITLE
)

func parseTime(timestring string) float64 {
	rawTime := strings.Trim(timestring, " \n")
	splitTime := strings.Split(rawTime, ":")
	var length float64
	hours, err := strconv.ParseInt(splitTime[0], 10, 8)
	if err != nil {
		panic(err)
	}
	length += float64(hours * 60 * 60)
	minutes, err := strconv.ParseInt(splitTime[1], 10, 8)
	length += float64(minutes * 60)
	seconds, err := strconv.ParseInt(splitTime[2], 10, 8)
	length += float64(seconds)
	return length
}

func parseInt(value string) int {
	if result, err := strconv.ParseInt(value, 10, 32); err != nil {
		panic(err)
	} else {
		return int(result)
	}
	return 0
}

func getLastAudioMeta(meta *HandBrakeMeta) *AudioMeta {
	if len(meta.Audio) == 0 {
		panic("No audio available!")
	}
	return meta.Audio[len(meta.Audio)-1]
}

func getLastSubtitleMeta(meta *HandBrakeMeta) *SubtitleMeta {
	if len(meta.Subtitle) == 0 {
		panic("No subtitle available!")
	}
	return &meta.Subtitle[len(meta.Subtitle)-1]
}

func getLastChapterMeta(meta *HandBrakeMeta) *ChapterMeta {
	if len(meta.Chapter) == 0 {
		panic("No subtitle available!")
	}
	return &meta.Chapter[len(meta.Chapter)-1]
}

func debug(format string, args ...interface{}) {
	if DebugEnabled {
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func addAudioMeta(meta *HandBrakeMeta) {
	meta.Audio = append(meta.Audio, &AudioMeta{})
}

func addSubtitleMeta(meta *HandBrakeMeta) {
	subtitle := SubtitleMeta{}
	meta.Subtitle = append(meta.Subtitle, subtitle)
}

func addChapterMeta(meta *HandBrakeMeta) {
	meta.Chapter = append(meta.Chapter, ChapterMeta{})
}
