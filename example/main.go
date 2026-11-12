// Command example prints a single MediaInfo parameter of a media file.
//
//	go run ./example movie.mp4 Video Width
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/donghquinn/go-mediainfo"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "usage: example <file> <General|Video|Audio|Text|Other|Image|Menu> <parameter>")
		os.Exit(2)
	}

	kind, ok := parseKind(os.Args[2])
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown stream kind %q\n", os.Args[2])
		os.Exit(2)
	}

	mi, err := mediainfo.Open(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer mi.Close()

	val, err := mi.Get(kind, 0, os.Args[3])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(val)
}

func parseKind(s string) (mediainfo.StreamKind, bool) {
	for k := mediainfo.General; k <= mediainfo.Menu; k++ {
		if strings.EqualFold(k.String(), s) {
			return k, true
		}
	}
	return 0, false
}
