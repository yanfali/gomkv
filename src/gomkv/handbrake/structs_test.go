package handbrake

import (
	"sort"
	"testing"
)

func Test_sort(t *testing.T) {
	atrack := AudioMeta{"English", "AC3", "5.1", 0, 0, 1}
	btrack := AudioMeta{"Spanish", "AC3", "5.1", 0, 0, 2}
	order := map[string]int{"Spanish": 0, "English": 1}
	metas := AudioMetas{&atrack, &btrack}
	sort.Sort(ByLanguage{metas, order})
	if metas[0].Language == "Spanish" && metas[1].Language == "English" {
		t.Log("ok")
	} else {
		t.Error("wrong order")
	}
}

func Test_sortUnSpecifiedLanguage(t *testing.T) {
	atrack := AudioMeta{"English", "AC3", "5.1", 0, 0, 1}
	btrack := AudioMeta{"Spanish", "AC3", "5.1", 0, 0, 2}
	ctrack := AudioMeta{"Japanese", "AC3", "5.1", 0, 0, 3}
	dtrack := AudioMeta{"Japanese", "AC3", "2.1", 0, 0, 4}
	order := map[string]int{"Spanish": 0, "English": 1}
	metas := AudioMetas{&atrack, &btrack, &ctrack, &dtrack}
	sort.Sort(ByLanguage{metas, order})
	if metas[0].Language == "Spanish" && metas[1].Language == "English" && metas[2].Language == "Japanese" {
		t.Log("ok")
	} else {
		t.Errorf("wrong order", metas)
	}
}

func Test_sortSplitSoundtracks(t *testing.T) {
	atrack := AudioMeta{"English", "AC3", "5.1", 0, 0, 1}
	btrack := AudioMeta{"Spanish", "AC3", "5.1", 0, 0, 2}
	ctrack := AudioMeta{"Japanese", "AC3", "5.1", 0, 0, 3}
	dtrack := AudioMeta{"English", "AC3", "2.1", 0, 0, 4}
	order := map[string]int{"Spanish": 0, "English": 1}
	metas := AudioMetas{&atrack, &btrack, &ctrack, &dtrack}
	sort.Sort(ByLanguage{metas, order})
	if metas[0].Language == "Spanish" &&
	   metas[1].Language == "English" &&
	   metas[2].Language == "English" &&
	   metas[2].Channels == "2.1" {
		t.Log("ok")
	} else {
		t.Errorf("wrong order %s", metas)
	}
}
