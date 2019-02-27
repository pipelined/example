package example

import (
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
	wavPump := wav.NewPump("_testdata/sample1.wav")

	vst2lib, err := vst2.Open("_testdata/Krush.vst")
	check(err)
	defer vst2lib.Close()

	vst2plugin, err := vst2lib.Open()
	check(err)
	defer vst2plugin.Close()
	vst2processor := vst2.NewProcessor(
		vst2plugin,
	)
	wavSink, err := wav.NewSink(
		"_testdata/out2.wav",
		signal.BitDepth16,
	)
	check(err)
	p, err := pipe.New(
		bufferSize,
		pipe.WithPump(wavPump),
		pipe.WithProcessors(vst2processor),
		pipe.WithSinks(wavSink),
	)
	check(err)
	defer p.Close()
	err = pipe.Wait(p.Run())
	check(err)
}
