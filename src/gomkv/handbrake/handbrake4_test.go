package handbrake

import (
	"testing"
)

var data4 = `+ title 1:
  + stream: title00.mkv
  + duration: 01:00:00
  + size: 720x480, pixel aspect: 32/27, display aspect: 1.78, 23.976 fps
  + autocrop: 0/0/0/0
  + chapters:
    + 1: cells 0->0, 0 blocks, duration 01:00:00
  + audio tracks:
    + 1, English (AC3) (5.1 ch) (iso639-2: eng), 48000Hz, 384000bps
    + 2, Francais (AC3) (5.1 ch) (iso639-2: fra), 48000Hz, 384000bps
    + 3, Portugues (AC3) (Dolby Surround) (iso639-2: por), 48000Hz, 192000bps
    + 4, Espanol (AC3) (5.1 ch) (iso639-2: spa), 48000Hz, 384000bps
    + 5, Thai (AC3) (5.1 ch) (iso639-2: tha), 48000Hz, 384000bps
  + subtitle tracks:
    + 1, English (iso639-2: eng) (Bitmap)(VOBSUB)
    + 2, Chinese (iso639-2: zho) (Bitmap)(VOBSUB)
    + 3, French (iso639-2: fra) (Bitmap)(VOBSUB)
    + 4, Portuguese (iso639-2: por) (Bitmap)(VOBSUB)
    + 5, Thai (iso639-2: tha) (Bitmap)(VOBSUB)
`

var meta4 = ParseOutput(data4)

func Test_HandleDolbySurroundAudio(t *testing.T) {
	if len(meta4.Audio) == 5 {
		t.Log("ok")
	} else {
		t.Errorf("expected 5 audio tracks, go %d", len(meta4.Audio))
	}
	if len(meta4.Subtitle) == 5 {
		t.Log("ok")
	} else {
		t.Error("expect 5 subtitles")
	}
	if meta4.Subtitle[0].Language == "English" &&
		meta4.Audio[2].Language == "Portugues" &&
		meta4.Audio[2].Channels == "Dolby Surround" {
		t.Log("ok")
	} else {
		t.Error("Expected Track 3 to be Portugues and dolby")
	}
}
