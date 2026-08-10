package views

import (
	"errors"

	"github.com/fdanctl/piggytron/internal/domain/user"
)

// ChangePasswordForm is the view-model for the change-password form.
type ChangePasswordForm struct {
	Form

	CurrentPassword    string
	NewPassword        string
	NewPasswordConfirm string
}

// NewChangePasswordForm returns a blank change-password form.
func NewChangePasswordForm() *ChangePasswordForm {
	f := ChangePasswordForm{}
	f.Initial = true
	return &f
}

func (v *ChangePasswordForm) ValidateCurrentPassword() (msgs []string) {
	if v.Initial {
		return
	}

	if errors.Is(v.CustomError, user.ErrWrongPassword) {
		msgs = append(msgs, v.ErrorMsg)
	}
	if v.CurrentPassword == "" {
		msgs = append(msgs, "Current password is required")
	}

	return msgs
}

func (v *ChangePasswordForm) CurrentPasswordHasError() bool {
	return len(v.ValidateCurrentPassword()) > 0
}

func (v *ChangePasswordForm) ValidateNewPassword() (msgs []string) {
	if v.Initial {
		return
	}

	if v.NewPassword == "" {
		msgs = append(msgs, "Password is required")
	}
	if v.NewPasswordConfirm == "" {
		msgs = append(msgs, "Password is required")
	}
	if v.NewPasswordConfirm != v.NewPassword {
		msgs = append(msgs, "Password doesn't match")
	}
	return msgs
}

func (v *ChangePasswordForm) NewPasswordHasError() bool {
	return len(v.ValidateNewPassword()) > 0
}

func (v *ChangePasswordForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateCurrentPassword()...)
	msgs = append(msgs, v.ValidateNewPassword()...)
	return msgs
}
