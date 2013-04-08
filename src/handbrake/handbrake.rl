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

func parseOutput(data string) HandBrakeMeta {
	cs, p, pe, eof := 0, 0, len(data), 0
	top, ts, te, act := 0,0,0,0
	_,_,_,_ = top, ts, te, act
	var stack = []int{0}
	var section = NONE
	var meta = HandBrakeMeta{}
	_ = eof
	line := 1
	csp := 0
	fmt.Printf("%02d: ", line)
	_ = csp
	%%{
		action newline { line +=1; fmt.Printf("\n%02d: ", line) }
		newline = any* '\n' @ newline;
		stitle := |*
			([^.])+[.]alnum+ => {
				meta.Title = strings.Trim(data[ts:te], " \n");
				fmt.Printf("%s", meta.Title);
				fret;
			};
		*|;
		sduration := |*
			space+;
			digit{2}[:]digit{2}[:]digit{2} => {
				meta.Duration = parseTime(data[ts:te])
				fmt.Printf("%f", meta.Duration);
				fret;
			};
		*|;
		picture := |*
			space*;
			digit{3,4} "x" digit{3,4} => { fmt.Printf("%s:", data[ts:te]); fret;};
		*|;
		paspect := |*
			space*;
			digit{1,4} "/" digit{1,4} => { fmt.Printf("%s:", data[ts:te]); fret; };
		*|;
		daspect := |*
			space*;
			digit . "." . digit{1,3} => { fmt.Printf("%s:", data[ts:te]); fret; };
		*|;
		sfps := |*
			"\n" => { ts -= 10; fmt.Printf("%s", data[ts:te-5]); p -= 1; fret; };
		*|;
		crop := |*
			space*;
			digit{1,3} "/" digit{1,3} "/" digit{1,3} "/" digit{1,3} => { fmt.Printf("%s", data[ts:te]); fret; };
		*|;
		atrack := |*
#			space+ "+" space+;
# track
#			digit => { fmt.Printf("%s,", data[ts:te]); };
#			"," space+;
# Language
			[A-Z] alpha+ => { fmt.Printf("a-%s:", data[ts:te]); };
			space;
			"(";
# Codec
			"AC3" | "DTS" | "aac" => { fmt.Printf("b-%s:", data[ts:te]); };
			")";
			space;
			"(";
# Channels
			digit . "." . digit => { fmt.Printf("c-%s:", data[ts:te]) };
			space;
			"ch) ";
# Ignore this bit
			"(iso" digit{3} "-" digit ":" space lower{3} "),";
# Hertz
			space;
			digit+ "Hz," => { fmt.Printf("d-%s:", data[ts:te-3]) };
			space;
# Bps
			digit+ "bps" => { fmt.Printf("e-%s", data[ts:te-3]); fret; };
			*|;
		subtype = "Bitmap" | "Text";
		format = "VOBSUB" | "UTF-8";
		subtitle := |*
			[A-Z] alpha+ - subtype - format => { fmt.Printf("a-%s:", data[ts:te]); };
			space;
			"(iso" digit{3} "-" digit ":" space lower{3} ")";
			space;
			"(" subtype ")" => { fmt.Printf("b-%s:", data[ts+1:te-1]); };
			"(" format ")" => { fmt.Printf("c-%s", data[ts+1:te-1]); fret; };
		*|;
		word = [a-z]+;
		prefix = space+ "+";
		prefixsp = prefix space;
		stream = prefixsp "stream:";
		duration = prefixsp "duration:";
		size = prefixsp "size:";
		pixelaspect = prefix any+ "pixel" space+ "aspect:";
		displayaspect = prefix any+ "display" space+ "aspect:";
		fps = prefix any+ "fps";
		autocrop = prefixsp "autocrop:";
		track = prefixsp digit "," space;
		main := ( 
			newline |
			stream @{ fcall stitle; } |
			duration @{ fcall sduration; } |
			size @{ fcall picture; } |
			pixelaspect @{ fcall paspect; } |
			displayaspect @{ fcall daspect; } |
			fps @{ fcall sfps; } |
			autocrop @{ fcall crop; } |
			prefixsp "chapters:" @{ section = CHAPTER; fmt.Printf("chapter"); } |
			prefixsp "audio tracks:" @{ section = AUDIO; fmt.Printf("audio"); } |
			prefixsp "subtitle tracks:" @{ section = SUBTITLE; fmt.Printf("subtitle"); } |
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
