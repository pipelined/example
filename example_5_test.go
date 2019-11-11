package example_test

import (
	"context"
	"log"
	"os"

	"github.com/pipelined/mixer"
	"github.com/pipelined/mp3"
	"github.com/pipelined/pipe"
	"github.com/pipelined/vst2"
	"github.com/pipelined/wav"
)

// This example demonstrates how to read two .wav files, mix,
// process them with vst2 plugin and save result as .mp3 file.
//
// NOTE: For simplicity both wav files have same characteristics i.e:
// sample rate and number of channels. In real life explicit
// conversion might be needed.
func Example_5() {
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
	mixer := mixer.New(2)

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

	// create a pipe with three lines.
	p, err := pipe.New(
		// line for mixing first wav file.
		&pipe.Line{
			// wav pump.
			Pump: &wav.Pump{ReadSeeker: inputFile1},
			// mixer sink.
			Sinks: pipe.Sinks(mixer),
		},
		// line for mixing second wav file.
		&pipe.Line{
			// wav pump.
			Pump: &wav.Pump{ReadSeeker: inputFile2},
			// mixer sink.
			Sinks: pipe.Sinks(mixer),
		},
		// line for sinking mp3.
		&pipe.Line{
			// mixer pump.
			Pump: mixer,
			// vst2 processor.
			Processors: pipe.Processors(&vst2.Processor{VST: vst}),
			Sinks: pipe.Sinks(
				// mp3 sink
				&mp3.Sink{
					Writer:      outputFile,
					BitRateMode: mp3.VBR(0),
					ChannelMode: mp3.JointStereo,
				},
			),
		},
	)
	if err != nil {
		log.Fatalf("failed to bind pipeline: %v", err)
	}
	defer p.Close()

	// run the pipeline.
	err = pipe.Wait(p.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %v", err)
	}
}
