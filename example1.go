package example

import (
	"context"
	"os"

	"github.com/pipelined/pipe"
	"github.com/pipelined/portaudio"
	"github.com/pipelined/wav"
)

// Example1 demonstrates:
//	* Read signal from .wav file
//	* Play signal with portaudio
func Example1() {
	bufferSize := 512
	// open source file
	wavFile, err := os.Open("_testdata/sample1.wav")
	check(err)
	defer wavFile.Close()

	// build line with a single pipe
	l, err := pipe.Line(
		&pipe.Pipe{
			// wav pump
			Pump: &wav.Pump{ReadSeeker: wavFile},
			// portaudio sink
			Sinks: pipe.Sinks(&portaudio.Sink{}),
		},
	)
	check(err)
	defer l.Close()

	// run the line
	err = pipe.Wait(l.Run(context.Background(), bufferSize))
	check(err)
}
