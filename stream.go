package mediainfo

import "strconv"

// StreamKind identifies a category of stream within a media file.
// The values match MediaInfo_stream_t in MediaInfoDLL.h.
type StreamKind int

// Stream kinds, in MediaInfo's numbering.
const (
	General StreamKind = iota // container-level info: format, duration, overall bitrate
	Video
	Audio
	Text  // subtitles and captions
	Other // timecode, chapters and other auxiliary tracks
	Image
	Menu
)

var _streamKindNames = [...]string{
	General: "General",
	Video:   "Video",
	Audio:   "Audio",
	Text:    "Text",
	Other:   "Other",
	Image:   "Image",
	Menu:    "Menu",
}

// String returns the MediaInfo name of the stream kind, e.g. "Video".
func (k StreamKind) String() string {
	if k >= 0 && int(k) < len(_streamKindNames) {
		return _streamKindNames[k]
	}
	return "StreamKind(" + strconv.Itoa(int(k)) + ")"
}
