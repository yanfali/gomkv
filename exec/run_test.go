package exec

import (
	"strings"
	"testing"
	"time"
)

func Test_execls(t *testing.T) {
	if std, err := Command("ls"); err != nil {
		t.Errorf("unexpected err %s", err)
	} else {
		if len(std.Out) == 0 {
			t.Errorf("expected output, got %s", std.Out)
		}
		t.Log("ok")
	}
}

func Test_execSlow(t *testing.T) {
	if _, err := CommandWithTimeout("sleep", time.Second, "3"); err == nil {
		t.Error("err signal: killed")
	} else {
		if strings.HasSuffix(err.Error(), "signal: killed") {
			t.Log("ok")
		} else {
			t.Errorf("unexpected error %s", err)
		}
	}
}
