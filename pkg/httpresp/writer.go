package httpresp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Writer struct {
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(writer http.ResponseWriter, statusCode int, payload any) error {
	writer.Header().Set("Content-Type", "application/json")

	writer.WriteHeader(statusCode)

	if payload != nil {
		err := json.NewEncoder(writer).Encode(payload)
		if err != nil {
			return fmt.Errorf("failed to encode payload")
		}
	}

	return nil
}
