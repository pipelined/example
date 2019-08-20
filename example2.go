package example

import (
	"context"
	"os"

	"github.com/pipelined/pipe"
	"github.com/pipelined/signal"
	"github.com/pipelined/vst2"
	"github.com/pipelined/wav"
)

// Example2 demonstrates:
//	* Read signal from .wav file
//	* Process signal with VST2 plugin
// 	* Save signal into new .wav file
func Example2() {
	bufferSize := 512
	// open input file
	inputFile, err := os.Open("_testdata/sample1.wav")
	check(err)
	defer inputFile.Close()

	// open vst library
	lib, err := vst2.Open("_testdata/Krush.vst")
	check(err)
	defer lib.Close()
	// open vst plugin
	plugin, err := lib.Open()
	check(err)
	defer plugin.Close()

	// create output file
	outputFile, err := os.Create("_testdata/out2.wav")
	check(err)
	defer outputFile.Close()

	// build a line with single pipe
	l, err := pipe.Line(
		&pipe.Pipe{
			// wav pump
			Pump: &wav.Pump{ReadSeeker: inputFile},
			// vst2 processor
			Processors: pipe.Processors(&vst2.Processor{Plugin: plugin}),
			// wav sink
			Sinks: pipe.Sinks(&wav.Sink{BitDepth: signal.BitDepth16, WriteSeeker: outputFile}),
		},
	)
	check(err)
	defer l.Close()

	// run the line
	err = pipe.Wait(l.Run(context.Background(), bufferSize))
	check(err)
}
