# go-mediainfo

[![CI](https://github.com/donghquinn/go-mediainfo/actions/workflows/ci.yml/badge.svg)](https://github.com/donghquinn/go-mediainfo/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/donghquinn/go-mediainfo.svg)](https://pkg.go.dev/github.com/donghquinn/go-mediainfo)

Go bindings for [MediaInfo](https://mediaarea.net/en/MediaInfo), for reading technical metadata
(codec, duration, resolution, bitrate, sample rate, ...) from audio and video files.

## Requirements

This package uses cgo. At build time it needs the MediaInfo headers (`MediaInfoDLL/MediaInfoDLL.h`). At run time it loads the `libmediainfo` shared library with `dlopen`.

| Platform      | Install                                         |
| ------------- | ----------------------------------------------- |
| Debian/Ubuntu | `sudo apt-get install libmediainfo-dev`         |
| Fedora/RHEL   | `sudo dnf install libmediainfo-devel`           |
| Alpine        | `apk add libmediainfo-dev`                      |
| macOS         | `brew install libmediainfo`                     |

On Apple Silicon, Homebrew installs into `/opt/homebrew`, which `dlopen` does not search by
default. Tell the dynamic loader where to look:

```sh
export DYLD_FALLBACK_LIBRARY_PATH="$(brew --prefix)/lib"
```

Homebrew headers in `/opt/homebrew/include` (arm64) and `/usr/local/include` (amd64) are found
automatically. For other locations, set `CGO_CFLAGS=-I/path/to/include`.

## Installation

```sh
go get github.com/donghquinn/go-mediainfo
```

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/donghquinn/go-mediainfo"
)

func main() {
	mi, err := mediainfo.Open("movie.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer mi.Close()

	width, err := mi.Get(mediainfo.Video, 0, "Width")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("width:", width)

	audioTracks, _ := mi.Count(mediainfo.Audio)
	fmt.Println("audio tracks:", audioTracks)

	// Full report as JSON.
	mi.Option("Inform", "JSON")
	report, _ := mi.Inform()
	fmt.Println(report)
}
```

Run `mediainfo --Info-Parameters` to list every parameter name you can pass to `Get`.

A runnable command lives in [`example/`](example/main.go):

```sh
go run ./example movie.mp4 Video Width
```

## API

| Function / method                         | Description                                                     |
| ----------------------------------------- | --------------------------------------------------------------- |
| `Init() error`                            | Load the shared library early (optional; `Open` does it lazily) |
| `Version() (string, error)`               | Version of the loaded MediaInfo library                         |
| `Open(path) (*MediaInfo, error)`          | Open and parse a file                                           |
| `(*MediaInfo).Get(kind, stream, key)`     | Value of one parameter for the n-th stream of a kind            |
| `(*MediaInfo).Count(kind)`                | Number of streams of a kind                                     |
| `(*MediaInfo).Option(name, value)`        | Set a MediaInfo option, e.g. `Option("Inform", "JSON")`         |
| `(*MediaInfo).Inform()`                   | Full report in the format chosen by the `Inform` option         |
| `(*MediaInfo).Close() error`              | Release the handle (idempotent)                                 |

Stream kinds: `General`, `Video`, `Audio`, `Text`, `Other`, `Image`, `Menu`.

Errors can be matched with `errors.Is`: `ErrLoad`, `ErrOpen`, `ErrNotFound`, `ErrClosed`.

A `*MediaInfo` is safe for concurrent use.

## Development

```sh
go vet ./...
go test -race ./...
golangci-lint run
```

Releases are cut by pushing a semver tag. CI runs first, then GoReleaser creates the GitHub
Release, and the version is published to the Go module proxy:

```sh
git tag v1.0.0
git push origin v1.0.0
```

## License

MIT. See [LICENSE](LICENSE).

Based on the original [dwbuiten/go-mediainfo](https://github.com/dwbuiten/go-mediainfo) by
Derek Buitenhuis, distributed under the BSD 2-Clause license. Its notice is kept in
[NOTICE](NOTICE).
