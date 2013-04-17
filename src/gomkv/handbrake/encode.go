package handbrake

import (
	"bytes"
	"errors"
	"fmt"
	"gomkv/config"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sort"
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

func isCopyLanguage(lang string, config *config.GomkvConfig) bool {
	if config.Languages == "" && lang == "English" {
		return true
	}
	if strings.Contains(config.Languages, lang) {
		return true
	}
	return false
}

func addAudioOpts(buf *bytes.Buffer, audiometas AudioMetas, config *config.GomkvConfig) error {
	audioTracks := []int{}
	audioOptions := []string{}
	for _, audio := range audiometas {
		if isCopyLanguage(audio.Language, config) {
			audioTracks = append(audioTracks, audio.Index)
		} else {
			continue
		}
		if config.AacOnly {
			audioOptions = append(audioOptions, "faac")
			continue
		}
		switch audio.Codec {
		case "AC3":
			audioOptions = append(audioOptions, "copy:ac3")
		case "DTS":
			audioOptions = append(audioOptions, "copy:dts")
		case "aac":
			audioOptions = append(audioOptions, "copy:aac")
		case "pcm_s24le":
			audioOptions = append(audioOptions, "ffac3")
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

var EmptyProfile = errors.New("Encoding profile is empty!")
var EmptyTitle = errors.New("Title is empty!")

func validateConfig(meta HandBrakeMeta, config *config.GomkvConfig) error {
	if config.Profile == "" {
		return EmptyProfile
	}
	if meta.Title == "" {
		return EmptyTitle
	}
	return nil
}

func FormatCLIOutput(meta HandBrakeMeta, config *config.GomkvConfig) (string, error) {
	if err := validateConfig(meta, config); err != nil {
		return "", err
	}
	buf := bytes.NewBuffer([]byte{})
	title := strings.Replace(meta.Title, " ", "\\ ", -1)
	fmt.Fprintf(buf, "%s", CLI)
	fmt.Fprintf(buf, " -Z \"%s\"", config.Profile)
	fmt.Fprintf(buf, " -i %s", title)
	fmt.Fprintf(buf, " -t1")

	// TODO Make this smarter
	// - deal with overwriting same path
	var output string
	var format string
	if config.M4v {
		format = ".m4v"
	} else {
		format = ".mkv"
	}

	if config.Prefix == "" {
		output = filepath.Base(title)
		i := strings.LastIndex(output, ".")
		if i == -1 {
			output += format
		} else {
			output = output[:i] + format
		}
	} else {
		if config.Episodic {
			output = fmt.Sprintf("%s_S%dE%02d%s", config.Prefix, config.SeasonOffset, config.EpisodeOffset, format)
			config.EpisodeOffset++
		} else {
			output = config.Prefix + format
		}
	}

	if len(meta.Audio) > 0 {
		if config.Languages != "" {
			sort.Sort(ByLanguage{meta.Audio, config.LanguageOrderMap()})
		}
		addAudioOpts(buf, meta.Audio, config)
	}

	if config.EnableSubs {
		addSubtitleOpts(buf, meta.Subtitle)
	}

	if config.DestDir != "" {
		output = fmt.Sprintf("%s%c%s", config.DestDir, os.PathSeparator, output)
	}
	output = strings.Replace(output, " ", "\\ ", -1)
	if output == title {
		index := strings.LastIndex(output, ".")
		output = output[:index] + "_new" + output[index:]
	}
	fmt.Fprintf(buf, " -o %s", output)

	return buf.String(), nil
}
