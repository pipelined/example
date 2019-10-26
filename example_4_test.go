package example_test

import (
	"context"
	"log"
	"os"

	"github.com/pipelined/audio"
	"github.com/pipelined/pipe"
	"github.com/pipelined/portaudio"
	"github.com/pipelined/signal"
	"github.com/pipelined/wav"
)

// This example demonstrates how to read signal from .wav file,
// compose a track with clips from that file and then sumalteniously
// save it to new .wav file and play it with portaudio.
func Example_4() {
	// open input file.
	inputFile, err := os.Open("_testdata/sample1.wav")
	if err != nil {
		log.Fatalf("failed to open input file: %w", err)
	}
	defer inputFile.Close()

	// asset sink.
	asset := &audio.Asset{}

	// read wav pipeline.
	wavFile, err := pipe.New(
		&pipe.Line{
			// wav pump.
			Pump: &wav.Pump{ReadSeeker: inputFile},
			// in-memory asset.
			Sinks: pipe.Sinks(asset),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind import pipeline: %w", err)
	}
	defer wavFile.Close()

	err = pipe.Wait(wavFile.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute import pipeline: %w", err)
	}

	// track pump.
	track := audio.NewTrack(asset.SampleRate(), asset.NumChannels())

	// add samples.
	track.AddClip(198450, asset.Clip(0, 44100))
	track.AddClip(66150, asset.Clip(44100, 44100))
	track.AddClip(132300, asset.Clip(0, 44100))

	// create output file.
	outputFile, err := os.Create("_testdata/out4.wav")
	if err != nil {
		log.Fatalf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// pipeline to process clips.
	p, err := pipe.New(
		&pipe.Line{
			// track with clips.
			Pump: track,
			Sinks: pipe.Sinks(
				// wav sink.
				&wav.Sink{
					WriteSeeker: outputFile,
					BitDepth:    signal.BitDepth16,
				},
				// portaudio sink.
				&portaudio.Sink{},
			),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind playback and save pipeline: %w", err)
	}
	defer p.Close()

	// run the pipeline.
	err = pipe.Wait(p.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute playback and save pipeline: %w", err)
	}
}
