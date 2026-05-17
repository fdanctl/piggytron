package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/fdanctl/piggytron/internal/application/budget"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/views"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type BudgetHandler struct {
	service *budget.Service
}

func NewBudgetHandler(s *budget.Service) *BudgetHandler {
	return &BudgetHandler{
		service: s,
	}
}

func (h *BudgetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BudgetHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionID, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r.ParseForm()
	params := r.Form
	amount := params.Get("amount")
	cid := params.Get("cid")
	bid := params.Get("bid")
	ptotal := params.Get("total-budget")

	total, err := strconv.Atoi(ptotal)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	i := strings.Index(amount, ".")
	cents := 0

	length := utf8.RuneCountInString(amount)
	if i == -1 {
		cents, err = strconv.Atoi(amount)
		if err != nil {
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}
		cents *= 100
	} else if length-1-i > 2 {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	} else {
		for length-i < 3 {
			amount += "0"
			length++
		}

		cents, err = strconv.Atoi(strings.Replace(amount, ".", "", 1))
		if err != nil {
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}
	}

	now := time.Now()

	fmt.Println(total)
	total += cents
	fmt.Println(total)

	if bid == "" || bid == "00000000-0000-0000-0000-000000000000" {
		_, err := h.service.CreateBudget(r.Context(), sessionID.UserID, cid, now, cents)
		if err != nil {
			logger.Error("error creating budget", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		err := h.service.UpdateBudgetAmount(r.Context(), bid, cents)
		if err != nil {
			logger.Error("error updating budget", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	str := views.FormatMoney(float64(total)/100, currency.EUR, language.AmericanEnglish)
	fmt.Fprint(
		w,
		`
		<input type="hidden" name="total-budget" value={ totalBudget }/>
		`,
		str,
	)
}
