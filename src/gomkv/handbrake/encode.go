package handbrake

import (
	"bytes"
	"fmt"
	"gomkv/config"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	CLI = "HandBrakeCLI"
)

func addSubtitleOpts(buf *bytes.Buffer, subtitlemeta []SubtitleMeta) error {
	if len(subtitlemeta) == 0 {
		return nil
	}
	subs := []int{}
	for i, subtitle := range subtitlemeta {
		if subtitle.Language == "English" {
			subs = append(subs, i+1)
		} else {
			continue
		}
	}
	toCopy := []string{}
	for _, track := range subs {
		toCopy = append(toCopy, strconv.Itoa(track))
	}

	if len(toCopy) > 0 {
		fmt.Fprintf(buf, " -s %s", strings.Join(toCopy, ","))
	}
	return nil
}

func addAudioOpts(buf *bytes.Buffer, audiometa []AudioMeta) error {
	if len(audiometa) == 0 {
		return nil
	}
	audioTracks := []int{}
	audioOptions := []string{}
	for i, audio := range audiometa {
		if audio.Language == "English" {
			audioTracks = append(audioTracks, i+1)
		} else {
			continue
		}
		switch audio.Codec {
		case "AC3":
			audioOptions = append(audioOptions, "copy:ac3")
		case "DTS":
			audioOptions = append(audioOptions, "copy:dts")
		}
	}
	tracks := []string{}
	for _, track := range audioTracks {
		tracks = append(tracks, strconv.Itoa(track))
	}

	if len(tracks) > 0 {
		fmt.Fprintf(buf, " -a%s", strings.Join(tracks, ","))
		fmt.Fprintf(buf, " -E %s", strings.Join(audioOptions, ","))
	}
	return nil
}

func FormatCLIOutput(meta HandBrakeMeta, config *config.GomkvConfig) (string, error) {
	buf := bytes.NewBuffer([]byte{})
	title := strings.Replace(meta.Title, " ", "\\ ", -1)
	fmt.Fprintf(buf, "%s", CLI)
	fmt.Fprintf(buf, " -Z %s", config.Profile)
	fmt.Fprintf(buf, " -i %s", title)
	fmt.Fprintf(buf, " -t1")

	// TODO Make this smarter
	// - deal with overwriting same path
	// - deal with episodes
	var output string
	if config.Prefix == "" {
		output = filepath.Base(title)
	} else {
		if config.Episodic {
			output = fmt.Sprintf("%s_S%dE%02d.mkv", config.Prefix, config.SeasonOffset, config.EpisodeOffset)
			config.EpisodeOffset++
		} else {
			output = config.Prefix + ".mkv"
		}
	}

	addAudioOpts(buf, meta.Audio)
	addSubtitleOpts(buf, meta.Subtitle)

	fmt.Fprintf(buf, " -o %s", output)
	fmt.Fprintf(buf, "\n")

	return buf.String(), nil
}
