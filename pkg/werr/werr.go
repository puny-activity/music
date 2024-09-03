package werr

import "fmt"

func WrapSE(highLevelMessage string, lowLevelError error) error {
	return fmt.Errorf("%s: %w", highLevelMessage, lowLevelError)
}

func WrapES(highLevelError error, lowLevelMessage string) error {
	return fmt.Errorf("%w: %s", highLevelError, lowLevelMessage)
}
