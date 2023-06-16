package offer

type CodecType string

const (
	Vp8      CodecType = "vp8"
	Vp9      CodecType = "vp9"
	Openh264 CodecType = "openh264"
	X264     CodecType = "x264"
)

type Codec struct {
	Type    CodecType `json:"type"`
	BitRate int       `json:"bitrate"`
}
