package example_test

import (
	"context"
	"log"
	"os"

	"github.com/pipelined/pipe"
	"github.com/pipelined/portaudio"
	"github.com/pipelined/wav"
)

// This example demonstrates how to play .wav file with portaudio.
func Example_1() {
	// open source wav file.
	wavFile, err := os.Open("_testdata/sample1.wav")
	if err != nil {
		log.Fatalf("failed to open wav file: %w", err)
	}
	defer wavFile.Close()

	// build line with a single pipe.
	l, err := pipe.Line(
		&pipe.Pipe{
			// wav pump.
			Pump: &wav.Pump{
				ReadSeeker: wavFile,
			},
			// portaudio sink.
			Sinks: pipe.Sinks(&portaudio.Sink{}),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind pipeline: %w", err)
	}
	defer l.Close()

	// run the line.
	err = pipe.Wait(l.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %w", err)
	}
}
