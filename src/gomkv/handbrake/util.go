package handbrake

import (
	"fmt"
	"strconv"
	"strings"
)

type Section int

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
	return &meta.Audio[len(meta.Audio)-1]
}

func getLastSubtitleMeta(meta *HandBrakeMeta) *SubtitleMeta {
	if len(meta.Subtitle) == 0 {
		panic("No subtitle available!")
	}
	return &meta.Subtitle[len(meta.Subtitle)-1]
}

var debugEnabled bool = true

func debug(format string, args ...interface{}) {
	if debugEnabled {
		fmt.Printf(format, args...)
	}
}

func addAudioMeta(meta *HandBrakeMeta) {
	audio := AudioMeta{}
	meta.Audio = append(meta.Audio, audio)
}

func addSubtitleMeta(meta *HandBrakeMeta) {
	subtitle := SubtitleMeta{}
	meta.Subtitle = append(meta.Subtitle, subtitle)
}
