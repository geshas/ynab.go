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
	ID              string           `json:"id"`
	CategoryID      string           `json:"category_id"`
	CategoryName    string           `json:"category_name"`
	Income          bool             `json:"income"`
	GoalTarget      *int64           `json:"goal_target"`
	GoalTargetDate  *time.Time       `json:"goal_target_date"`
	GoalUnderfunded *bool            `json:"goal_underfunded"`
	GoalOverspent   *bool            `json:"goal_overspent"`
	MoneyMovements  []*MoneyMovement `json:"money_movements"`
}
