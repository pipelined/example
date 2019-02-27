package example

import (
	"github.com/pipelined/mixer"
	"github.com/pipelined/pipe"
	"github.com/pipelined/signal"
	"github.com/pipelined/wav"
)

// Example3 demonstrates:
//	* Read signals from .wav files
//	* Mix signals with mixer
//	* Save mixed signal into new .wav file
//
// NOTE: For example both wav files have same characteristics i.e: sample rate and number of channels.
// In real life implicit conversion will be needed.
func Example3() {
	bufferSize := 512

	wavPump1 := wav.NewPump("_testdata/sample1.wav")

	wavPump2 := wav.NewPump("_testdata/sample2.wav")

	wavSink, err := wav.NewSink(
		"_testdata/out3.wav",
		signal.BitDepth16,
	)
	check(err)
	mixer := mixer.New()

	track1, err := pipe.New(
		bufferSize,
		pipe.WithPump(wavPump1),
		pipe.WithSinks(mixer),
	)
	check(err)
	defer track1.Close()
	track2, err := pipe.New(
		bufferSize,
		pipe.WithPump(wavPump2),
		pipe.WithSinks(mixer),
	)
	check(err)
	defer track2.Close()
	out, err := pipe.New(
		bufferSize,
		pipe.WithPump(mixer),
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
