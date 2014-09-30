package handbrake

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/yanfali/gomkv/config"
)

var (
	EmptyProfile   = errors.New("Encoding profile is empty!")
	EmptyTitle     = errors.New("Title is empty!")
	ParameterOrder = []string{"-Z ", "-i ", "-t", "-a", "-c", "-E ", "-s", "--subtitle-default ", "-o "}
)

const (
	CLI         = "HandBrakeCLI"
	ENCODE_FAAC = "faac"
	ENCODE_AC3  = "ffac3"
	COPY_AC3    = "copy:ac3"
	COPY_DTS    = "copy:dts"
	COPY_AAC    = "copy:aac"
	AC3         = "AC3"
	DTS         = "DTS"
	AAC         = "aac"
	FAAC        = "faac"
)

func addSubtitleOpts(hbcmd *handbrakeCommand, subtitlemeta []SubtitleMeta, config *config.GomkvConfig) error {
	if len(subtitlemeta) == 0 {
		return nil
	}

	subs := []int{}
	subdef := 0
	for i, subtitle := range subtitlemeta {
		if isCopyLanguage(subtitle.Language, config) {
			subs = append(subs, i+1)
			if subdef == 0 && config.DefaultSub == subtitle.Language {
				subdef = subs[i]
			}

		} else {
			continue
		}
	}
	toCopy := []string{}
	for _, track := range subs {
		toCopy = append(toCopy, strconv.Itoa(track))
	}

	if len(toCopy) > 0 {
		hbcmd.Params["-s"] = fmt.Sprintf("%s", strings.Join(toCopy, ","))
	}
	if subdef > 0 {
		hbcmd.Params["--subtitle-default "] = fmt.Sprintf("%s", strconv.Itoa(subdef))
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

func addAudioOpts(hbcmd *handbrakeCommand, audiometas AudioMetas, config *config.GomkvConfig) error {
	audioTracks := []int{}
	audioOptions := []string{}
	for _, audio := range audiometas {
		if isCopyLanguage(audio.Language, config) {
			audioTracks = append(audioTracks, audio.Index)
		} else {
			continue
		}
		if config.AacOnly {
			audioOptions = append(audioOptions, ENCODE_FAAC)
			continue
		}
		encoder := ""
		switch audio.Codec {
		case AC3:
			encoder = COPY_AC3
		case DTS:
			encoder = COPY_DTS
		case AAC:
			encoder = COPY_AAC
		default:
			encoder = ENCODE_AC3
		}
		audioOptions = append(audioOptions, encoder)
	}
	if !config.DisableAAC && len(audioTracks) > 0 && audioOptions[0] != FAAC {
		audioTracks = append(audioTracks, audioTracks[0])
		audioOptions = append(audioOptions, ENCODE_FAAC)
	}
	tracks := []string{}
	for _, track := range audioTracks {
		tracks = append(tracks, strconv.Itoa(track))
	}

	if len(tracks) > 0 {
		hbcmd.Params["-a"] = fmt.Sprintf("%s", strings.Join(tracks, ","))
		hbcmd.Params["-E "] = fmt.Sprintf("%s", strings.Join(audioOptions, ","))
	}
	return nil
}

func validateConfig(meta HandBrakeMeta, config *config.GomkvConfig) error {
	if config.Profile == "" {
		return EmptyProfile
	}
	if meta.Title == "" {
		return EmptyTitle
	}
	return nil
}

func FormatCLIOutput(meta HandBrakeMeta, config *config.GomkvConfig, session *config.GomkvSession) ([]string, error) {
	results := []string{}
	for {
		result, err := FormatCLIOutputEntry(meta, config, session)
		if err != nil {
			return results, err
		}
		results = append(results, result)
		if session.Chapter == 0 {
			// normal
			return results, nil
		}
		// split original into chapter groups until you run out of chapters
		session.Chapter += config.SplitFileEvery
		chpEnd := session.Chapter + config.SplitFileEvery - 1
		if chpEnd > len(meta.Chapter) {
			return results, nil
		}
	}
	panic("unreachable")
}

type handbrakeCommand struct {
	Params map[string]string
}

func ContainsString(search string, matches []string) bool {
	for _, match := range matches {
		if strings.Contains(search, match) {
			return true
		}
	}
	return false
}

func MakeFormatString(height int, title string, config *config.GomkvConfig) (format string) {
	if config.Profile != "Universal" && !(ContainsString(title, []string{".480p.", ".720p.", ".1080p.", ".4k."})) {
		switch {
		case height <= 480:
			format = ".480p"
		case height <= 720:
			format = ".720p"
		case height <= 1080:
			format = ".1080p"
		default:
			format = ".4k"
		}
	}
	if config.M4v {
		format += ".m4v"
	} else {
		format += ".mkv"
	}
	return
}

func FormatCLIOutputEntry(meta HandBrakeMeta, config *config.GomkvConfig, session *config.GomkvSession) (string, error) {
	hbcmd := handbrakeCommand{
		make(map[string]string),
	}

	if err := validateConfig(meta, config); err != nil {
		return "", err
	}
	title := meta.Title
	hbcmd.Params["-Z "] = fmt.Sprintf("%q", config.Profile)
	hbcmd.Params["-i "] = fmt.Sprintf("%q", title)
	hbcmd.Params["-t"] = "1"

	// TODO Make this smarter
	// - deal with overwriting same path
	var format string = MakeFormatString(meta.Height, title, config)

	var output string
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
			output = fmt.Sprintf("%s_S%dE%02d%s", config.Prefix, config.SeasonOffset, session.Episode, format)
			session.Episode += 1
			if session.Chapter > 0 {
				end := session.Chapter + config.SplitFileEvery - 1
				if (len(meta.Chapter)) >= end {
					hbcmd.Params["-c"] = fmt.Sprintf("%d-%d", session.Chapter, session.Chapter+config.SplitFileEvery-1)
				}
			}
		} else {
			output = config.Prefix + format
		}
	}

	if len(meta.Audio) > 0 {
		if config.Languages != "" {
			sort.Sort(ByLanguage{meta.Audio, config.LanguageOrderMap()})
		}
		addAudioOpts(&hbcmd, meta.Audio, config)
	}

	if config.EnableSubs {
		addSubtitleOpts(&hbcmd, meta.Subtitle, config)
	}

	if config.DestDir != "" {
		output = fmt.Sprintf("%s%c%s", config.DestDir, os.PathSeparator, output)
	}
	if output != "" && output == title {
		index := strings.LastIndex(output, ".")
		output = output[:index] + "_new" + output[index:]
	}
	hbcmd.Params["-o "] = fmt.Sprintf("%q", output)
	buf := bytes.NewBuffer([]byte{})
	fmt.Fprintf(buf, "%s ", CLI)
	for _, key := range ParameterOrder {
		value, ok := hbcmd.Params[key]
		if ok {
			fmt.Fprintf(buf, "%s%s ", key, value)
		}
	}
	return strings.Trim(buf.String(), " "), nil
}
