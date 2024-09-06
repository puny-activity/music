package contenttype

import "fmt"

type Type string

const (
	Unknown                Type = "unknown"
	ApplicationOctetStream Type = "application/octet-stream"
	// Audio
	AudioAac      Type = "audio/aac"
	AudioFlac     Type = "audio/flac"
	AudioMatroska Type = "audio/matroska"
	AudioMp4      Type = "audio/mp4"
	AudioMpeg     Type = "audio/mpeg"
	AudioMpeg3    Type = "audio/mpeg3"
	AudioOgg      Type = "audio/ogg"
	AudioWav      Type = "audio/wav"
	AudioWebm     Type = "audio/webm"
	AudioXFlac    Type = "audio/x-flac"
	AudioXMpeg3   Type = "audio/x-mpeg-3"
	AudioXWav     Type = "audio/x-wav"
	// Image
	ImageGif  Type = "image/gif"
	ImageJpeg Type = "image/jpeg"
	ImagePng  Type = "image/png"
)

var fileContentTypes = map[string]Type{
	ApplicationOctetStream.String(): ApplicationOctetStream,
	AudioAac.String():               AudioAac,
	AudioFlac.String():              AudioFlac,
	AudioMatroska.String():          AudioMatroska,
	AudioMp4.String():               AudioMp4,
	AudioMpeg.String():              AudioMpeg,
	AudioMpeg3.String():             AudioMpeg3,
	AudioOgg.String():               AudioOgg,
	AudioWav.String():               AudioWav,
	AudioWebm.String():              AudioWebm,
	AudioXFlac.String():             AudioXFlac,
	AudioXMpeg3.String():            AudioXMpeg3,
	AudioXWav.String():              AudioXWav,
	ImageGif.String():               ImageGif,
	ImageJpeg.String():              ImageJpeg,
	ImagePng.String():               ImagePng,
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
	case AudioAac, AudioFlac, AudioMatroska, AudioMp4, AudioMpeg, AudioMpeg3, AudioOgg, AudioWav, AudioWebm, AudioXFlac,
		AudioXMpeg3, AudioXWav:
		return true
	default:
		return false
	}
}

func (e Type) IsImage() bool {
	switch e {
	case ImageGif, ImageJpeg, ImagePng:
		return true
	default:
		return false
	}
}

func (e Type) String() string {
	return string(e)
}
