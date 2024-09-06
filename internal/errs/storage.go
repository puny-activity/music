package errs

import "errors"

type internalError struct {
	error error
	code  string
}

type Storage struct {
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Error(err error) string {
	for i := range errorList {
		if errors.Is(err, errorList[i].error) {
			return errorList[i].code
		}
	}
	return unexpectedError.code
}

func (s *Storage) Code(err error) string {
	for i := range errorList {
		if errors.Is(err, errorList[i].error) {
			return errorList[i].code
		}
	}
	return unexpectedError.code
}
