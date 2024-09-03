package zerologger

import (
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"time"
)

func newConsoleWriter() zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
}

func NewLogger(logLevel string) (*zerolog.Logger, error) {
	if logLevel == "" {
		logLevel = "trace"
	}
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	consoleWriter := newConsoleWriter()
	consoleWriter.FormatTimestamp = func(i interface{}) string {
		return ""
	}
	consoleWriter.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("\n    %s: ", i)
	}
	consoleWriter.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	logger := zerolog.New(consoleWriter).Level(level).With().Logger()
	return &logger, nil
}
