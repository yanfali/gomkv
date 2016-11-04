package handbrake

import (
	"testing"
)

var data3 = `+ title 1:
  + stream: GODOT_t00.mkv
  + duration: 01:05:29
  + size: 720x480, pixel aspect: 1/1, display aspect: 1.78, 23.976 fps
  + autocrop: 0/0/0/0
  + chapters:
    + 1: cells 0->0, 0 blocks, duration 01:05:29
  + audio tracks:
    + 1, Japanese (pcm_s24le) (1.0 ch) (iso639-2: jpn)
    + 2, English (AC3) (1.0 ch) (iso639-2: eng), 48000Hz, 192000bps
  + subtitle tracks:
`

var meta3 = ParseOutput(data3)

func Test_parseTitleOddAudio(t *testing.T) {
	audio := meta3.Audio
	if len(audio) == 2 {
		t.Log("ok")
	} else {
		t.Errorf("not enough tracks, expected 2 got %d\n", len(audio))
	}
}

func Test_parseTitleOddAudioTracks(t *testing.T) {
	audio := meta3.Audio
	if audio[0].Language == "Japanese" &&
		audio[0].Codec == "pcm_s24le" &&
		audio[1].Language == "English" &&
		audio[1].Codec == "AC3" {
		t.Log("ok")
	} else {
		t.Errorf("Not parsed correctly %+v %+v", audio[0], audio[1])
	}
}
