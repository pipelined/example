module pipelined.dev/example

go 1.13

require (
	pipelined.dev/audio v0.2.2-0.20191204072949-aab07b1e55dd
	pipelined.dev/audio/mp3 v0.5.0
	pipelined.dev/audio/portaudio v0.3.0
	pipelined.dev/audio/vst2 v0.6.1
	pipelined.dev/audio/wav v0.5.0
	pipelined.dev/pipe v0.8.4
	pipelined.dev/signal v0.8.0
)

replace (
	pipelined.dev/audio => ../audio
	pipelined.dev/audio/vst2 => ../vst2
)
