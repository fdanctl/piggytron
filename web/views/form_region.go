package views

type RegionForm struct {
	Form
	Currency string
	Language string
}

func NewRegionForm() *RegionForm {
	f := RegionForm{}
	f.Initial = true
	return &f
}

func (v *RegionForm) ValidateCurrency() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Currency == "" {
		msgs = append(msgs, "Currency is required")
	}
	if len(v.Currency) > 50 {
		msgs = append(msgs, "Max length is 50 character")
	}
	return msgs
}

func (v *RegionForm) CurrencyHasError() bool {
	return len(v.ValidateCurrency()) > 0
}

func (v *RegionForm) ValidateLanguage() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Language == "" {
		msgs = append(msgs, "Language is required")
	}
	if len(v.Language) > 50 {
		msgs = append(msgs, "Max length is 50 character")
	}
	return msgs
}

func (v *RegionForm) LanguageHasError() bool {
	return len(v.ValidateLanguage()) > 0
}

func (v *RegionForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateCurrency()...)
	return msgs
}
