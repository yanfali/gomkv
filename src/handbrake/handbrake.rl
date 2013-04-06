package handbrake

import (
	"fmt"
	"strings"
)

%%{
	machine handbrake;
	write data;
}%%

func parseOutput(data string) (HandBrakeMeta) {
	cs, p, pe, eof := 0, 0, len(data), 0
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
		action titleStart {
			 //fmt.Printf("%q\n", p)
			 fmt.Printf("S")
			 capture=true
			 csp = p
			}
		action titleEnd {
			 if capture {
			 meta.Title = strings.Trim(data[csp:p], " \n")
			 fmt.Printf("E%s", meta.Title)
			 capture = false
			 }
			  }
		newline = any* '\n' @ newline;
		prefix = space+ . "+" . space+;
		title = prefix . "stream:" . space+ >titleStart;
		titlee = title . any* . "\n" %titleEnd;
		main := ( 
			title |
			titlee |
			newline
		)*;
		write init;
		write exec;
	}%%
	return meta
}
