gomkv automating the tedious parts of Transcoding
=================================================

gomkv is the result of my having encoded lots of lots of video using HandBrakeCLI. Nothing gomkv does is hard or complicated, it just takes my preferences and makes them the default while automating one of the most tedious aspects of converting mkvs, parsing the output of ```-t0```

The output is an executable shell script which you can then run in batch.

gomkv is not the simplest way to do this. I also wanted to play more with golang and ragel. I use ragel to parse the output of t0 as a simple state machine. I could have done this with regex but that would have been too easy and I wouldn't have learned very much.

Requires Ragel v6.8.

Planning of what it's going to do, is going on [here](https://www.evernote.com/shard/s28/sh/7e79e2e8-925e-4aec-8852-a71954d63040/2327eb4a7245582ddc6822f5d5b1be8a)

Building
--------
		cd src/gomkv/handbrake && ragel -Z handbrake.rl
		cd - && go install gomkv

Usage
-----
	gomkv
		-aac=false: Encode audio using aac, instead of copying
		-debug=0: Debug level 1..3
		-dest-dir="": directory you want video files to be created
		-episode=1: Episode starting offset.
		-languages="": list of languages and order to copy, comma separated e.g. English,Japanese
		-mobile=false: Use mobile friendly settings
		-prefix="": Default Prefix for output filename(s)
		-profile="High Profile": Default Encoding Profile. Defaults to 'High Profile'
		-season=1: Season starting offset.
		-series=false: Videos are episodes of a series
		-source-dir="": directory containing video files. Defaults to current working directory.
		-split-chapters=0: Create one file for every N chapters. Only works with --series. e.g. -split-chapters 5
		-subs=true: Copy subtitles
		-subtitle-default="": Enable subtitles by default for the language matching this value. e.g. -subtitle-default=English
	

Examples
--------

1. Simple:

		gomkv --source-dir=/my/videos --prefix NEW_VIDEOS

2. TV Series:

		gomkv --source-dir=/my/tvseries --prefix TV_SERIES --series --dest-dir=/tmp --season=2 --episode=4

3. Mobile Friendly:

		gomkv --source-dir=/my/videos --dest-dir=/my/mobile --mobile

4. Japanese Then English:

		gomkv --source-dir=/my/videos --dest-dir=/tmp --languages=Japanese,English

5. Japanese Then English and Subtitles:

		gomkv --source-dir=/my/videos --dest-dir=/tmp --languages=Japanese,English --subtitle-default=English

6. Split a single file into many files based on chapter count:

		gomkv --source-dir=/my/videos --dest-dir=/tmp --series --split-chapters=5 --prefix TV_SERIES

Example output:
---------------

		gomkv --source-dir /my/videos --languages="Japanese,English" --prefix SF --series --episode=8


		HandBrakeCLI -Z "High Profile" -i /my/videos/title00.mkv -t1 -a2,1 -E copy:ac3,copy:ac3 -s 1,2 -o /home/yanfali/SF_S1E08.mkv
		HandBrakeCLI -Z "High Profile" -i /my/videos/title01.mkv -t1 -a2,1 -E copy:ac3,copy:ac3 -s 1,2 -o /home/yanfali/SF_S1E09.mkv
		HandBrakeCLI -Z "High Profile" -i /my/videos/title02.mkv -t1 -a2,1 -E copy:ac3,copy:ac3 -s 1,2 -o /home/yanfali/SF_S1E10.mkv
		HandBrakeCLI -Z "High Profile" -i /my/videos/title03.mkv -t1 -a2,1 -E copy:ac3,copy:ac3 -s 1,2 -o /home/yanfali/SF_S1E11.mkv
		HandBrakeCLI -Z "High Profile" -i /my/videos/title04.mkv -t1 -a2,1 -E copy:ac3,copy:ac3 -s 1,2 -o /home/yanfali/SF_S1E12.mkv
		HandBrakeCLI -Z "High Profile" -i /my/videos/title05.mkv -t1 -a2,1 -E copy:ac3,copy:ac3 -s 1,2 -o /home/yanfali/SF_S1E13.mkv

