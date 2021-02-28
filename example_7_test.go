package example_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"pipelined.dev/audio/fileformat"
	"pipelined.dev/audio/mp3"
	"pipelined.dev/audio/wav"
	"pipelined.dev/pipe"
	"pipelined.dev/signal"
)

// This example demonstrates how to walk file system and process files
// depending on audio file format. It converts all mp3 files to 16-bit wav
// and removes all flac files.
func Example_seven() {
	err := filepath.Walk("_testdata", fileformat.Walk(
		func(f *fileformat.Format, path string, fi os.FileInfo) error {
			switch f {
			case fileformat.MP3():
				// open mp3 input file
				input, err := os.Open(path)
				if err != nil {
					return err
				}
				defer input.Close()

				// create wav output file
				output, err := os.Create(fmt.Sprintf("%s.wav", path))
				defer output.Close()

				// bind the line
				p, err := pipe.New(
					512,
					pipe.Line{
						Source: mp3.Source(input),
						Sink:   wav.Sink(output, signal.BitDepth16),
					},
				)
				if err != nil {
					return err
				}

				// execute the pipe with single line
				return pipe.Wait(p.Start(context.Background()))
			case fileformat.FLAC():
				// remove flac file
				return os.Remove(path)
			}
			return nil
		}, true))
	if err != nil {
		log.Fatalf("failed to walk: %v", err)
	}
}
