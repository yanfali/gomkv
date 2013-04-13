package handbrake

import (
	"gomkv/config"
	"strings"
	"testing"
)

func Test_ValidateProfile(t *testing.T) {
	meta := HandBrakeMeta{}
	conf := config.GomkvConfig{}
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
	meta := HandBrakeMeta{}
	conf := config.GomkvConfig{}
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
	meta := HandBrakeMeta{}
	conf := config.GomkvConfig{}
	conf.Profile = "Universal"
	conf.Prefix = ""
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
	meta := HandBrakeMeta{}
	conf := config.GomkvConfig{}
	conf.Profile = "Universal"
	conf.Prefix = ""
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
