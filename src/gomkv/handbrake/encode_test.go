package handbrake

import (
	"gomkv/config"
	"strings"
	"testing"
)

func harness() (HandBrakeMeta, config.GomkvConfig) {
	return HandBrakeMeta{}, config.GomkvConfig{}
}

func Test_ValidateProfile(t *testing.T) {
	meta, conf := harness()
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
	conf.Profile = "Universal"
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
	conf.Profile = "Universal"
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
	conf.Profile = "Universal"
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

func Test_ValidateBasicEncoding(t *testing.T) {
	meta, conf := harness()
	conf.Profile = "Universal"
	meta.Title = "a.mkv"
	expected := "HandBrakeCLI -Z Universal -i a.mkv -t1 -o a.mkv"
	if result, err := FormatCLIOutput(meta, &conf); err != nil {
		t.Errorf("unexpected error %s", err)
	} else {
		if result == expected {
			t.Log("ok")
			return
		}
		t.Errorf("Expected %s got '%s'", expected, result)
	}
}

func Test_ValidatePrefix(t *testing.T) {
	meta, conf := harness()
	conf.Profile = "Universal"
	conf.Prefix = "b"
	meta.Title = "a.mkv"
	expected := "HandBrakeCLI -Z Universal -i a.mkv -t1 -o b.mkv"
	if result, err := FormatCLIOutput(meta, &conf); err != nil {
		t.Errorf("unexpected error %s", err)
	} else {
		if result == expected {
			t.Log("ok")
			return
		}
		t.Errorf("Expected %s got '%s'", expected, result)
	}
}

func Test_ValidateEpisodes(t *testing.T) {
	meta, conf := harness()
	conf.Profile = "Universal"
	conf.Prefix = "b"
	conf.Episodic = true
	meta.Title = "a.mkv"
	expected := "HandBrakeCLI -Z Universal -i a.mkv -t1 -o b_S0E00.mkv"
	if result, err := FormatCLIOutput(meta, &conf); err != nil {
		t.Errorf("unexpected error %s", err)
	} else {
		if result == expected {
			t.Log("ok")
			return
		}
		t.Errorf("Expected %s got '%s'", expected, result)
	}
}

func Test_ValidateSeasonOffset(t *testing.T) {
	meta, conf := harness()
	conf.Profile = "Universal"
	conf.Prefix = "b"
	conf.Episodic = true
	conf.SeasonOffset = 3
	meta.Title = "a.mkv"
	expected := "HandBrakeCLI -Z Universal -i a.mkv -t1 -o b_S3E00.mkv"
	if result, err := FormatCLIOutput(meta, &conf); err != nil {
		t.Errorf("unexpected error %s", err)
	} else {
		if result == expected {
			t.Log("ok")
			return
		}
		t.Errorf("Expected %s got '%s'", expected, result)
	}
}

func Test_ValidateEpisodeOffset(t *testing.T) {
	meta, conf := harness()
	conf.Profile = "Universal"
	conf.Prefix = "b"
	conf.Episodic = true
	conf.EpisodeOffset = 15 
	meta.Title = "a.mkv"
	expected := "HandBrakeCLI -Z Universal -i a.mkv -t1 -o b_S0E15.mkv"
	if result, err := FormatCLIOutput(meta, &conf); err != nil {
		t.Errorf("unexpected error %s", err)
	} else {
		if result == expected {
			t.Log("ok")
			return
		}
		t.Errorf("Expected %s got '%s'", expected, result)
	}
}
