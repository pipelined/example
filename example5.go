package example

import (
	"github.com/pipelined/mixer"
	"github.com/pipelined/pipe"
	"github.com/pipelined/signal"
	"github.com/pipelined/vst2"
	"github.com/pipelined/wav"
)

// Example:
//		Read two .wav files
//		Mix them
// 		Process with vst2
//		Save result into new .wav file
//
// NOTE: For example both wav files have same characteristics i.e: sample rate, bit depth and number of channels.
// In real life implicit conversion will be needed.
func five() {
	bufferSize := 512

	// wav pump 1
	wavPump1 := wav.NewPump("_testdata/sample1.wav")

	// wav pump 2
	wavPump2 := wav.NewPump("_testdata/sample2.wav")

	// mixer
	mixer := mixer.New()

	// track 1
	track1, err := pipe.New(
		bufferSize,
		pipe.WithPump(wavPump1),
		pipe.WithSinks(mixer),
	)
	check(err)
	defer track1.Close()
	// track 2
	track2, err := pipe.New(
		bufferSize,
		pipe.WithPump(wavPump2),
		pipe.WithSinks(mixer),
	)
	check(err)
	defer track2.Close()

	// vst2 processor
	vst2lib, err := vst2.Open("_testdata/Krush.vst")
	check(err)
	defer vst2lib.Close()

	vst2plugin, err := vst2lib.Open()
	check(err)
	defer vst2plugin.Close()

	vst2processor := vst2.NewProcessor(vst2plugin)

	// wav sink
	wavSink, err := wav.NewSink(
		"_testdata/out5.wav",
		signal.BitDepth16,
	)
	check(err)

	// out pipe
	out, err := pipe.New(
		bufferSize,
		pipe.WithPump(mixer),
		pipe.WithProcessors(vst2processor),
		pipe.WithSinks(wavSink),
	)
	check(err)
	defer out.Close()

	track1Errc := track1.Run()
	check(err)
	track2Errc := track2.Run()
	check(err)
	outErrc := out.Run()
	check(err)

	err = pipe.Wait(track1Errc)
	check(err)
	err = pipe.Wait(track2Errc)
	check(err)
	err = pipe.Wait(outErrc)
	check(err)
}
