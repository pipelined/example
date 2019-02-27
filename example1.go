package example

import (
	"github.com/pipelined/pipe"
	"github.com/pipelined/portaudio"
	"github.com/pipelined/wav"
)

// Example:
//		Read .wav file
//		Play it with portaudio
func one() {
	bufferSize := 512
	// wav pump
	wavPump := wav.NewPump("_testdata/sample1.wav")

	// portaudio sink
	paSink := portaudio.NewSink()

	// build pipe
	p, err := pipe.New(
		bufferSize,
		pipe.WithPump(wavPump),
		pipe.WithSinks(paSink),
	)
	check(err)
	defer p.Close()

	// run pipe
	err = pipe.Wait(p.Run())
	check(err)
}
