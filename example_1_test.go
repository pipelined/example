package example_test

import (
	"context"
	"log"
	"os"

	"pipelined.dev/audio/portaudio"
	"pipelined.dev/audio/wav"
	"pipelined.dev/pipe"
)

// This example demonstrates how to play .wav file with portaudio.
func Example_1() {
	bufferSize := 512
	// open source wav file.
	wavFile, err := os.Open("_testdata/sample1.wav")
	if err != nil {
		log.Fatalf("failed to open wav file: %v", err)
	}
	defer wavFile.Close()

	err = portaudio.Initialize()
	if err != nil {
		log.Fatalf("failed to initialize portaudio: %v", err)
	}
	defer portaudio.Terminate()

	device, err := portaudio.DefaultOutputDevice()
	if err != nil {
		log.Fatalf("failed to get default system device: %v", err)
	}
	// build pipe with a single line.
	l, err := pipe.Routing{
		Source: wav.Source(wavFile),
		Sink:   portaudio.Sink(device),
	}.Line(bufferSize)
	if err != nil {
		log.Fatalf("failed to bind line: %v", err)
	}

	err = pipe.New(
		context.Background(),
		pipe.WithLines(l),
	).Wait()
	if err != nil {
		log.Fatalf("failed to execute pipeline: %v", err)
	}
}
