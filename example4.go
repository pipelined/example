package example

import (
	"context"
	"os"

	"github.com/pipelined/audio"
	"github.com/pipelined/pipe"
	"github.com/pipelined/portaudio"
	"github.com/pipelined/signal"
	"github.com/pipelined/wav"
)

// Example4 demonstrates:
//	* Read signal from .wav file
// 	* Slice signal into Clips
// 	* Put Clips to Track
//	* Save track signal into .wav and play it with portaudio at the same time
func Example4() {
	bufferSize := 512
	// open input file
	inputFile, err := os.Open("_testdata/sample1.wav")
	check(err)
	defer inputFile.Close()

	// asset sink
	asset := &audio.Asset{}

	// import line
	importAsset, err := pipe.Line(
		&pipe.Pipe{
			// wav pump
			Pump: &wav.Pump{ReadSeeker: inputFile},
			// in-memory asset
			Sinks: pipe.Sinks(asset),
		},
	)
	check(err)
	defer importAsset.Close()

	err = pipe.Wait(importAsset.Run(context.Background(), bufferSize))
	check(err)

	// track pump
	track := audio.NewTrack(44100, asset.NumChannels())

	// add samples
	track.AddClip(198450, asset.Clip(0, 44100))
	track.AddClip(66150, asset.Clip(44100, 44100))
	track.AddClip(132300, asset.Clip(0, 44100))

	// create output file
	outputFile, err := os.Create("_testdata/out4.wav")
	check(err)
	defer outputFile.Close()

	// line to process clips
	l, err := pipe.Line(
		&pipe.Pipe{
			// track with clips
			Pump: track,
			Sinks: pipe.Sinks(
				// wav sink
				&wav.Sink{
					WriteSeeker: outputFile,
					BitDepth:    signal.BitDepth16,
				},
				// portaudio sink
				&portaudio.Sink{},
			),
		},
	)
	check(err)

	err = pipe.Wait(l.Run(context.Background(), bufferSize))
	check(err)
}
