package views

import "github.com/fdanctl/piggytron/internal/domain/account"

type Bank struct {
	ID   string
	Name string
	Type string
	// Amount       string
}

func NewBank(a *account.Account) Bank {
	return Bank{
		ID:   string(a.ID()),
		Name: a.Name(),
		Type: string(a.Type()),
	}
}
