package handbrake

import (
	"testing"
)

var data2 string = `+ title 1:
  + stream: TRANSFORMERS_1984_S1E01.mkv
  + duration: 00:22:52
  + size: 704x480, pixel aspect: 355/396, display aspect: 1.31, 29.970 fps
  + autocrop: 0/0/0/0
  + chapters:
    + 1: cells 0->0, 0 blocks, duration 00:00:34
    + 2: cells 0->0, 0 blocks, duration 00:05:14
    + 3: cells 0->0, 0 blocks, duration 00:02:47
    + 4: cells 0->0, 0 blocks, duration 00:06:41
    + 5: cells 0->0, 0 blocks, duration 00:06:32
    + 6: cells 0->0, 0 blocks, duration 00:01:01
  + audio tracks:
    + 1, English (AC3) (5.1 ch) (iso639-2: eng), 48000Hz, 448000bps
  + subtitle tracks:
`

var meta2 = ParseOutput(data2)

func Test_parseTitleTransformers(t *testing.T) {
	exp := "TRANSFORMERS_1984_S1E01.mkv"
	if meta2.Title == exp {
		t.Log("ok")
	} else {
		t.Errorf("expected '%s' - got '%s'", exp, meta2.Title)
	}
}

func Test_subTitleMissing(t *testing.T) {
	if len(meta2.Subtitle) == 0 {
		t.Log("ok")
	} else {
		t.Error("unexpected length of Subtitle %d", len(meta2.Subtitle))
	}
}

func Test_audioLength(t *testing.T) {
	if len(meta2.Audio) == 1 {
		t.Log("ok")
	} else {
		t.Error("unexpected length of Audio tracks %d", len(meta2.Audio))
	}
}
