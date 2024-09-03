package contenttype

import "fmt"

type Type string

const (
	Unknown                Type = "unknown"
	ApplicationOctetStream Type = "application/octet-stream"
	AudioMpeg              Type = "audio/mpeg"
	AudioWav               Type = "audio/wav"
	AudioXWav              Type = "audio/x-wav"
	AudioXFlac             Type = "audio/x-flac"
	ImagePng               Type = "image/png"
	ImageJpeg              Type = "image/jpeg"
	ImageGif               Type = "image/gif"
)

var fileContentTypes = map[string]Type{
	ApplicationOctetStream.String(): ApplicationOctetStream,
	AudioMpeg.String():              AudioMpeg,
	AudioWav.String():               AudioWav,
	AudioXWav.String():              AudioXWav,
	AudioXFlac.String():             AudioXFlac,
	ImagePng.String():               ImagePng,
	ImageJpeg.String():              ImageJpeg,
	ImageGif.String():               ImageGif,
}

func New(typeName string) (Type, error) {
	fileContentType, ok := fileContentTypes[typeName]
	if !ok {
		return Unknown, fmt.Errorf("unknown content type: %s", typeName)
	}
	return fileContentType, nil
}

func (e Type) IsAudio() bool {
	switch e {
	case AudioMpeg, AudioWav, AudioXWav, AudioXFlac:
		return true
	default:
		return false
	}
}

func (e Type) IsImage() bool {
	switch e {
	case ImagePng, ImageJpeg, ImageGif:
		return true
	default:
		return false
	}
}

func (e Type) String() string {
	return string(e)
}
