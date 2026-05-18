package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/budget"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
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
	ps := params.Get("prev-amount")
	catType := params.Get("ctype")
	ptotalBudgeted := params.Get("total-budgeted")
	ptotalRowBudget := params.Get("total-row-budget")
	ptotalRowLeft := params.Get("total-row-left")
	pcatLeft := params.Get("cat-left")
	pleftToBudget := params.Get("ltb")
	pleftToSpent := params.Get("lts")
	pincome := params.Get("income")
	poverspent := params.Get("overspent")

	prev, err := strconv.Atoi(ps)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	i := strings.Index(amount, ".")
	cents := 0

	length := utf8.RuneCountInString(amount)
	if amount == "" {
		cents = 0
	} else if i == -1 {
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

	if cents == prev {
		logger.Debug("nothing to do")
		fmt.Fprintf(
			w,
			`
		<input id="prev-amount" type="hidden" name="prev-amount" value="%d"/>
		`,
			cents,
		)
		return
	}

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

	addedAmount := cents - prev
	fmt.Println(addedAmount)

	leftToSpent, err := strconv.Atoi(pleftToSpent)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	leftToSpent += addedAmount

	leftToBudget, err := strconv.Atoi(pleftToBudget)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	leftToBudget += addedAmount

	income, err := strconv.Atoi(pincome)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	income += addedAmount

	totalBudgeted, err := strconv.Atoi(ptotalBudgeted)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	totalBudgeted += addedAmount

	overspent, err := strconv.Atoi(poverspent)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	catLeft, err := strconv.Atoi(pcatLeft)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// if prev left is overspent reset this spent
	if catLeft < 0 {
		overspent += catLeft
	}
	// update
	catLeft += addedAmount
	// if new left is overspent add to overspent
	if catLeft < 0 {
		overspent -= catLeft
	}

	totalRowBudget, err := strconv.Atoi(ptotalRowBudget)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	totalRowBudget += addedAmount

	totalRowLeft, err := strconv.Atoi(ptotalRowLeft)
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	totalRowLeft += addedAmount

	fmt.Fprintf(
		w,
		`
		<input id="prev-amount" type="hidden" name="prev-amount" value="%d"/>
		`,
		cents,
	)

	obb := templ.Join(
		partials.CatRowLeftCell(cid, catLeft, templ.Attributes{
			"hx-swap-oob": "outerHTML",
		}),
		partials.BudgetStats(
			totalBudgeted,
			leftToBudget,
			income,
			leftToSpent,
			overspent,
			templ.Attributes{
				"hx-swap-oob": "outerHTML",
			},
		),
		partials.TotalRow(catType, totalRowBudget, totalRowLeft, templ.Attributes{
			"hx-swap-oob": "outerHTML",
		}),
	)

	obb.Render(r.Context(), w)

	fmt.Fprintf(
		w,
		`
		<span id="%s-pct" hx-swap-oob="outerHTML class="font-mono text-xs font-normal">
		(%s%%)
		</span>
		`,
		catType,
		views.FormatFloat((float64(totalRowBudget)/float64(totalBudgeted))*100),
	)

	fmt.Fprintf(
		w,
		`
		<span id="%s-pct" hx-swap-oob="outerHTML class="font-mono text-xs font-normal">
		(%s%%)
		</span>
		`,
		catType,
		views.FormatFloat((float64(totalRowBudget)/float64(totalBudgeted))*100),
	)
}
