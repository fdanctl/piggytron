package views

type LoginView struct {
	Initial bool

	Name     string
	Password string
	ErrorMsg string
}

func NewLoginView() *LoginView {
	return &LoginView{
		Initial: true,
	}
}

func (v *LoginView) ValidateName() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Name == "" {
		msgs = append(msgs, "Name is required")
	}
	return msgs
}

func (v *LoginView) NameHasError() bool {
	return len(v.ValidateName()) > 0 || v.ErrorMsg != ""
}

func (v *LoginView) ValidatePassword() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Password == "" {
		msgs = append(msgs, "Password is required")
	}
	return msgs
}

func (v *LoginView) PasswordHasError() bool {
	return len(v.ValidatePassword()) > 0 || v.ErrorMsg != ""
}

func (v *LoginView) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateName()...)
	msgs = append(msgs, v.ValidatePassword()...)
	return msgs
}
