package example

import (
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

	// wav pump
	wavPump := wav.NewPump("_testdata/sample1.wav")

	// asset sink
	asset := &audio.Asset{}

	// import pipe
	importAsset, err := pipe.New(
		bufferSize,
		pipe.WithPump(wavPump),
		pipe.WithSinks(asset),
	)
	check(err)
	defer importAsset.Close()

	err = pipe.Wait(importAsset.Run())
	check(err)

	// track pump
	track := audio.NewTrack(44100, asset.NumChannels())

	// add samples
	track.AddClip(198450, asset.Clip(0, 44100))
	track.AddClip(66150, asset.Clip(44100, 44100))
	track.AddClip(132300, asset.Clip(0, 44100))

	// wav sink
	wavSink, err := wav.NewSink(
		"_testdata/out4.wav",
		signal.BitDepth16,
	)
	// portaudio sink
	paSink := portaudio.NewSink()

	// final pipe
	p, err := pipe.New(
		bufferSize,
		pipe.WithPump(track),
		pipe.WithSinks(wavSink, paSink),
	)
	check(err)

	err = pipe.Wait(p.Run())
	check(err)
}
