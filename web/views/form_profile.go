package views

import (
	"errors"

	"github.com/fdanctl/piggytron/internal/domain/user"
)

// ProfileForm is the view-model for the profile (name) form.
type ProfileForm struct {
	Form
	Name string
}

// NewProfileForm returns a profile form pre-filled with the current
// username.
func NewProfileForm(username string) *ProfileForm {
	f := ProfileForm{}
	f.Initial = true
	f.Name = username
	return &f
}

func (v *ProfileForm) ValidateName() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Name == "" {
		msgs = append(msgs, "Name is required")
	}
	if errors.Is(v.CustomError, user.ErrDuplicate) {
		msgs = append(msgs, "User already exists")
	}
	if len(v.Name) > 50 {
		msgs = append(msgs, "Max length is 50 character")
	}
	return msgs
}

func (v *ProfileForm) NameHasError() bool {
	return len(v.ValidateName()) > 0
}

func (v *ProfileForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateName()...)
	return msgs
}
