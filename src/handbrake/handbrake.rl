package handbrake

import (
	"fmt"
	"strings"
	"strconv"
)

%%{
	machine handbrake;
	write data;
}%%

type Section int
const (
	NONE Section = iota
	CHAPTER
	AUDIO
	SUBTITLE
)

func parseTime(timestring string) float64 {
	rawTime := strings.Trim(timestring, " \n");
	splitTime := strings.Split(rawTime, ":")
	var length float64
	hours, err := strconv.ParseInt(splitTime[0], 10, 8)
	if err != nil {
		panic(err)
	}
	length += float64(hours * 60 * 60)
	minutes, err := strconv.ParseInt(splitTime[1], 10, 8)
	length += float64(minutes * 60)
	seconds, err := strconv.ParseInt(splitTime[2], 10, 8)
	length += float64(seconds)
	return length
}

func parseInt(value string) int {
	if result, err := strconv.ParseInt(value, 10, 32); err != nil {
		panic(err)
	} else {
		return int(result)
	}
	return 0
}

func getLastAudioMeta(meta *HandBrakeMeta) *AudioMeta {
	if len(meta.Audio) == 0 {
		panic("No audio available!")
	}
	return &meta.Audio[len(meta.Audio)-1]
}

func getLastSubtitleMeta(meta *HandBrakeMeta) *SubtitleMeta {
	if len(meta.Subtitle) == 0 {
		panic("No subtitle available!")
	}
	return &meta.Subtitle[len(meta.Subtitle)-1]
}

var debugEnabled bool = true
func debug(format string, args... interface{}) {
	if debugEnabled {
		fmt.Printf(format, args...)
	}
}

func addAudioMeta(meta *HandBrakeMeta) {
	audio := AudioMeta{}
	meta.Audio = append(meta.Audio, audio)
}

func addSubtitleMeta(meta *HandBrakeMeta) {
	subtitle := SubtitleMeta{}
	meta.Subtitle = append(meta.Subtitle, subtitle)
}

func ParseOutput(data string) HandBrakeMeta {
	cs, p, pe, eof := 0, 0, len(data), 0
	top, ts, te, act := 0,0,0,0
	var stack = []int{0}
	var section = NONE
	var meta = HandBrakeMeta{}
	line := 1
	debug("%02d: ", line)

	%%{
		action newline { line +=1; debug("\n%02d: ", line) }
		newline = any* '\n' @ newline;
		stitle := |*
			([^.])+[.]alnum+ => {
				meta.Title = strings.Trim(data[ts:te], " \n");
				debug("%s", meta.Title);
				fret;
			};
		*|;
		sduration := |*
			space+;
			digit{2}[:]digit{2}[:]digit{2} => {
				meta.Duration = parseTime(data[ts:te])
				debug("%f", meta.Duration);
				fret;
			};
		*|;
		picture := |*
			space*;
			digit{3,4} "x" digit{3,4} "," => {
				raw := data[ts:te-1];
				values := strings.Split(raw, "x");
				meta.Width = parseInt(values[0]);
				meta.Height = parseInt(values[1]);
				debug("1-%s:", data[ts:te-1]);
			};
			space*;
			"pixel" space+ "aspect:";
			digit{1,4} "/" digit{1,4} "," => {
				meta.Pixelaspect = data[ts:te-1]
				debug("2-%s:", meta.Pixelaspect);
			};
			space*;
			"display" space+ "aspect:" space*;
			digit . "." . digit{1,3} "," => {
				meta.Aspect = data[ts:te-1]
				debug("3-%s:", data[ts:te-1])
			};
			space*;
			digit{2} . "." digit{3} space+ "fps" "\n" => {
				raw := data[ts:te-5]
				meta.Fps = raw
				debug("4-%s", meta.Fps)
				p -= 1; fret;
			};
		*|;
		crop := |*
			space*;
			digit{1,3} "/" digit{1,3} "/" digit{1,3} "/" digit{1,3} => { debug("%s", data[ts:te]); fret; };
		*|;
		atrack := |*
# Language
			[A-Z] alpha+ => {
				addAudioMeta(&meta)
				audio := getLastAudioMeta(&meta)
				audio.Language = data[ts:te]
				debug("a-%s:", audio.Language)
			};
			space;
# Codec
			"(" ("AC3" | "DTS" ) ")" => {
				audio := getLastAudioMeta(&meta)
				audio.Codec = data[ts+1:te-1]
				debug("b-%s:", audio.Codec);
			};

			space;
# Channels
			"(" digit . "." . digit space "ch)" space => {
				audio := getLastAudioMeta(&meta)
				audio.Channels = data[ts+1:te-5]
				debug("c-%s:", audio.Channels)
			};
# Ignore this bit
			"(iso" digit{3} "-" digit ":" space lower{3} "),";
# Hertz
			space;
			digit+ "Hz," => {
				audio := getLastAudioMeta(&meta)
				audio.Frequency = parseInt(data[ts:te-3])
				debug("d-%d:", audio.Frequency)
			};
			space;
# Bps
			digit+ "bps" => {
				audio := getLastAudioMeta(&meta)
				audio.Bps = parseInt(data[ts:te-3])
		 		debug("e-%d", audio.Bps)
				fret;
			};
		*|;
		subtype = "Bitmap" | "Text";
		format = "VOBSUB" | "UTF-8";
		subtitle := |*
			[A-Z] alpha+ - subtype - format => {
				addSubtitleMeta(&meta)
				subtitle := getLastSubtitleMeta(&meta)
				subtitle.Language = data[ts:te];
				debug("a-%s:", subtitle.Language);
			};
			space;
			"(iso" digit{3} "-" digit ":" space lower{3} ")";
			space;
			"(" subtype ")" => {
				subtitle := getLastSubtitleMeta(&meta)
				subtitle.Type = data[ts+1:te-1]
				debug("b-%s:", subtitle.Type)
			};
			"(" format ")" => {
				subtitle := getLastSubtitleMeta(&meta)
				subtitle.Format = data[ts+1:te-1]
				debug("c-%s", subtitle.Format)
				fret;
			};
		*|;
		word = [a-z]+;
		prefix = space+ "+";
		prefixsp = prefix space;
		stream = prefixsp "stream:";
		duration = prefixsp "duration:";
		video = prefixsp "size:";
		autocrop = prefixsp "autocrop:";
		track = prefixsp digit "," space;
		main := ( 
			newline |
			stream @{ fcall stitle; } |
			duration @{ fcall sduration; } |
			video @{ fcall picture; } |
			autocrop @{ fcall crop; } |
			prefixsp "chapters:" @{ section = CHAPTER; debug("chapter"); } |
			prefixsp "audio tracks:" @{ section = AUDIO; debug("audio"); } |
			prefixsp "subtitle tracks:" @{ section = SUBTITLE; debug("subtitle"); } |
			track @{
				switch section {
				case AUDIO:
					fcall atrack;
				case SUBTITLE:
					fcall subtitle;
				}
			}
		)*;
		write init;
		write exec;
	}%%
	return meta
}
