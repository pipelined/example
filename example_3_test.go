package example_test

import (
	"context"
	"log"
	"os"

	"pipelined.dev/audio"
	"pipelined.dev/audio/wav"
	"pipelined.dev/pipe"
	"pipelined.dev/signal"
)

// This example demonstrates how to read two .wav files,
// mix them together and save result to new .wav file.
//
// NOTE: For simplicity both wav files have same characteristics i.e:
// sample rate and number of channels. In real life explicit
// conversion might be needed.
func Example_three() {
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
	mix := audio.Mixer{}

	bufferSize := 512
	p, err := pipe.New(
		bufferSize,
		pipe.Line{
			Source: wav.Source(inputFile1),
			Sink:   mix.Sink(),
		},
		pipe.Line{
			Source: wav.Source(inputFile2),
			Sink:   mix.Sink(),
		},
		pipe.Line{
			Source: mix.Source(),
			Sink:   wav.Sink(outputFile, signal.BitDepth16),
		},
	)

	// run the pipeline.
	err = pipe.Wait(p.Start(context.Background()))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %v", err)
	}
}
