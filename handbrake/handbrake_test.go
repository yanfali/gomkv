package handbrake

import (
	"testing"
)

var data = `+ title 1:
  + stream: source code.480p.x264.mkv
  + duration: 01:33:09
  + size: 720x480, pixel aspect: 2560/2151, display aspect: 1.79, 29.970 fps
  + autocrop: 0/0/0/0
  + chapters:
    + 1: cells 0->0, 0 blocks, duration 00:06:52
    + 2: cells 0->0, 0 blocks, duration 00:04:11
    + 3: cells 0->0, 0 blocks, duration 00:06:18
    + 4: cells 0->0, 0 blocks, duration 00:03:43
    + 5: cells 0->0, 0 blocks, duration 00:07:30
    + 6: cells 0->0, 0 blocks, duration 00:03:13
    + 7: cells 0->0, 0 blocks, duration 00:04:28
    + 8: cells 0->0, 0 blocks, duration 00:05:16
    + 9: cells 0->0, 0 blocks, duration 00:05:46
    + 10: cells 0->0, 0 blocks, duration 00:05:08
    + 11: cells 0->0, 0 blocks, duration 00:02:02
    + 12: cells 0->0, 0 blocks, duration 00:03:28
    + 13: cells 0->0, 0 blocks, duration 00:05:00
    + 14: cells 0->0, 0 blocks, duration 00:06:37
    + 15: cells 0->0, 0 blocks, duration 00:05:32
    + 16: cells 0->0, 0 blocks, duration 00:06:37
    + 17: cells 0->0, 0 blocks, duration 00:02:51
    + 18: cells 0->0, 0 blocks, duration 00:08:28
  + audio tracks:
    + 1, English (AC3) (5.1 ch) (iso639-2: eng), 48000Hz, 448000bps
  + subtitle tracks:
    + 1, English (iso639-2: eng) (Bitmap)(VOBSUB)
    + 2, Spanish (iso639-2: spa) (Bitmap)(VOBSUB)
    + 3, English (iso639-2: eng) (Bitmap)(VOBSUB)
`

var meta = ParseOutput(data)

func Test_parseTitle(t *testing.T) {
	exp := "source code.480p.x264.mkv"
	if meta.Title == exp {
		t.Log("ok")
	} else {
		t.Errorf("expected '%s' - got '%s'", exp, meta.Title)
	}
}

func Test_parseTime(t *testing.T) {
	if meta.Duration == 5589.0 {
		t.Log("ok")
	} else {
		t.Errorf("expected 5589, got %f", meta.Duration)
	}
}

func Test_parseWidthHeight(t *testing.T) {
	expHeight, expWidth := 480, 720
	if meta.Height == expHeight && meta.Width == expWidth {
		t.Log("ok")
	} else {
		t.Errorf("expected %d x %d, got %d x %d", expWidth, expHeight, meta.Width, meta.Height)
	}
}

func Test_parsePixelAspect(t *testing.T) {
	expected := "2560/2151"
	if meta.Pixelaspect == expected {
		t.Log("ok")
	} else {
		t.Errorf("expected %s, got %s", expected, meta.Pixelaspect)
	}
}

func Test_aspectRatio(t *testing.T) {
	ar := "1.79"
	if meta.Aspect == ar {
		t.Log("ok")
	} else {
		t.Errorf("expected %s, got %s", ar, meta.Aspect)
	}
}

func Test_fps(t *testing.T) {
	fps := "29.970"
	if meta.Fps == fps {
		t.Log("ok")
	} else {
		t.Errorf("expected %s, got %s", fps, meta.Fps)
	}
}

func Test_AudioLanguage(t *testing.T) {
	if len(meta.Audio) != 1 {
		t.Error("No audio track found")
	}
	if meta.Audio[0].Language == "English" {
		t.Log("ok")
	} else {
		t.Errorf("expected English, got '%s'", meta.Audio[0].Language)
	}
}

func Test_AudioFrequency(t *testing.T) {
	if meta.Audio[0].Frequency == 48000 {
		t.Log("ok")
	} else {
		t.Errorf("expected 48KHz, got %d", meta.Audio[0].Frequency)
	}
}

func Test_AudioBps(t *testing.T) {
	if meta.Audio[0].Bps == 448000 {
		t.Log("ok")
	} else {
		t.Errorf("expected 448KBps, got %d", meta.Audio[0].Bps)
	}
}

func Test_SubtitleCount(t *testing.T) {
	if len(meta.Subtitle) == 3 {
		t.Log("ok")
	} else {
		t.Errorf("expected 3 subtitles, found %d", len(meta.Subtitle))
	}
}

func Test_SubtitleLanguage(t *testing.T) {
	expected := []string{"English", "Spanish", "English"}

	for index, subtitle := range meta.Subtitle {
		if subtitle.Language != expected[index] {
			t.Errorf("expected %s as language %d, got '%s'", expected[index], index, subtitle.Language)
		}
	}
	t.Log("ok")
}
