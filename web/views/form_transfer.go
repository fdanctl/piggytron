package views

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

type TransferForm struct {
	Initial bool

	Amount         string
	Description    string
	Currency       string
	Date           string
	Category       string
	SourceAcc      string
	DestinationAcc string

	ErrorMsg    string
	CustomError error
}

func NewTransferForm() *TransferForm {
	return &TransferForm{
		Initial:  true,
		Date:     time.Now().Format("02/01/2006"),
		Currency: currency.EUR.String(),
	}
}

func (v *TransferForm) ValidateAmount() (msgs []string) {
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

func (v *TransferForm) AmountHasError() bool {
	return len(v.ValidateAmount()) > 0
}

func (v *TransferForm) ValidateCurrency() (msgs []string) {
	if v.Initial {
		return
	}

	_, err := currency.ParseISO(v.Currency)
	if err != nil {
		msgs = append(msgs, v.Currency+" is not a valid currency")
	}

	return msgs
}

func (v *TransferForm) CurrencyHasError() bool {
	return len(v.ValidateCurrency()) > 0
}

func (v *TransferForm) ValidateDescription() (msgs []string) {
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

func (v *TransferForm) DescriptionHasError() bool {
	return len(v.ValidateDescription()) > 0
}

func (v *TransferForm) ValidateDate() (msgs []string) {
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

func (v *TransferForm) DateHasError() bool {
	return len(v.ValidateDate()) > 0
}

func (v *TransferForm) ValidateCategory() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Category == "" {
		return
	}

	if _, err := uuid.Parse(v.Category); err != nil {
		msgs = append(msgs, "Invalid category")
	}

	return msgs
}

func (v *TransferForm) CategoryHasError() bool {
	return len(v.ValidateCategory()) > 0
}

func (v *TransferForm) ValidateSourceAcc() (msgs []string) {
	if v.Initial {
		return
	}

	if v.SourceAcc == "" {
		return append(msgs, "Souce is required")
	}

	if _, err := uuid.Parse(v.SourceAcc); err != nil {
		msgs = append(msgs, "Invalid source account")
	}

	if v.SourceAcc == v.DestinationAcc {
		msgs = append(msgs, "Source account can't be the same has destination")
	}

	return msgs
}

func (v *TransferForm) SourceAccHasError() bool {
	return len(v.ValidateSourceAcc()) > 0
}

func (v *TransferForm) ValidateDestinationAcc() (msgs []string) {
	if v.Initial {
		return
	}

	if v.DestinationAcc == "" {
		return append(msgs, "Destination is required")
	}

	if _, err := uuid.Parse(v.DestinationAcc); err != nil {
		msgs = append(msgs, "Invalid destination account")
	}

	if v.SourceAcc == v.DestinationAcc {
		msgs = append(msgs, "Source account can't be the same has destination")
	}

	return msgs
}

func (v *TransferForm) DestinationAccHasError() bool {
	return len(v.ValidateDestinationAcc()) > 0
}

func (v *TransferForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateAmount()...)
	msgs = append(msgs, v.ValidateCurrency()...)
	msgs = append(msgs, v.ValidateDescription()...)
	msgs = append(msgs, v.ValidateDate()...)
	msgs = append(msgs, v.ValidateCategory()...)
	msgs = append(msgs, v.ValidateSourceAcc()...)
	msgs = append(msgs, v.ValidateDestinationAcc()...)
	return msgs
}
