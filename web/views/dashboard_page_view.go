package views

import (
	"math"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/query"
)

type DashboardPage struct {
	BudgetReadyToAssign int
	Income              int
	Spent               int
	SpentPct            float64 // based on budget
	Savings             int
	Balance             int
	Budget              int
	AvgSpent            int
	AvgLeft             int
	CurrDay             int
	MonthDayLength      int

	Goals              []Goal
	RecentTransactions []Transaction
}

func NewDashboardPage(
	a []query.AccountWithSum,
	t []query.LedgerEntryDTO,
	b []query.CategoryBudgetValue,
) DashboardPage {
	var savings int
	var balance int

	var income int
	var spent int

	var budget int
	var totalAvailable int

	for _, v := range b {
		if v.Type == "income" {
			income += v.Value
			continue
		}
		budget += v.Budgeted
		spent += v.Value
		carryover := v.PrevTotalBudget - v.PrevTotalSpent
		available := v.Budgeted - v.Value + carryover
		totalAvailable += available
	}

	var goals []Goal
	for _, v := range a {
		if v.Type == "bank" {
			if *v.IsSaving {
				savings += v.Sum
			} else {
				balance += v.Sum
			}
		}

		if v.Type == string(account.GoalType) {
			goals = append(goals, NewGoal(v))
		}
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	nextMonthFirst := time.Date(today.Year(), today.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	monthLen := time.Date(
		nextMonthFirst.Year(),
		nextMonthFirst.Month(),
		nextMonthFirst.Day()-1,
		0,
		0,
		0,
		0,
		time.UTC,
	).Day()

	var tviews []Transaction
	for _, v := range t {
		tviews = append(tviews, NewTransaction(v))
	}

	var spentPct float64
	if budget == 0 {
		spentPct = 100
	} else {
		spentPct = float64(spent) / float64(budget) * 100
	}

	return DashboardPage{
		BudgetReadyToAssign: balance - totalAvailable,
		Budget:              budget,
		Income:              income,
		Spent:               spent,
		SpentPct:            spentPct,
		Savings:             savings,
		Balance:             balance,
		AvgSpent:            int(math.Round(float64(spent) / float64(now.Day()))),
		AvgLeft: int(
			math.Round(float64(income-spent) / float64(monthLen-today.Day()-1)),
		),
		CurrDay:            today.Day(),
		MonthDayLength:     monthLen,
		Goals:              goals,
		RecentTransactions: tviews,
	}
}
