// Package account implements account entities and services
package account // import "github.com/geshas/ynab.go/api/account"

// LoanAccountPeriodicValue represents periodic values for loan accounts
// keyed by date strings in YYYY-MM-DD format (e.g., "2024-01-01").
// Values are int64 amounts in milliunits format.
//
// Example:
//
//	{
//	  "2024-01-01": 425000,  // $425.00 payment starting Jan 1, 2024
//	  "2024-06-01": 450000   // $450.00 payment starting Jun 1, 2024
//	}
type LoanAccountPeriodicValue map[string]int64

// Account represents an account for a budget
type Account struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     Type   `json:"type"`
	OnBudget bool   `json:"on_budget"`
	// Balance The current balance of the account in milliunits format
	Balance int64 `json:"balance"`
	// BalanceFormatted The current balance formatted in the plan's currency format
	BalanceFormatted *string `json:"balance_formatted"`
	// BalanceCurrency The current balance as a decimal currency amount
	BalanceCurrency *float64 `json:"balance_currency"`
	// ClearedBalance The current cleared balance of the account in milliunits format
	ClearedBalance int64 `json:"cleared_balance"`
	// ClearedBalanceFormatted The current cleared balance formatted in the plan's currency format
	ClearedBalanceFormatted *string `json:"cleared_balance_formatted"`
	// ClearedBalanceCurrency The current cleared balance as a decimal currency amount
	ClearedBalanceCurrency *float64 `json:"cleared_balance_currency"`
	// UnclearedBalance The current uncleared balance of the account in milliunits format
	UnclearedBalance int64 `json:"uncleared_balance"`
	// UnclearedBalanceFormatted The current uncleared balance formatted in the plan's currency format
	UnclearedBalanceFormatted *string `json:"uncleared_balance_formatted"`
	// UnclearedBalanceCurrency The current uncleared balance as a decimal currency amount
	UnclearedBalanceCurrency *float64 `json:"uncleared_balance_currency"`
	// TransferPayeeID The payee id which should be used when transferring to this account
	TransferPayeeID *string `json:"transfer_payee_id"`
	Closed          bool    `json:"closed"`
	// Deleted Deleted accounts will only be included in delta requests
	Deleted bool `json:"deleted"`

	Note                *string                   `json:"note"`
	DirectImportLinked  bool                      `json:"direct_import_linked"`
	DirectImportInError bool                      `json:"direct_import_in_error"`
	LastReconciledAt    *string                   `json:"last_reconciled_at"`
	DebtOriginalBalance *int64                    `json:"debt_original_balance"`
	DebtInterestRates   *LoanAccountPeriodicValue `json:"debt_interest_rates"`
	DebtMinimumPayments *LoanAccountPeriodicValue `json:"debt_minimum_payments"`
	DebtEscrowAmounts   *LoanAccountPeriodicValue `json:"debt_escrow_amounts"`
}

// SearchResultSnapshot represents a versioned snapshot for an account search
type SearchResultSnapshot struct {
	Accounts        []*Account
	ServerKnowledge uint64
}
