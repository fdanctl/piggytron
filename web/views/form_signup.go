package views

import (
	"errors"

	"github.com/fdanctl/piggytron/internal/application/user"
)

type SignupView struct {
	Initial bool

	Name            string
	Password        string
	PasswordConfirm string
	ErrorMsg        string
	CustomError     error
}

func NewSignupView() *SignupView {
	return &SignupView{
		Initial: true,
	}
}

func (v *SignupView) ValidateName() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Name == "" {
		msgs = append(msgs, "Name is required")
	}
	if errors.Is(v.CustomError, user.ErrUserExists) {
		msgs = append(msgs, v.CustomError.Error())
	}
	return msgs
}

func (v *SignupView) NameHasError() bool {
	return len(v.ValidateName()) > 0
}

func (v *SignupView) ValidatePassword() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Password == "" {
		msgs = append(msgs, "Password is required")
	}
	if v.PasswordConfirm == "" {
		msgs = append(msgs, "Password is required")
	}
	if v.PasswordConfirm != v.Password {
		msgs = append(msgs, "Password doesn't match")
	}
	return msgs
}

func (v *SignupView) PasswordHasError() bool {
	return len(v.ValidatePassword()) > 0
}

func (v *SignupView) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateName()...)
	msgs = append(msgs, v.ValidatePassword()...)
	return msgs
}

func StringArrToStr(arr []string) string {
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}
