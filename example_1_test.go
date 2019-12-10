package example_test

import (
	"context"
	"log"
	"os"

	"pipelined.dev/pipe"
	"pipelined.dev/portaudio"
	"pipelined.dev/wav"
)

// This example demonstrates how to play .wav file with portaudio.
func Example_1() {
	// open source wav file.
	wavFile, err := os.Open("_testdata/sample1.wav")
	if err != nil {
		log.Fatalf("failed to open wav file: %v", err)
	}
	defer wavFile.Close()

	// build pipe with a single line.
	p, err := pipe.New(
		&pipe.Line{
			// wav pump.
			Pump: &wav.Pump{
				ReadSeeker: wavFile,
			},
			// portaudio sink.
			Sinks: pipe.Sinks(&portaudio.Sink{}),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind pipeline: %v", err)
	}
	defer p.Close()

	// run the pipe.
	err = pipe.Wait(p.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %v", err)
	}
}
