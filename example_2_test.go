package example_test

import (
	"context"
	"log"
	"os"

	"pipelined.dev/audio/vst2"
	"pipelined.dev/audio/wav"
	"pipelined.dev/pipe"
	"pipelined.dev/signal"
)

// This example demonstrates how to process .wav file with
// vst2 plugin and write result to a new .wav file.
func Example_two() {
	// open input file.
	inputFile, err := os.Open("_testdata/sample1.wav")
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	// open vst2 library.
	vst, err := vst2.Open("_testdata/Krush.vst")
	if err != nil {
		log.Fatalf("failed to open vst2 plugin: %v", err)
	}
	defer vst.Close()

	// create output file.
	outputFile, err := os.Create("_testdata/out2.wav")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	bufferSize := 512
	// build the line.
	p, err := pipe.New(bufferSize, pipe.Line{
		// wav pump.
		Source: wav.Source(inputFile),
		// vst2 processor.
		Processors: pipe.Processors(
			vst.Processor(vst2.Host{}).Allocator(nil),
		),
		// wav sink.
		Sink: wav.Sink(outputFile, signal.BitDepth16),
	})
	if err != nil {
		log.Fatalf("failed to bind line: %v", err)
	}

	// run the pipe with a single line.
	err = pipe.Wait(p.Start(context.Background()))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %v", err)
	}
}
