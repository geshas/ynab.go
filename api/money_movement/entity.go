package money_movement

import (
	"time"

	"github.com/geshas/ynab.go/api"
)

// MoneyMovement represents a money movement
type MoneyMovement struct {
	ID                    string    `json:"id"`
	CategoryID            string    `json:"category_id"`
	CategoryName          string    `json:"category_name"`
	Date                  *api.Date `json:"date"`
	Amount                int64     `json:"amount"`
	PayeeID               *string   `json:"payee_id"`
	PayeeName             *string   `json:"payee_name"`
	RecurringJobID        *string   `json:"recurring_job_id"`
	RecurringJobType      *string   `json:"recurring_job_type"`
	ScheduledFlag         bool      `json:"scheduled_flag"`
	Approved              bool      `json:"approved"`
	FlagColor             *string   `json:"flag_color"`
	TransferAccountID     *string   `json:"transfer_account_id"`
	TransferTransactionID *string   `json:"transfer_transaction_id"`
	MatchedTransactionID  *string   `json:"matched_transaction_id"`
	ImportID              *string   `json:"import_id"`
	Type                  string    `json:"type"`
	Isrenamed             bool      `json:"isrenamed"`
}

// MoneyMovementGroup represents a group of money movements
type MoneyMovementGroup struct {
	GroupCreatedAt    time.Time        `json:"group_created_at"`
	Month             *api.Date        `json:"month"`
	Note              *string          `json:"note,omitempty"`
	PerformedByUserID string           `json:"performed_by_user_id"`
	MoneyMovements    []*MoneyMovement `json:"money_movements"`
}
