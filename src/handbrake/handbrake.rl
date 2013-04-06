package handbrake

import (
)

%%{
	machine handbrake;
	write data;
}%%

func parseOutput(data string) (HandBrakeMeta) {
	cs, p, pe := 0, 0, len(data)
	%%{
		action stub {}
		main := ();
		write init;
		write exec;
	}%%
	return HandBrakeMeta{}
}
