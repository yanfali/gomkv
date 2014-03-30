package handbrake

import (
	"gomkv/config"
	"strings"
	"testing"
)

func harness() (HandBrakeMeta, config.GomkvConfig, config.GomkvSession) {
	return HandBrakeMeta{Title: "a.mkv"}, config.GomkvConfig{Profile: "Universal"}, config.GomkvSession{}
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
	meta, conf, sess := HandBrakeMeta{}, config.GomkvConfig{}, config.GomkvSession{}
	if _, err := FormatCLIOutputEntry(meta, &conf, &sess); err != nil {
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
	meta, conf, sess := harness()
	meta.Title = ""
	if _, err := FormatCLIOutputEntry(meta, &conf, &sess); err != nil {
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
	meta, conf, sess := harness()
	meta.Title = "RubberDuck"
	expected := "RubberDuck"
	if result, err := FormatCLIOutputEntry(meta, &conf, &sess); err != nil {
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
	meta, conf, sess := harness()
	meta.Title = "Rubber Duck"
	expected := "Rubber\\ Duck"
	if result, err := FormatCLIOutputEntry(meta, &conf, &sess); err != nil {
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
	meta, conf, sess := harness()
	conf.DestDir = "/home/yanfali/My Video"
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o /home/yanfali/My\\ Video/a.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateBasicEncoding(t *testing.T) {
	meta, conf, sess := harness()
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateFormatM4v(t *testing.T) {
	meta, conf, sess := harness()
	conf.M4v = true
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o a.m4v"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidatePrefix(t *testing.T) {
	meta, conf, sess := harness()
	conf.Prefix = "b"
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateEpisodes(t *testing.T) {
	meta, conf, sess := harness()
	conf.Prefix = "b"
	conf.Episodic = true
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b_S0E00.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateSeasonOffset(t *testing.T) {
	meta, conf, sess := harness()
	conf.Prefix = "b"
	conf.Episodic = true
	conf.SeasonOffset = 3
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b_S3E00.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateEpisodeOffset(t *testing.T) {
	meta, conf, sess := harness()
	sess.Episode = 15
	conf.Prefix = "b"
	conf.Episodic = true
	conf.EpisodeOffset = 15
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b_S0E15.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateMultiEpisodeOffset(t *testing.T) {
	meta, conf, sess := harness()
	sess.Episode = 15
	conf.Prefix = "b"
	conf.Episodic = true
	conf.EpisodeOffset = 15
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b_S0E15.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
	expected = "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o b_S0E16.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateBasicAudio(t *testing.T) {
	meta, conf, sess := harness()
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	meta.Audio = AudioMetas{atrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1,1 -E copy:ac3,faac -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}
func Test_ValidateDisableAACAudio(t *testing.T) {
	meta, conf, sess := harness()
	conf.DisableAAC = true;
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	meta.Audio = AudioMetas{atrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1 -E copy:ac3 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}
func Test_AACAudioNoTracks(t *testing.T) {
	meta, conf, sess := harness()
	conf.DisableAAC = true;
	meta.Audio = AudioMetas{}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}
func Test_ValidateBasicAudioAacOnly(t *testing.T) {
	meta, conf, sess := harness()
	conf.AacOnly = true
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	meta.Audio = AudioMetas{atrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1 -E faac -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateMobile(t *testing.T) {
	meta, conf, sess := harness()
	conf.Mobile()
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o a.m4v"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateMobileWithAudio(t *testing.T) {
	meta, conf, sess := harness()
	conf.Mobile()
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1 -E faac -o a.m4v"
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	meta.Audio = AudioMetas{atrack}
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateSrcDir(t *testing.T) {
	meta, conf, sess := harness()
	conf.Mobile()
	conf.DestDir = "/tmp"
	meta.Title = "/home/beagle/a.mkv"
	expected := "HandBrakeCLI -Z \"Universal\" -i /home/beagle/a.mkv -t1 -o /tmp/a.m4v"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateAudio2TracksInEnglish(t *testing.T) {
	meta, conf, sess := harness()
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	btrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 2}
	meta.Audio = AudioMetas{atrack, btrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1,2,1 -E copy:ac3,copy:ac3,faac -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateAudio2TracksInEnglishOneInFrench(t *testing.T) {
	meta, conf, sess := harness()
	atrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 1}
	btrack := &AudioMeta{Language: "French", Codec: "AC3", Index: 2}
	ctrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 3}
	meta.Audio = AudioMetas{atrack, btrack, ctrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a1,3,1 -E copy:ac3,copy:ac3,faac -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateAudio2TracksInEnglishOneInJapanese(t *testing.T) {
	meta, conf, sess := harness()
	conf.Languages = "Japanese,English"
	atrack := &AudioMeta{Language: "English", Codec: "DTS", Index: 1}
	btrack := &AudioMeta{Language: "Japanese", Codec: "AC3", Index: 2}
	ctrack := &AudioMeta{Language: "English", Codec: "AC3", Index: 3}
	meta.Audio = AudioMetas{atrack, btrack, ctrack}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a2,1,3,2 -E copy:ac3,copy:dts,copy:ac3,faac -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateSubtitleWithLangugageOption(t *testing.T) {
	meta, conf, sess := harness()
	conf.Languages = "Japanese,English"
	conf.EnableSubs = true
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,2,3 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_ValidateSubtitleDefault(t *testing.T) {
	meta, conf, sess := harness()
	conf.EnableSubs = true
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,3 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_DefaultSubtitle(t *testing.T) {
	meta, conf, sess := harness()
	conf.EnableSubs = true
	conf.DefaultSub = "English"
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,3 --subtitle-default 1 -o a_new.mkv"
	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_DefaultSubtitleJapaneseWithNoLanguage(t *testing.T) {
	meta, conf, sess := harness()
	conf.EnableSubs = true
	conf.DefaultSub = "Japanese"
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,3 -o a_new.mkv"

	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_DefaultSubtitleJapaneseWithLanguage(t *testing.T) {
	meta, conf, sess := harness()
	conf.EnableSubs = true
	conf.DefaultSub = "Japanese"
	conf.Languages = "Japanese,English"
	asub := SubtitleMeta{Language: "English"}
	bsub := SubtitleMeta{Language: "Japanese"}
	csub := SubtitleMeta{Language: "English"}
	meta.Subtitle = []SubtitleMeta{asub, bsub, csub}
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -s 1,2,3 --subtitle-default 2 -o a_new.mkv"

	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func Test_DefaultSubtitleJapaneseWithLanguageAndAudio(t *testing.T) {
	meta, conf, sess := harness()
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
	expected := "HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -a2,1,2 -E copy:ac3,copy:ac3,faac -s 1,2,3 --subtitle-default 2 -o a_new.mkv"

	equals_harness(func() (string, error) {
		return FormatCLIOutputEntry(meta, &conf, &sess)
	}, t, expected)
}

func TestChapterSplitting(t *testing.T) {
	meta, conf, sess := harness()
	meta.Chapter = []ChapterMeta{ChapterMeta{Index: 1}, ChapterMeta{Index: 2}, ChapterMeta{Index: 3}, ChapterMeta{Index: 4}, ChapterMeta{Index: 5}, ChapterMeta{Index: 6}}
	conf.Prefix = "b"
	conf.Episodic = true
	conf.SplitFileEvery = 2
	sess.Episode = 1
	sess.Chapter = 1
	results, err := FormatCLIOutput(meta, &conf, &sess)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	} else {
		t.Log("ok - 3 results")
	}
	expected := []string{
		"HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -c1-2 -o b_S0E01.mkv",
		"HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -c3-4 -o b_S0E02.mkv",
		"HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -c5-6 -o b_S0E03.mkv",
	}
	for i := 0; i < len(expected); i += 1 {
		if results[i] == expected[i] {
			t.Logf("ok - found entry %d", i)
		} else {
			t.Errorf("Expected '%s', got '%s'", expected[i], results[i])
		}
	}
}

func TestChapterSplittingIgnoreRemainder(t *testing.T) {
	meta, conf, sess := harness()
	meta.Chapter = []ChapterMeta{ChapterMeta{Index: 1}, ChapterMeta{Index: 2}, ChapterMeta{Index: 3}, ChapterMeta{Index: 4}, ChapterMeta{Index: 5}}
	conf.Prefix = "c"
	conf.Episodic = true
	conf.SplitFileEvery = 2
	sess.Episode = 1
	sess.Chapter = 1
	results, err := FormatCLIOutput(meta, &conf, &sess)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	exp_len := 2
	if len(results) != exp_len {
		t.Errorf("Expected %d results, got %d", exp_len, len(results))
	} else {
		t.Logf("ok - %d results", exp_len)
	}
	expected := []string{
		"HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -c1-2 -o c_S0E01.mkv",
		"HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -c3-4 -o c_S0E02.mkv",
	}
	for i := 0; i < len(expected); i += 1 {
		if results[i] == expected[i] {
			t.Logf("ok - found entry %d", i)
		} else {
			t.Errorf("Expected '%s', got '%s'", expected[i], results[i])
		}
	}
}

func TestChapterSplittingIgnoredIfTooFewChapters(t *testing.T) {
	meta, conf, sess := harness()
	meta.Chapter = []ChapterMeta{ChapterMeta{Index: 1}}
	conf.Prefix = "c"
	conf.Episodic = true
	conf.SplitFileEvery = 2
	sess.Episode = 1
	sess.Chapter = 1
	results, err := FormatCLIOutput(meta, &conf, &sess)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	exp_len := 1
	if len(results) != exp_len {
		t.Errorf("Expected %d results, got %d", exp_len, len(results))
	} else {
		t.Logf("ok - %d results", exp_len)
	}
	expected := []string{
		"HandBrakeCLI -Z \"Universal\" -i a.mkv -t1 -o c_S0E01.mkv",
	}
	for i := 0; i < len(expected); i += 1 {
		if results[i] == expected[i] {
			t.Logf("ok - found entry %d", i)
		} else {
			t.Errorf("Expected '%s', got '%s'", expected[i], results[i])
		}
	}
}
