package handbrake

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type section int

// exports
var (
	DebugEnabled = false
)

// constants
const (
	None section = iota
	Chapter
	Audio
	Subtitle
	minuteInSeconds = 60
	hourInMinutes   = 60
)

func parseDurationIntoSeconds(timestring string) float64 {
	rawTime := strings.Trim(timestring, " \n")
	splitTime := strings.Split(rawTime, ":")
	var totalSeconds float64
	hours, err := strconv.ParseInt(splitTime[0], 10, 8)
	if err != nil {
		panic(err)
	}
	totalSeconds += float64(hours * minuteInSeconds * hourInMinutes)
	minutes, err := strconv.ParseInt(splitTime[1], 10, 8)
	totalSeconds += float64(minutes * minuteInSeconds)
	seconds, err := strconv.ParseInt(splitTime[2], 10, 8)
	totalSeconds += float64(seconds)
	return totalSeconds
}

func parseInt(value string) int {
	if result, err := strconv.ParseInt(value, 10, 32); err != nil {
		panic(err)
	} else {
		return int(result)
	}
}

func getLastAudioMeta(meta *Meta) *AudioMeta {
	if len(meta.Audio) == 0 {
		panic("No audio available!")
	}
	return meta.Audio[len(meta.Audio)-1]
}

func getLastSubtitleMeta(meta *Meta) *SubtitleMeta {
	if len(meta.Subtitle) == 0 {
		panic("No subtitle available!")
	}
	return &meta.Subtitle[len(meta.Subtitle)-1]
}

func getLastChapterMeta(meta *Meta) *ChapterMeta {
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

func addAudioMeta(meta *Meta) {
	meta.Audio = append(meta.Audio, &AudioMeta{})
}

func addSubtitleMeta(meta *Meta) {
	subtitle := SubtitleMeta{}
	meta.Subtitle = append(meta.Subtitle, subtitle)
}

func addChapterMeta(meta *Meta) {
	meta.Chapter = append(meta.Chapter, ChapterMeta{})
}
