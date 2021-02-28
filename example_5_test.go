package example_test

import (
	"context"
	"log"
	"os"

	"pipelined.dev/audio"
	"pipelined.dev/audio/mp3"
	"pipelined.dev/audio/vst2"
	"pipelined.dev/audio/wav"
	"pipelined.dev/pipe"
)

// This example demonstrates how to read two .wav files, mix,
// process them with vst2 plugin and save result as .mp3 file.
//
// NOTE: For simplicity both wav files have same characteristics i.e:
// sample rate and number of channels. In real life explicit
// conversion might be needed.
func Example_five() {
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

	// create new mixer with 2 channels.
	mixer := audio.Mixer{}

	// open vst library.
	vst, err := vst2.Open("_testdata/Krush.vst")
	if err != nil {
		log.Fatalf("failed to open vst2 plugin: %v", err)
	}
	defer vst.Close()

	// create output file.
	outputFile, err := os.Create("_testdata/out5.mp3")
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	defer outputFile.Close()

	bufferSize := 512
	p, err := pipe.New(bufferSize,
		// line for mixing first wav file.
		pipe.Line{
			// wav pump.
			Source: wav.Source(inputFile1),
			// mixer sink.
			Sink: mixer.Sink(),
		},
		// line for mixing second wav file.
		pipe.Line{
			// wav pump.
			Source: wav.Source(inputFile2),
			// mixer sink.
			Sink: mixer.Sink(),
		},
		// line for sinking mp3.
		pipe.Line{
			// mixer pump.
			Source: mixer.Source(),
			// vst2 processor.
			Processors: pipe.Processors(
				vst.Processor(vst2.Host{}).Allocator(nil),
			),
			Sink: mp3.Sink(
				outputFile,
				mp3.VBR(0),
				mp3.JointStereo,
				mp3.DefaultEncodingQuality,
			),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind lines: %v", err)
	}

	// execute the pipe with three lines.
	err = pipe.Wait(p.Start(context.Background()))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %v", err)
	}
}
