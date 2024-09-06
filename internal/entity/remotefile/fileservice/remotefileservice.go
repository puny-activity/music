package fileservice

import (
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
)

type ID uuid.UUID

func NewID(uuidID uuid.UUID) ID {
	return ID(uuidID)
}

func GenerateID() ID {
	return ID(uuid.New())
}

type FileService struct {
	ID        *ID
	Address   string
	ScannedAt *carbon.Carbon
}