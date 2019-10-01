package example_test

import (
	"context"
	"log"
	"os"

	"github.com/pipelined/mixer"
	"github.com/pipelined/pipe"
	"github.com/pipelined/signal"
	"github.com/pipelined/wav"
)

// This example demonstrates how to read two .wav files,
// mix them together and save result to new .wav file.
//
// NOTE: For simplicity both wav files have same characteristics i.e:
// sample rate and number of channels. In real life explicit
// conversion might be needed.
func Example_3() {
	// open first wav input.
	inputFile1, err := os.Open("_testdata/sample1.wav")
	if err != nil {
		log.Fatalf("failed to open first input file: %w", err)
	}
	defer inputFile1.Close()

	// open second wav input.
	inputFile2, err := os.Open("_testdata/sample2.wav")
	if err != nil {
		log.Fatalf("failed to open second input file: %w", err)
	}
	defer inputFile2.Close()

	// create output file.
	outputFile, err := os.Create("_testdata/out3.wav")
	if err != nil {
		log.Fatalf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// create mixer.
	mix := mixer.New()

	// create a line with three pipes.
	l, err := pipe.Line(
		// pipe for first input.
		&pipe.Pipe{
			// wav pump.
			Pump: &wav.Pump{ReadSeeker: inputFile1},
			// mixer sink.
			Sinks: pipe.Sinks(mix),
		},
		// pipe for second input
		&pipe.Pipe{
			// wav pump.
			Pump: &wav.Pump{ReadSeeker: inputFile2},
			// mixer sink.
			Sinks: pipe.Sinks(mix),
		},
		// pipe for output.
		&pipe.Pipe{
			// mixer pump.
			Pump: mix,
			// wav sink.
			Sinks: pipe.Sinks(
				&wav.Sink{
					WriteSeeker: outputFile,
					BitDepth:    signal.BitDepth16,
				},
			),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind pipeline: %w", err)
	}

	// run the pipeline.
	err = pipe.Wait(l.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %w", err)
	}
}
