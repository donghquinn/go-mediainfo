package mediainfo_test

import (
	"testing"

	"github.com/donghquinn/go-mediainfo"
)

func TestStreamKind_String(t *testing.T) {
	tests := []struct {
		give mediainfo.StreamKind
		want string
	}{
		{give: mediainfo.General, want: "General"},
		{give: mediainfo.Video, want: "Video"},
		{give: mediainfo.Audio, want: "Audio"},
		{give: mediainfo.Text, want: "Text"},
		{give: mediainfo.Other, want: "Other"},
		{give: mediainfo.Image, want: "Image"},
		{give: mediainfo.Menu, want: "Menu"},
		{give: mediainfo.StreamKind(-1), want: "StreamKind(-1)"},
		{give: mediainfo.StreamKind(42), want: "StreamKind(42)"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.give.String(); got != tt.want {
				t.Errorf("StreamKind(%d).String() = %q, want %q", int(tt.give), got, tt.want)
			}
		})
	}
}
