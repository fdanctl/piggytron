package views

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

type ExpenseForm struct {
	Initial bool

	Amount      string
	Description string
	Currency    string
	Date        string
	Category    string
	SourceAcc   string

	ErrorMsg    string
	CustomError error
}

func NewExpenseForm() *ExpenseForm {
	return &ExpenseForm{
		Initial:  true,
		Date:     time.Now().Format("02/01/2006"),
		Currency: currency.EUR.String(),
	}
}

func (v *ExpenseForm) ValidateAmount() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Amount == "" {
		msgs = append(msgs, "Amount is required")
	}

	str := strings.ReplaceAll(v.Amount, ",", "")
	str = strings.Replace(str, ".", "", 1)

	n, err := strconv.Atoi(str)
	if err != nil {
		return append(msgs, "Not a valid number")
	}

	if n < 0 {
		msgs = append(msgs, "Amount can't be negative")
	}
	return msgs
}

func (v *ExpenseForm) AmountHasError() bool {
	return len(v.ValidateAmount()) > 0
}

func (v *ExpenseForm) ValidateCurrency() (msgs []string) {
	if v.Initial {
		return
	}

	_, err := currency.ParseISO(v.Currency)
	if err != nil {
		msgs = append(msgs, v.Currency+" is not a valid currency")
	}

	return msgs
}

func (v *ExpenseForm) CurrencyHasError() bool {
	return len(v.ValidateCurrency()) > 0
}

func (v *ExpenseForm) ValidateDescription() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Description == "" {
		msgs = append(msgs, "Description is required")
	}
	if len(v.Description) > 50 {
		msgs = append(msgs, "Max length is 50 character")
	}
	return msgs
}

func (v *ExpenseForm) DescriptionHasError() bool {
	return len(v.ValidateDescription()) > 0
}

func (v *ExpenseForm) ValidateDate() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Date == "" {
		msgs = append(msgs, "Date is required")
	}

	_, err := time.Parse("02/01/2006", v.Date)
	if err != nil {
		return append(msgs, "Invalid date")
	}

	return msgs
}

func (v *ExpenseForm) DateHasError() bool {
	return len(v.ValidateDate()) > 0
}

func (v *ExpenseForm) ValidateCategory() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Category == "" {
		return append(msgs, "Category is required")
	}

	if _, err := uuid.Parse(v.Category); err != nil {
		msgs = append(msgs, "Invalid category")
	}

	return msgs
}

func (v *ExpenseForm) CategoryHasError() bool {
	return len(v.ValidateCategory()) > 0
}

func (v *ExpenseForm) ValidateSourceAcc() (msgs []string) {
	if v.Initial {
		return
	}

	if v.SourceAcc == "" {
		return append(msgs, "Souce is required")
	}

	if _, err := uuid.Parse(v.SourceAcc); err != nil {
		msgs = append(msgs, "Invalid source account")
	}

	return msgs
}

func (v *ExpenseForm) SourceAccHasError() bool {
	return len(v.ValidateSourceAcc()) > 0
}

func (v *ExpenseForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateAmount()...)
	msgs = append(msgs, v.ValidateCurrency()...)
	msgs = append(msgs, v.ValidateDescription()...)
	msgs = append(msgs, v.ValidateDate()...)
	msgs = append(msgs, v.ValidateCategory()...)
	msgs = append(msgs, v.ValidateSourceAcc()...)
	return msgs
}
