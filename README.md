gomkv automating the tedious parts of Transcoding
-------------------------------------------------

gomkv is the result of my having encoded lots of lots of video using HandBrakeCLI. Nothing gomkv does is hard or complicated, it just takes my preferences and makes them the default while automating one of the most tedious aspects of converting mkvs, parsing the output of ```-t0```

gomkv is not the simplest way to do this. I also wanted to play more with golang and ragel. I use ragel to parse the output of t0 as a simple state machine. I could have done this with regex but that would have been too easy and I wouldn't have learned very much.

Requires Ragel v6.8.

Planning of what it's going to do, is going on [here](https://www.evernote.com/shard/s28/sh/7e79e2e8-925e-4aec-8852-a71954d63040/2327eb4a7245582ddc6822f5d5b1be8a)
