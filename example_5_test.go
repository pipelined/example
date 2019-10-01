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
		log.Fatalf("failed to open first input file: %w", err)
	}
	defer inputFile1.Close()

	// open second wav input.
	inputFile2, err := os.Open("_testdata/sample2.wav")
	if err != nil {
		log.Fatalf("failed to open second input file: %w", err)
	}
	defer inputFile2.Close()

	// create new mixer.
	mixer := mixer.New()

	// open vst library.
	vst, err := vst2.Open("_testdata/Krush.vst")
	if err != nil {
		log.Fatalf("failed to open vst2 plugin: %w", err)
	}
	defer vst.Close()

	// create output file.
	outputFile, err := os.Create("_testdata/out5.mp3")
	if err != nil {
		log.Fatalf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	// create a line with three pipes.
	l, err := pipe.Line(
		// pipe for mixing first wav file.
		&pipe.Pipe{
			// wav pump.
			Pump: &wav.Pump{ReadSeeker: inputFile1},
			// mixer sink.
			Sinks: pipe.Sinks(mixer),
		},
		// pipe for mixing second wav file.
		&pipe.Pipe{
			// wav pump.
			Pump: &wav.Pump{ReadSeeker: inputFile2},
			// mixer sink.
			Sinks: pipe.Sinks(mixer),
		},
		// pipe for sinking mp3.
		&pipe.Pipe{
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
		log.Fatalf("failed to bind pipeline: %w", err)
	}
	defer l.Close()

	// run the pipeline.
	err = pipe.Wait(l.Run(context.Background(), 512))
	if err != nil {
		log.Fatalf("failed to execute pipeline: %w", err)
	}
}
