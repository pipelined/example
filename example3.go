package example

import (
	"context"
	"os"

	"github.com/pipelined/mixer"
	"github.com/pipelined/pipe"
	"github.com/pipelined/signal"
	"github.com/pipelined/wav"
)

// Example3 demonstrates:
//	* Read signals from .wav files
//	* Mix signals with mixer
//	* Save mixed signal into new .wav file
//
// NOTE: For example both wav files have same characteristics i.e: sample rate and number of channels.
// In real life implicit conversion will be needed.
func Example3() {
	bufferSize := 512
	// open first wav input
	inputFile1, err := os.Open("_testdata/sample1.wav")
	check(err)
	defer inputFile1.Close()

	// open second wav input
	inputFile2, err := os.Open("_testdata/sample2.wav")
	check(err)
	defer inputFile2.Close()

	// create output file
	outputFile, err := os.Create("_testdata/out3.wav")
	check(err)
	defer outputFile.Close()

	// create mixer
	mix := mixer.New()

	// create a line with three pipes
	l, err := pipe.Line(
		// pipe for first input
		&pipe.Pipe{
			// wav pump
			Pump: &wav.Pump{ReadSeeker: inputFile1},
			// mixer sink
			Sinks: pipe.Sinks(mix),
		},
		// pipe for second input
		&pipe.Pipe{
			// wav pump
			Pump: &wav.Pump{ReadSeeker: inputFile2},
			// mixer sink
			Sinks: pipe.Sinks(mix),
		},
		// pipe for output
		&pipe.Pipe{
			// mixer pump
			Pump: mix,
			// wav sink
			Sinks: pipe.Sinks(
				&wav.Sink{
					WriteSeeker: outputFile,
					BitDepth:    signal.BitDepth16,
				},
			),
		},
	)
	check(err)

	err = pipe.Wait(l.Run(context.Background(), bufferSize))
	check(err)
}
