package handbrake

import (
	"testing"
)

var data5 = `+ title 1:
  + stream: title00.mkv
  + duration: 02:01:42
  + size: 720x480, pixel aspect: 32/27, display aspect: 1.78, 29.970 fps
  + autocrop: 0/0/2/0
  + chapters:
    + 1: cells 0->0, 0 blocks, duration 00:02:05
    + 2: cells 0->0, 0 blocks, duration 00:10:33
    + 3: cells 0->0, 0 blocks, duration 00:10:21
    + 4: cells 0->0, 0 blocks, duration 00:01:30
    + 5: cells 0->0, 0 blocks, duration 00:00:16
    + 6: cells 0->0, 0 blocks, duration 00:01:30
    + 7: cells 0->0, 0 blocks, duration 00:11:11
    + 8: cells 0->0, 0 blocks, duration 00:09:43
    + 9: cells 0->0, 0 blocks, duration 00:01:30
    + 10: cells 0->0, 0 blocks, duration 00:00:16
    + 11: cells 0->0, 0 blocks, duration 00:01:30
    + 12: cells 0->0, 0 blocks, duration 00:11:10
    + 13: cells 0->0, 0 blocks, duration 00:09:45
    + 14: cells 0->0, 0 blocks, duration 00:01:29
    + 15: cells 0->0, 0 blocks, duration 00:00:16
    + 16: cells 0->0, 0 blocks, duration 00:01:29
    + 17: cells 0->0, 0 blocks, duration 00:10:14
    + 18: cells 0->0, 0 blocks, duration 00:10:41
    + 19: cells 0->0, 0 blocks, duration 00:01:30
    + 20: cells 0->0, 0 blocks, duration 00:00:16
    + 21: cells 0->0, 0 blocks, duration 00:01:30
    + 22: cells 0->0, 0 blocks, duration 00:11:35
    + 23: cells 0->0, 0 blocks, duration 00:09:19
    + 24: cells 0->0, 0 blocks, duration 00:01:29
    + 25: cells 0->0, 0 blocks, duration 00:00:21
    + 26: cells 0->0, 0 blocks, duration 00:00:02
  + audio tracks:
    + 1, English (AC3) (2.0 ch) (iso639-2: eng), 48000Hz, 192000bps
`

var meta5 = ParseOutput(data5)

func Test_Chapters(t *testing.T) {
	if len(meta5.Chapter) == 26 {
		t.Log("ok")
	} else {
		t.Errorf("expected 26 chapters, got %d", len(meta5.Chapter))
	}
}
