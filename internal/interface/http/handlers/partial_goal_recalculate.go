package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/web/views"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// RemainderRecalculeHandler recalculates the remainder in goal complete form
type RemainderRecalculeHandler struct{}

func NewRemainderRecalculeHandler() *RemainderRecalculeHandler {
	return &RemainderRecalculeHandler{}
}

func (h *RemainderRecalculeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *RemainderRecalculeHandler) Get(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	balance := q.Get("balance")
	amountUsed := q.Get("amount-used")

	bal, err := strconv.Atoi(balance)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid amount", balance),
			fmt.Errorf("failed to convert amount '%s' to cents: %w", balance, err),
			"RemainderRecalculeHandler.Get",
		)
		httperror.SendError(w, r, err)
		return
	}

	used, err := convertAmountStrToInt(amountUsed)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid amount", amountUsed),
			fmt.Errorf("failed to convert amount '%s' to cents: %w", amountUsed, err),
			"RemainderRecalculeHandler.Get",
		)
		httperror.SendError(w, r, err)
		return
	}

	fmt.Fprint(w, views.FormatMoney(float64(bal-used)/100, currency.EUR, language.AmericanEnglish))
}
