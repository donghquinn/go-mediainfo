// Package mediainfo provides Go bindings for MediaArea's MediaInfo library,
// for reading technical metadata (codecs, duration, resolution, bitrate, ...)
// from audio and video files.
//
// The libmediainfo shared library is loaded at runtime with dlopen, so it
// must be installed on the machine running the program, and its headers
// (MediaInfoDLL/MediaInfoDLL.h) must be available at build time. cgo is
// required.
//
// Basic usage:
//
//	mi, err := mediainfo.Open("movie.mp4")
//	if err != nil {
//		return err
//	}
//	defer mi.Close()
//
//	width, err := mi.Get(mediainfo.Video, 0, "Width")
package mediainfo
