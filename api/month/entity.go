// Package month implements month entities and services
package month // import "github.com/geshas/ynab.go/api/month"

import (
	"github.com/geshas/ynab.go/api"
	"github.com/geshas/ynab.go/api/category"
)

// Month represents a month for a budget
// Each budget contains one or more months, which is where To be Budgeted,
// Age of Money and Category (budgeted / activity / balances)
// amounts are available.
type Month struct {
	Month      api.Date             `json:"month"`
	Categories []*category.Category `json:"categories"`
	Deleted    bool                 `json:"deleted"`

	Note         *string `json:"note"`
	ToBeBudgeted *int64  `json:"to_be_budgeted"`
	AgeOfMoney   *int64  `json:"age_of_money"`

	// Income the total amount in transactions categorized to "Inflow: To be Budgeted"
	// in the month (milliunits format)
	Income *int64 `json:"income"`
	// IncomeFormatted total income formatted in the plan's currency format
	IncomeFormatted *string `json:"income_formatted"`
	// IncomeCurrency total income as a decimal currency amount
	IncomeCurrency *float64 `json:"income_currency"`
	// Budgeted the total amount budgeted in the month (milliunits format)
	Budgeted *int64 `json:"budgeted"`
	// BudgetedFormatted total budgeted amount formatted in the plan's currency format
	BudgetedFormatted *string `json:"budgeted_formatted"`
	// BudgetedCurrency total budgeted amount as a decimal currency amount
	BudgetedCurrency *float64 `json:"budgeted_currency"`
	// Activity the total amount in transactions in the month, excluding those
	// categorized to "Inflow: To be Budgeted" (milliunits format)
	Activity *int64 `json:"activity"`
	// ActivityFormatted total activity amount formatted in the plan's currency format
	ActivityFormatted *string `json:"activity_formatted"`
	// ActivityCurrency total activity amount as a decimal currency amount
	ActivityCurrency *float64 `json:"activity_currency"`
	// ToBeBudgetedFormatted ready to assign amount formatted in the plan's currency format
	ToBeBudgetedFormatted *string `json:"to_be_budgeted_formatted"`
	// ToBeBudgetedCurrency ready to assign amount as a decimal currency amount
	ToBeBudgetedCurrency *float64 `json:"to_be_budgeted_currency"`
}

// Summary represents the summary of a month for a budget
// Each budget contains one or more months, which is where To be Budgeted,
// Age of Money and Category (budgeted / activity / balances)
// amounts are available.
type Summary struct {
	Month   api.Date `json:"month"`
	Deleted bool     `json:"deleted"`

	Note         *string `json:"note"`
	ToBeBudgeted *int64  `json:"to_be_budgeted"`
	AgeOfMoney   *int64  `json:"age_of_money"`

	// Income the total amount in transactions categorized to "Inflow: To be Budgeted"
	// in the month (milliunits format)
	Income *int64 `json:"income"`
	// IncomeFormatted total income formatted in the plan's currency format
	IncomeFormatted *string `json:"income_formatted"`
	// IncomeCurrency total income as a decimal currency amount
	IncomeCurrency *float64 `json:"income_currency"`
	// Budgeted the total amount budgeted in the month (milliunits format)
	Budgeted *int64 `json:"budgeted"`
	// BudgetedFormatted total budgeted amount formatted in the plan's currency format
	BudgetedFormatted *string `json:"budgeted_formatted"`
	// BudgetedCurrency total budgeted amount as a decimal currency amount
	BudgetedCurrency *float64 `json:"budgeted_currency"`
	// Activity the total amount in transactions in the month, excluding those
	// categorized to "Inflow: To be Budgeted" (milliunits format)
	Activity *int64 `json:"activity"`
	// ActivityFormatted total activity amount formatted in the plan's currency format
	ActivityFormatted *string `json:"activity_formatted"`
	// ActivityCurrency total activity amount as a decimal currency amount
	ActivityCurrency *float64 `json:"activity_currency"`
	// ToBeBudgetedFormatted ready to assign amount formatted in the plan's currency format
	ToBeBudgetedFormatted *string `json:"to_be_budgeted_formatted"`
	// ToBeBudgetedCurrency ready to assign amount as a decimal currency amount
	ToBeBudgetedCurrency *float64 `json:"to_be_budgeted_currency"`
}

// SearchResultSnapshot represents a versioned snapshot for a month search
type SearchResultSnapshot struct {
	Months          []*Summary
	ServerKnowledge uint64
}
