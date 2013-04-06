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
			space*;
			alnum+ space* alnum+[.]"mkv" => { fmt.Printf("%s", data[ts:te]); fret;};
			"\n" => { fret; };
		*|;
		duration := |*
			space*;
			digit{2}[:]digit{2}[:]digit{2} => { fmt.Printf("%s", data[ts:te]); fret;};
			"\n" => { fret; };
		*|;
		main := ( 
			newline |
			space+ "+" space+ . "stream:" @{ fcall title; } |
			space+ "+" space+ . "duration:" @{ fcall duration; }
		)*;
		write init;
		write exec;
	}%%
	return meta
}
