package handbrake

import (
	"gomkv/config"
	"strings"
	"testing"
)

func harness() (HandBrakeMeta, config.GomkvConfig) {
	return HandBrakeMeta{Title: "a.mkv"}, config.GomkvConfig{Profile: "Universal"}
}

type SimpleFunc func() (string, error)

func equals_harness(fn SimpleFunc, t *testing.T, expected string) {
	if result, err := fn(); err != nil {
		t.Errorf("unexpected error %s", err)
	} else {
		if result == expected {
			t.Log("ok")
			return
		}
		t.Errorf("Expected %s got '%s'", expected, result)
	}
}

func Test_ValidateProfile(t *testing.T) {
	meta, conf := HandBrakeMeta{}, config.GomkvConfig{}
	if _, err := FormatCLIOutput(meta, &conf); err != nil {
		if err == EmptyProfile {
			t.Log("ok")
			return
		}
		t.Errorf("unexpected error %s", err)
	} else {
		t.Error("expected encoding profile error")
	}
}

func Test_ValidateTitle(t *testing.T) {
	meta, conf := harness()
	meta.Title = ""
	if _, err := FormatCLIOutput(meta, &conf); err != nil {
		if err == EmptyTitle {
			t.Log("ok")
			return
		}
		t.Errorf("unexpected error %s", err)
	} else {
		t.Error("expected error title was empty")
	}
}

func Test_ValidatePrefixEmptyPassesTitle(t *testing.T) {
	meta, conf := harness()
	meta.Title = "RubberDuck"
	expected := "RubberDuck"
	if result, err := FormatCLIOutput(meta, &conf); err != nil {
		t.Errorf("unexpected error %s", err)
	} else {
		if strings.Contains(result, expected) {
			t.Log("ok")
			return
		}
		t.Errorf("Expected %s got '%s'", expected, result)
	}
}

func Test_ValidateTitleSpaceEscaped(t *testing.T) {
	meta, conf := harness()
	meta.Title = "Rubber Duck"
	expected := "Rubber\\ Duck"
	if result, err := FormatCLIOutput(meta, &conf); err != nil {
		t.Errorf("unexpected error %s", err)
	} else {
		if strings.Contains(result, expected) {
			t.Log("ok")
			return
		}
		t.Errorf("Expected %s got '%s'", expected, result)
	}
}

func Test_ValidateDestPathEscaped(t *testing.T) {
	meta, conf := harness()
	conf.DestDir = "/home/yanfali/My Video"
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o /home/yanfali/My\\ Video/a.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateBasicEncoding(t *testing.T) {
	meta, conf := harness()
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateFormatM4v(t *testing.T) {
	meta, conf := harness()
	conf.M4v = true
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o a.m4v"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidatePrefix(t *testing.T) {
	meta, conf := harness()
	conf.Prefix = "b"
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateEpisodes(t *testing.T) {
	meta, conf := harness()
	conf.Prefix = "b"
	conf.Episodic = true
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b_S0E00.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateSeasonOffset(t *testing.T) {
	meta, conf := harness()
	conf.Prefix = "b"
	conf.Episodic = true
	conf.SeasonOffset = 3
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b_S3E00.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateEpisodeOffset(t *testing.T) {
	meta, conf := harness()
	conf.Prefix = "b"
	conf.Episodic = true
	conf.EpisodeOffset = 15
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b_S0E15.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateBasicAudio(t *testing.T) {
	meta, conf := harness()
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	meta.Audio = AudioMetas{atrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1 -E copy:ac3 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateBasicAudioAacOnly(t *testing.T) {
	meta, conf := harness()
	conf.AacOnly = true
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	meta.Audio = AudioMetas{atrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1 -E faac -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateMobile(t *testing.T) {
	meta, conf := harness()
	conf.Mobile()
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o a.m4v"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateMobileWithAudio(t *testing.T) {
	meta, conf := harness()
	conf.Mobile()
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1 -E faac -o a.m4v"
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	meta.Audio = AudioMetas{atrack}
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateSrcDir(t *testing.T) {
	meta, conf := harness()
	conf.Mobile()
	conf.DestDir = "/tmp"
	meta.Title = "/home/beagle/a.mkv"
	expected := "HandBrakeCLI -Z \"Universal\" -i /home/beagle/a.mkv -t1 -o /tmp/a.m4v"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateAudio2TracksInEnglish(t *testing.T) {
	meta, conf := harness()
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	btrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 2}
	meta.Audio = AudioMetas{atrack, btrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1,2 -E copy:ac3,copy:ac3 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateAudio2TracksInEnglishOneInFrench(t *testing.T) {
	meta, conf := harness()
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	btrack := &AudioMeta{Language: "French", Codec: "AC3", Index: 2}
	ctrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 3}
	meta.Audio = AudioMetas{atrack, btrack, ctrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1,3 -E copy:ac3,copy:ac3 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateAudio2TracksInEnglishOneInJapanese(t *testing.T) {
	meta, conf := harness()
	conf.Languages = "Japanese,English"
	atrack := &AudioMeta{Language: "English", Codec: "DTS", Index: 1}
	btrack := &AudioMeta{Language: "Japanese", Codec: "AC3", Index: 2}
	ctrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 3}
	meta.Audio = AudioMetas{atrack, btrack, ctrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a2,1,3 -E copy:ac3,copy:dts,copy:ac3 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateSubtitleWithLangugageOption(t *testing.T) {
	meta, conf := harness()
	conf.Languages = "Japanese,English"
	conf.EnableSubs = true
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,2,3 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_ValidateSubtitleDefault(t *testing.T) {
	meta, conf := harness()
	conf.EnableSubs = true
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,3 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_DefaultSubtitle(t *testing.T) {
	meta, conf := harness()
	conf.EnableSubs = true
	conf.DefaultSub = "English"
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,3 --subtitle-default 1 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_DefaultSubtitleJapaneseWithNoLanguage(t *testing.T) {
	meta, conf := harness()
	conf.EnableSubs = true
	conf.DefaultSub = "Japanese"
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,3 -o a_new.mkv"

	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_DefaultSubtitleJapaneseWithLanguage(t *testing.T) {
	meta, conf := harness()
	conf.EnableSubs = true
	conf.DefaultSub = "Japanese"
	conf.Languages = "Japanese,English"
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,2,3 --subtitle-default 2 -o a_new.mkv"

	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}

func Test_DefaultSubtitleJapaneseWithLanguageAndAudio(t *testing.T) {
	meta, conf := harness()
	conf.EnableSubs = true
	conf.DefaultSub = "Japanese"
	conf.Languages = "Japanese,English"
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	btrack := &AudioMeta{Language: "Japanese", Codec: "AC3", Index: 2}
	meta.Audio = AudioMetas{atrack, btrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a2,1 -E copy:ac3,copy:ac3 -s 1,2,3 --subtitle-default 2 -o a_new.mkv"

	equals_harness(func() (string, error) {
		return FormatCLIOutput(meta, &conf)
	}, t, expected)
}
