package example_test

import (
	"context"
	"log"
	"os"
	"time"

	"pipelined.dev/audio"
	"pipelined.dev/audio/mp3"
	"pipelined.dev/pipe"
	"pipelined.dev/signal"
)

// This example demonstrates how to cut a clip from .mp3 file
// and save the result to new .mp3 file.
func Example_six() {
	// open source wav file.
	mp3File, err := os.Open("_testdata/sample.mp3")
	if err != nil {
		log.Fatalf("failed to open mp3 file: %v", err)
	}
	defer mp3File.Close()

	// create output file.
	outputFile, err := os.Create("_testdata/out6.mp3")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	// asset to keep mp3 signal.
	a := &audio.Asset{}

	bufferSize := 512
	p, err := pipe.New(
		bufferSize,
		pipe.Line{
			// mp3 pump.
			Source: mp3.Source(mp3File),
			// asset sink.
			Sink: a.Sink(),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind asset pipeline: %v", err)
	}

	// execute pipe with a single line.
	err = pipe.Wait(p.Start(context.Background()))
	if err != nil {
		log.Fatalf("failed to execute asset pipeline: %v", err)
	}

	// get the sample rate of the asset.
	// it will be needed to convert duration.
	sampleRate := a.SampleRate()
	// cut the clip that starts at 1 second and lasts 2.5 seconds.
	clip := signal.Slice(
		a.Signal,
		sampleRate.Events(time.Millisecond*1000),
		sampleRate.Events(time.Millisecond*3500),
	)

	p, err = pipe.New(
		bufferSize,
		pipe.Line{
			// clip source
			Source: audio.Source(sampleRate, clip),
			// mp3 sink
			Sink: mp3.Sink(
				outputFile,
				mp3.CBR(320),
				mp3.JointStereo,
				mp3.DefaultEncodingQuality,
			),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind output pipeline: %v", err)
	}
	// build pipe with a single line.
	err = pipe.Wait(p.Start(context.Background()))
	if err != nil {
		log.Fatalf("failed to execute output pipeline: %v", err)
	}
}
