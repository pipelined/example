package example_test

import (
	"context"
	"log"
	"os"
	"time"

	"pipelined.dev/audio"
	"pipelined.dev/audio/portaudio"
	"pipelined.dev/audio/wav"
	"pipelined.dev/pipe"
	"pipelined.dev/signal"
)

// This example demonstrates how to read signal from .wav file,
// compose a track with clips from that file and then sumalteniously
// save it to new .wav file and play it with portaudio.
func Example_4() {
	// open input file.
	inputFile, err := os.Open("_testdata/sample1.wav")
	if err != nil {
		log.Fatalf("failed to open input file: %v", err)
	}
	defer inputFile.Close()

	bufferSize := 512
	// asset sink.
	asset := &audio.Asset{}

	// read wav line.
	l, err := pipe.Routing{
		Source: wav.Source(inputFile),
		Sink:   asset.Sink(),
	}.Line(bufferSize)
	if err != nil {
		log.Fatalf("failed to bind import pipeline: %v", err)
	}

	err = pipe.New(context.Background(), pipe.WithLines(l)).Wait()
	if err != nil {
		log.Fatalf("failed to execute import pipeline: %v", err)
	}

	sampleRate := asset.SampleRate()
	// track pump.
	track := audio.Track{
		SampleRate: asset.SampleRate(),
		Channels:   asset.Channels(),
	}

	// add samples.
	track.AddClip(198450, signal.Slice(asset.Signal, 0, sampleRate.SamplesIn(1*time.Second)))
	track.AddClip(66150, signal.Slice(asset.Signal, sampleRate.SamplesIn(1*time.Second), sampleRate.SamplesIn(2*time.Second)))
	track.AddClip(132300, signal.Slice(asset.Signal, 0, sampleRate.SamplesIn(1*time.Second)))

	repeater := audio.Repeater{}

	// create output file.
	outputFile, err := os.Create("_testdata/out4.wav")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	err = portaudio.Initialize()
	if err != nil {
		log.Fatalf("failed to init portaudio: %v", err)
	}
	defer portaudio.Terminate()
	device, err := portaudio.DefaultInputDevice()
	if err != nil {
		log.Fatalf("failed to get default system device: %v", err)
	}

	lines, err := pipe.Lines(
		bufferSize,
		pipe.Routing{
			Source: track.Source(0, 0),
			Sink:   repeater.Sink(),
		},
		pipe.Routing{
			Source: repeater.Source(),
			Sink:   wav.Sink(outputFile, signal.BitDepth16),
		},
		pipe.Routing{
			Source: repeater.Source(),
			Sink:   portaudio.Sink(device),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind lines: %v", err)
	}

	// execute the pipe with three lines.
	err = pipe.New(context.Background(), pipe.WithLines(lines...)).Wait()
	if err != nil {
		log.Fatalf("failed to execute playback and save pipeline: %v", err)
	}
}
