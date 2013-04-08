package handbrake

import (
	"fmt"
)

%%{
	machine handbrake;
	write data;
}%%

func parseOutput(data string) (HandBrakeMeta) {
	cs, p, pe, eof := 0, 0, len(data), 0
	top, ts, te, act := 0,0,0,0
	_,_,_,_ = top, ts, te, act
	stack := []int{0}
	_ = eof
	line := 1
	capture := false
	csp := 0
	meta := HandBrakeMeta{}
	fmt.Printf("%02d: ", line)
	_ = capture
	_ = csp
	%%{
		action newline { line +=1; fmt.Printf("\n%02d: ", line) }
		newline = any* '\n' @ newline;
		title := |*
			(alnum|space)+[.]*alnum* => { fmt.Printf("%s", data[ts:te]); fret;};
			"\n" => { fret; };
		*|;
		duration := |*
			space*;
			digit{2}[:]digit{2}[:]digit{2} => { fmt.Printf("%s", data[ts:te]); fret;};
			"\n" => { fret; };
		*|;
		picture := |*
			space*;
			digit{3,4} "x" digit{3,4} => { fmt.Printf("%s", data[ts:te]); fret;};
		*|;
		paspect := |*
			space*;
			digit{1,4} "/" digit{1,4} => { fmt.Printf("%s", data[ts:te]); fret; };
		*|;
		daspect := |*
			space*;
			digit{1} . "." . digit{1,3} => { fmt.Printf("%s", data[ts:te]); fret; };
		*|;
		fps := |*
			"\n" => { ts -= 10; fmt.Printf("%s", data[ts:te-5]); p -= 1; fret; };
		*|;
		crop := |*
			space*;
			digit{1,3} "/" digit{1,3} "/" digit{1,3} "/" digit{1,3} => { fmt.Printf("%s", data[ts:te]); fret; };
		*|;
		main := ( 
			newline |
			space+ "+" space+ "stream:" space* @{ fcall title; } |
			space+ "+" space+ "duration:" @{ fcall duration; } |
			space+ "+" space+ "size:" @{ fcall picture; } |
			space+ "+" any+ "pixel" space+ "aspect:" @{ fcall paspect; } |
			space+ "+" any+ "display" space+ "aspect:" @{ fcall daspect; } |
			space+ "+" any+ "fps" @{ fcall fps; } |
			space+ "+" space+ "autocrop:" @{ fcall crop; }
		)*;
		write init;
		write exec;
	}%%
	return meta
}
