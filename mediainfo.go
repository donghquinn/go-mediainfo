package mediainfo

/*
#cgo linux LDFLAGS: -ldl
#cgo darwin LDFLAGS: -framework CoreFoundation
#cgo darwin,arm64 CFLAGS: -I/opt/homebrew/include
#cgo darwin,amd64 CFLAGS: -I/usr/local/include
#include <stdlib.h>
#include "mediainfo_wrapper.h"
*/
import "C"

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"
)

var (
	// ErrLoad is returned when the MediaInfo shared library cannot be loaded.
	ErrLoad = errors.New("mediainfo: cannot load libmediainfo shared library")
	// ErrOpen is returned when a file cannot be opened or parsed.
	ErrOpen = errors.New("mediainfo: cannot open file")
	// ErrNotFound is returned when a parameter has no value for the requested stream.
	ErrNotFound = errors.New("mediainfo: no value for parameter")
	// ErrClosed is returned when a method is called on a closed MediaInfo.
	ErrClosed = errors.New("mediainfo: use of closed handle")
)

var (
	_loadOnce sync.Once
	_loadErr  error
)

// Init loads the MediaInfo shared library. It is safe to call multiple
// times and from multiple goroutines; only the first call does any work.
//
// Calling Init is optional: Open and Version load the library on demand.
// Call it explicitly to detect a missing library early, e.g. at startup.
func Init() error {
	_loadOnce.Do(func() {
		if C.mediainfo_c_load() == 0 {
			_loadErr = ErrLoad
		}
	})
	return _loadErr
}

// Version returns the version string of the loaded MediaInfo library,
// e.g. "MediaInfoLib - v24.12".
func Version() (string, error) {
	if err := Init(); err != nil {
		return "", err
	}
	return option(nil, "Info_Version", ""), nil
}

// MediaInfo is a handle to a parsed media file.
//
// A MediaInfo is safe for concurrent use. It must be released with Close.
type MediaInfo struct {
	mu  sync.Mutex
	ptr unsafe.Pointer
}

// Open opens and parses the file at path.
func Open(path string) (*MediaInfo, error) {
	if err := Init(); err != nil {
		return nil, err
	}

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	ptr := C.mediainfo_c_open(cpath)
	if ptr == nil {
		return nil, fmt.Errorf("%w: %s", ErrOpen, path)
	}

	return &MediaInfo{ptr: ptr}, nil
}

// Get returns the value of parameter key for the stream-th stream of the
// given kind. Streams are numbered from 0.
//
// The available parameters are listed by:
//
//	mediainfo --Info-Parameters
func (mi *MediaInfo) Get(kind StreamKind, stream int, key string) (string, error) {
	if stream < 0 {
		return "", fmt.Errorf("mediainfo: invalid stream number %d", stream)
	}

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	mi.mu.Lock()
	defer mi.mu.Unlock()

	if mi.ptr == nil {
		return "", ErrClosed
	}

	// The returned buffer is owned by the handle and only valid until the
	// next call, so it is copied while the lock is held.
	val := C.GoString(C.mediainfo_c_get(mi.ptr, C.int(kind), C.size_t(stream), ckey))
	if val == "" {
		return "", fmt.Errorf("%w: %s %s #%d", ErrNotFound, kind, key, stream)
	}

	return val, nil
}

// Count returns the number of streams of the given kind.
func (mi *MediaInfo) Count(kind StreamKind) (int, error) {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	if mi.ptr == nil {
		return 0, ErrClosed
	}

	return int(C.mediainfo_c_count(mi.ptr, C.int(kind))), nil
}

// Option sets a MediaInfo option on this handle and returns the library's
// response. For example, Option("Inform", "JSON") makes Inform return JSON.
func (mi *MediaInfo) Option(name, value string) (string, error) {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	if mi.ptr == nil {
		return "", ErrClosed
	}

	return option(mi.ptr, name, value), nil
}

// Inform returns a report of all streams, formatted according to the
// "Inform" option (plain text by default).
func (mi *MediaInfo) Inform() (string, error) {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	if mi.ptr == nil {
		return "", ErrClosed
	}

	return C.GoString(C.mediainfo_c_inform(mi.ptr)), nil
}

// Close releases the handle. Calling Close more than once is a no-op.
func (mi *MediaInfo) Close() error {
	mi.mu.Lock()
	defer mi.mu.Unlock()

	if mi.ptr != nil {
		C.mediainfo_c_close(mi.ptr)
		mi.ptr = nil
	}

	return nil
}

func option(ptr unsafe.Pointer, name, value string) string {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	return C.GoString(C.mediainfo_c_option(ptr, cname, cvalue))
}
