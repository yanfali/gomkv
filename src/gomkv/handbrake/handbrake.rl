package handbrake

import (
	"strings"
)

%%{
	machine handbrake;
	write data;
}%%

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
				audio.Index = len(meta.Audio)
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
