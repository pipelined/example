package example

import (
	"context"
	"os"

	"github.com/pipelined/mixer"
	"github.com/pipelined/mp3"
	"github.com/pipelined/pipe"
	"github.com/pipelined/vst2"
	"github.com/pipelined/wav"
)

// Example5 demonstrates:
//	* Read signals from .wav files
//	* Mix signals with mixer
// 	* Process signal with VST2 plugin
//	* Save signal into .mp3 file
//
// NOTE: For example both wav files have same characteristics i.e: sample rate and number of channels.
// In real life implicit conversion will be needed.
func Example5() {
	bufferSize := 512

	// open first wav input
	inputFile1, err := os.Open("_testdata/sample1.wav")
	check(err)
	defer inputFile1.Close()

	// open second wav input
	inputFile2, err := os.Open("_testdata/sample2.wav")
	check(err)
	defer inputFile2.Close()

	// mixer
	mixer := mixer.New()

	// open vst library
	lib, err := vst2.Open("_testdata/Krush.vst")
	check(err)
	defer lib.Close()

	// open vst plugin
	plugin, err := lib.Open()
	check(err)
	defer plugin.Close()

	// create output file
	outputFile, err := os.Create("_testdata/out5.mp3")
	check(err)
	defer outputFile.Close()

	// create a line with three pipes
	l, err := pipe.Line(
		// pipe for mixing first wav file
		&pipe.Pipe{
			// wav pump
			Pump: &wav.Pump{ReadSeeker: inputFile1},
			// mixer sink
			Sinks: pipe.Sinks(mixer),
		},
		// pipe for mixing second wav file
		&pipe.Pipe{
			// wav pump
			Pump: &wav.Pump{ReadSeeker: inputFile2},
			// mixer sink
			Sinks: pipe.Sinks(mixer),
		},
		// pipe for sinking mp3
		&pipe.Pipe{
			// mixer pump
			Pump: mixer,
			// vst2 processor
			Processors: pipe.Processors(&vst2.Processor{Plugin: plugin}),
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
	check(err)
	defer l.Close()

	err = pipe.Wait(l.Run(context.Background(), bufferSize))
	check(err)
}
