package email

import (
	"errors"
	"github.com/puny-activity/music/pkg/werr"
	"net/mail"
)

type Email struct {
	email string
}

func New(email string) (Email, error) {
	if email == "" {
		return Email{}, errors.New("NotProvidedEmail")
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return Email{}, werr.WrapES(errors.New("InvalidEmail"), err.Error())
	}

	return Email{
		email: email,
	}, nil
}

func (e Email) String() string {
	return e.email
}
