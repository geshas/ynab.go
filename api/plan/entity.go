// Package plan implements plan entities and services
package plan // import "github.com/geshas/ynab.go/api/plan"

import (
	"time"

	"github.com/geshas/ynab.go/api"
	"github.com/geshas/ynab.go/api/account"
	"github.com/geshas/ynab.go/api/category"
	"github.com/geshas/ynab.go/api/month"
	"github.com/geshas/ynab.go/api/payee"
	"github.com/geshas/ynab.go/api/transaction"
)

// Plan represents a plan
type Plan struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Accounts                 []*account.Account                     `json:"accounts"`
	Payees                   []*payee.Payee                         `json:"payees"`
	PayeeLocations           []*payee.Location                      `json:"payee_locations"`
	Categories               []*category.Category                   `json:"categories"`
	CategoryGroups           []*category.Group                      `json:"category_groups"`
	Months                   []*month.Month                         `json:"months"`
	Transactions             []*transaction.Summary                 `json:"transactions"`
	SubTransactions          []*transaction.SubTransaction          `json:"subtransactions"`
	ScheduledTransactions    []*transaction.ScheduledSummary        `json:"scheduled_transactions"`
	ScheduledSubTransactions []*transaction.ScheduledSubTransaction `json:"scheduled_sub_transactions"`

	DateFormat     *DateFormat     `json:"date_format"`
	CurrencyFormat *CurrencyFormat `json:"currency_format"`
	LastModifiedOn *time.Time      `json:"last_modified_on"`
	FirstMonth     *api.Date       `json:"first_month"`
	LastMonth      *api.Date       `json:"last_month"`
}

// Summary represents the summary of a plan
type Summary struct {
	ID       string             `json:"id"`
	Name     string             `json:"name"`
	Accounts []*account.Account `json:"accounts"`

	DateFormat     *DateFormat     `json:"date_format"`
	CurrencyFormat *CurrencyFormat `json:"currency_format"`
	LastModifiedOn *time.Time      `json:"last_modified_on"`
	FirstMonth     *api.Date       `json:"first_month"`
	LastMonth      *api.Date       `json:"last_month"`
}

// Snapshot represents a versioned snapshot for a plan
type Snapshot struct {
	Plan            *Plan
	ServerKnowledge uint64
}

// Settings represents the settings for a plan
type Settings struct {
	DateFormat     *DateFormat     `json:"date_format"`
	CurrencyFormat *CurrencyFormat `json:"currency_format"`
}

// DateFormat represents date format for a plan
type DateFormat struct {
	Format string `json:"format"`
}

// CurrencyFormat represents a currency format for a plan settings
type CurrencyFormat struct {
	ISOCode          string `json:"iso_code"`
	ExampleFormat    string `json:"example_format"`
	DecimalDigits    uint64 `json:"decimal_digits"`
	DecimalSeparator string `json:"decimal_separator"`
	GroupSeparator   string `json:"group_separator"`
	SymbolFirst      bool   `json:"symbol_first"`
	CurrencySymbol   string `json:"currency_symbol"`
	DisplaySymbol    bool   `json:"display_symbol"`
}
