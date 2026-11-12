package mediainfo_test

import (
	"fmt"
	"log"

	"github.com/donghquinn/go-mediainfo"
)

func Example() {
	mi, err := mediainfo.Open("testdata/sample.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer mi.Close()

	format, err := mi.Get(mediainfo.Video, 0, "Format")
	if err != nil {
		log.Fatal(err)
	}
	width, err := mi.Get(mediainfo.Video, 0, "Width")
	if err != nil {
		log.Fatal(err)
	}
	height, err := mi.Get(mediainfo.Video, 0, "Height")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %sx%s\n", format, width, height)
	// Output: AVC 64x48
}

func ExampleMediaInfo_Inform() {
	mi, err := mediainfo.Open("testdata/sample.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer mi.Close()

	// Switch the report format to JSON before calling Inform.
	if _, err := mi.Option("Inform", "JSON"); err != nil {
		log.Fatal(err)
	}
	report, err := mi.Inform()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(len(report) > 0)
	// Output: true
}

func ExampleMediaInfo_Count() {
	mi, err := mediainfo.Open("testdata/sample.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer mi.Close()

	for _, kind := range []mediainfo.StreamKind{mediainfo.Video, mediainfo.Audio} {
		n, err := mi.Count(kind)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %d\n", kind, n)
	}
	// Output:
	// Video: 1
	// Audio: 1
}
