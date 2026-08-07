package views

type ChangePasswordForm struct {
	Form

	CurrentPassword    string
	NewPassword        string
	NewPasswordConfirm string
}

func NewChangePasswordForm() *ChangePasswordForm {
	f := ChangePasswordForm{}
	f.Initial = true
	return &f
}

func (v *ChangePasswordForm) ValidateCurrentPassword() (msgs []string) {
	if v.Initial {
		return
	}
	if v.CurrentPassword == "" {
		msgs = append(msgs, "Current password is required")
	}
	// if errors.Is(v.CustomError, user.ErrWrongPassword) {
	// 	msgs = append(msgs, "Unable to verify password")
	// }
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
