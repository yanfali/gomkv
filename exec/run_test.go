package exec

import (
	"testing"
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
