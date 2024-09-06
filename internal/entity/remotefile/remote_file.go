package remotefile

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/puny-activity/music/internal/entity/remotefile/contenttype"
	"github.com/puny-activity/music/pkg/werr"
)

type ID uuid.UUID

func ParseID(id string) (ID, error) {
	idUUID, err := uuid.Parse(id)
	if err != nil {
		return ID{}, werr.WrapSE("failed to parse uuid", err)
	}
	return ID(idUUID), nil
}

func (i ID) String() string {
	return uuid.UUID(i).String()
}

type File struct {
	ID   ID
	Name string
	Path string
}

type Changes struct {
	Created []FileInfo
	Updated []FileInfo
	Deleted []ID
}

type FileInfo struct {
	ID          ID
	Name        string
	ContentType contenttype.Type
	Path        string
	Size        int64
	Metadata    json.RawMessage
	MD5         string
}

func (e *FileInfo) GetAudioMetadata() *AudioMetadata {
	var m AudioMetadata
	err := json.Unmarshal(e.Metadata, &m)
	if err != nil {
		return nil
	}
	return &m
}

func (e *FileInfo) GetImageMetadata() *ImageMetadata {
	var m ImageMetadata
	err := json.Unmarshal(e.Metadata, &m)
	if err != nil {
		return nil
	}
	return &m
}

type AudioMetadata struct {
	Title        *string `json:"title,omitempty"`
	DurationNs   int64   `json:"durationNs,omitempty"`
	Artist       *string `json:"artist,omitempty"`
	Album        *string `json:"album,omitempty"`
	Genre        *string `json:"genre,omitempty"`
	Year         *int    `json:"year,omitempty"`
	TrackNumber  *int    `json:"trackNumber,omitempty"`
	Comment      *string `json:"comment,omitempty"`
	Channels     int     `json:"channels,omitempty"`
	BitrateKbps  int     `json:"bitrateKbps,omitempty"`
	SampleRateHz int     `json:"sampleRateHz,omitempty"`
}

type ImageMetadata struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}
