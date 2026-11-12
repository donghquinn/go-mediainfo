package mediainfo_test

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/donghquinn/go-mediainfo"
)

// testdata/sample.mp4 is a 1s, 64x48, 10fps H.264 video with an AAC audio track.
const _sampleFile = "testdata/sample.mp4"

func openSample(t *testing.T) *mediainfo.MediaInfo {
	t.Helper()

	mi, err := mediainfo.Open(_sampleFile)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", _sampleFile, err)
	}
	t.Cleanup(func() { _ = mi.Close() })

	return mi
}

func TestInit(t *testing.T) {
	for range 2 {
		if err := mediainfo.Init(); err != nil {
			t.Fatalf("Init() error = %v", err)
		}
	}
}

func TestVersion(t *testing.T) {
	got, err := mediainfo.Version()
	if err != nil {
		t.Fatalf("Version() error = %v", err)
	}
	if !strings.HasPrefix(got, "MediaInfoLib") {
		t.Errorf("Version() = %q, want prefix %q", got, "MediaInfoLib")
	}
}

func TestOpen(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		mi, err := mediainfo.Open("testdata/does-not-exist.mp4")
		if !errors.Is(err, mediainfo.ErrOpen) {
			t.Errorf("Open() error = %v, want %v", err, mediainfo.ErrOpen)
		}
		if mi != nil {
			t.Errorf("Open() = %v, want nil", mi)
		}
	})
}

func TestMediaInfo_Get(t *testing.T) {
	mi := openSample(t)

	tests := []struct {
		name       string
		giveKind   mediainfo.StreamKind
		giveStream int
		giveKey    string
		want       string
	}{
		{name: "container format", giveKind: mediainfo.General, giveKey: "Format", want: "MPEG-4"},
		{name: "video format", giveKind: mediainfo.Video, giveKey: "Format", want: "AVC"},
		{name: "video width", giveKind: mediainfo.Video, giveKey: "Width", want: "64"},
		{name: "video height", giveKind: mediainfo.Video, giveKey: "Height", want: "48"},
		{name: "audio format", giveKind: mediainfo.Audio, giveKey: "Format", want: "AAC"},
		{name: "audio sampling rate", giveKind: mediainfo.Audio, giveKey: "SamplingRate", want: "8000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := mi.Get(tt.giveKind, tt.giveStream, tt.giveKey)
			if err != nil {
				t.Fatalf("Get(%v, %d, %q) error = %v", tt.giveKind, tt.giveStream, tt.giveKey, err)
			}
			if got != tt.want {
				t.Errorf("Get(%v, %d, %q) = %q, want %q", tt.giveKind, tt.giveStream, tt.giveKey, got, tt.want)
			}
		})
	}
}

func TestMediaInfo_Get_errors(t *testing.T) {
	mi := openSample(t)

	tests := []struct {
		name       string
		giveKind   mediainfo.StreamKind
		giveStream int
		giveKey    string
		wantErr    error
	}{
		{name: "unknown key", giveKind: mediainfo.Video, giveKey: "NoSuchParameter", wantErr: mediainfo.ErrNotFound},
		{name: "missing stream", giveKind: mediainfo.Text, giveKey: "Format", wantErr: mediainfo.ErrNotFound},
		{name: "stream out of range", giveKind: mediainfo.Video, giveStream: 5, giveKey: "Format", wantErr: mediainfo.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mi.Get(tt.giveKind, tt.giveStream, tt.giveKey)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get(%v, %d, %q) error = %v, want %v", tt.giveKind, tt.giveStream, tt.giveKey, err, tt.wantErr)
			}
		})
	}

	t.Run("negative stream", func(t *testing.T) {
		if _, err := mi.Get(mediainfo.Video, -1, "Format"); err == nil {
			t.Error("Get(Video, -1, \"Format\") error = nil, want non-nil")
		}
	})
}

func TestMediaInfo_Get_concurrent(t *testing.T) {
	mi := openSample(t)

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 50 {
				if got, err := mi.Get(mediainfo.Video, 0, "Width"); err != nil || got != "64" {
					t.Errorf("Get(Video, 0, \"Width\") = %q, %v, want %q, nil", got, err, "64")
					return
				}
			}
		}()
	}
	wg.Wait()
}

func TestMediaInfo_Count(t *testing.T) {
	mi := openSample(t)

	tests := []struct {
		giveKind mediainfo.StreamKind
		want     int
	}{
		{giveKind: mediainfo.General, want: 1},
		{giveKind: mediainfo.Video, want: 1},
		{giveKind: mediainfo.Audio, want: 1},
		{giveKind: mediainfo.Text, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.giveKind.String(), func(t *testing.T) {
			got, err := mi.Count(tt.giveKind)
			if err != nil {
				t.Fatalf("Count(%v) error = %v", tt.giveKind, err)
			}
			if got != tt.want {
				t.Errorf("Count(%v) = %d, want %d", tt.giveKind, got, tt.want)
			}
		})
	}
}

func TestMediaInfo_Inform(t *testing.T) {
	t.Run("text", func(t *testing.T) {
		mi := openSample(t)

		got, err := mi.Inform()
		if err != nil {
			t.Fatalf("Inform() error = %v", err)
		}
		if !strings.Contains(got, "AVC") {
			t.Errorf("Inform() = %q, want it to contain %q", got, "AVC")
		}
	})

	t.Run("json", func(t *testing.T) {
		mi := openSample(t)

		if _, err := mi.Option("Inform", "JSON"); err != nil {
			t.Fatalf("Option(\"Inform\", \"JSON\") error = %v", err)
		}
		got, err := mi.Inform()
		if err != nil {
			t.Fatalf("Inform() error = %v", err)
		}

		var report struct {
			Media struct {
				Track []struct {
					Type string `json:"@type"`
				} `json:"track"`
			} `json:"media"`
		}
		if err := json.Unmarshal([]byte(got), &report); err != nil {
			t.Fatalf("Inform() returned invalid JSON: %v\n%s", err, got)
		}
		if n := len(report.Media.Track); n != 3 {
			t.Errorf("Inform() JSON has %d tracks, want 3", n)
		}
	})
}

func TestMediaInfo_Close(t *testing.T) {
	mi, err := mediainfo.Open(_sampleFile)
	if err != nil {
		t.Fatalf("Open(%q) error = %v", _sampleFile, err)
	}

	for range 2 {
		if err := mi.Close(); err != nil {
			t.Fatalf("Close() error = %v", err)
		}
	}

	if _, err := mi.Get(mediainfo.Video, 0, "Width"); !errors.Is(err, mediainfo.ErrClosed) {
		t.Errorf("Get() after Close error = %v, want %v", err, mediainfo.ErrClosed)
	}
	if _, err := mi.Count(mediainfo.Video); !errors.Is(err, mediainfo.ErrClosed) {
		t.Errorf("Count() after Close error = %v, want %v", err, mediainfo.ErrClosed)
	}
	if _, err := mi.Option("Inform", "JSON"); !errors.Is(err, mediainfo.ErrClosed) {
		t.Errorf("Option() after Close error = %v, want %v", err, mediainfo.ErrClosed)
	}
	if _, err := mi.Inform(); !errors.Is(err, mediainfo.ErrClosed) {
		t.Errorf("Inform() after Close error = %v, want %v", err, mediainfo.ErrClosed)
	}
}
