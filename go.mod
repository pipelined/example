module github.com/pipelined/example

go 1.12

require (
	github.com/kr/pty v1.1.8 // indirect
	github.com/pipelined/audio v0.0.0-20190816061705-53ca487f40df
	github.com/pipelined/mixer v0.1.0
	github.com/pipelined/mp3 v0.2.1
	github.com/pipelined/pipe v0.5.2
	github.com/pipelined/portaudio v0.0.0-20190820055226-d77195b43cc4
	github.com/pipelined/signal v0.2.0
	github.com/pipelined/vst2 v0.4.0
	github.com/pipelined/wav v0.2.0
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/viert/lame v0.0.0-20190823071122-49a063e7d5e6 // indirect
)

replace github.com/pipelined/audio => ../audio
