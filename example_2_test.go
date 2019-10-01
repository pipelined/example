package example_test

import (
	"context"
	"log"
	"os"

	"github.com/pipelined/pipe"
	"github.com/pipelined/signal"
	"github.com/pipelined/vst2"
	"github.com/pipelined/wav"
)

// This example demonstrates how to process .wav file with
// vst2 plugin and write result to a new .wav file.
func Example_2() {
	// open input file.
	inputFile, err := os.Open("_testdata/sample1.wav")
	if err != nil {
		log.Fatalf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// open vst2 library.
	vst, err := vst2.Open("_testdata/Krush.vst")
	if err != nil {
		log.Fatalf("failed to open vst2 plugin: %w", err)
	}
	defer vst.Close()

	// create output file.
	outputFile, err := os.Create("_testdata/out2.wav")
	if err != nil {
		log.Fatalf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// build a line with single pipe.
	l, err := pipe.Line(
		&pipe.Pipe{
			// wav pump.
			Pump: &wav.Pump{
				ReadSeeker: inputFile,
			},
			// vst2 processor.
			Processors: pipe.Processors(
				&vst2.Processor{
					VST: vst,
				},
			),
			// wav sink.
			Sinks: pipe.Sinks(
				&wav.Sink{
					BitDepth:    signal.BitDepth16,
					WriteSeeker: outputFile,
				},
			),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind pipeline: %w", err)
	}
	defer l.Close()

	// run the pipeline.
	err = pipe.Wait(l.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %w", err)
	}
}
