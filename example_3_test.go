package example_test

import (
	"context"
	"log"
	"os"

	"pipelined.dev/audio"
	"pipelined.dev/pipe"
	"pipelined.dev/signal"
	"pipelined.dev/wav"
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
		log.Fatalf("failed to open first input file: %v", err)
	}
	defer inputFile1.Close()

	// open second wav input.
	inputFile2, err := os.Open("_testdata/sample2.wav")
	if err != nil {
		log.Fatalf("failed to open second input file: %v", err)
	}
	defer inputFile2.Close()

	// create output file.
	outputFile, err := os.Create("_testdata/out3.wav")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// create mixer with 2 channels.
	mix := audio.NewMixer(2)

	// create a pipe with three pipes.
	p, err := pipe.New(
		// line for first input.
		&pipe.Line{
			// wav pump.
			Pump: &wav.Pump{ReadSeeker: inputFile1},
			// mixer sink.
			Sinks: pipe.Sinks(mix),
		},
		// line for second input
		&pipe.Line{
			// wav pump.
			Pump: &wav.Pump{ReadSeeker: inputFile2},
			// mixer sink.
			Sinks: pipe.Sinks(mix),
		},
		// line for output.
		&pipe.Line{
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
		log.Fatalf("failed to bind pipeline: %v", err)
	}
	defer p.Close()

	// run the pipeline.
	err = pipe.Wait(p.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %v", err)
	}
}
