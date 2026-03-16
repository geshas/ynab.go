package money_movement

import (
	"time"

	"github.com/geshas/ynab.go/api"
)

// MoneyMovement represents a money movement
type MoneyMovement struct {
	// Fields that are part of the documented MoneyMovement schema
	ID                   string     `json:"id"`
	Month                *api.Date  `json:"month"`
	MovedAt              *time.Time `json:"moved_at"`
	Note                 string     `json:"note"`
	MoneyMovementGroupID string     `json:"money_movement_group_id"`
	PerformedByUserID    string     `json:"performed_by_user_id"`
	FromCategoryID       string     `json:"from_category_id"`
	ToCategoryID         string     `json:"to_category_id"`
	Amount               int64      `json:"amount"`

	// Legacy/internal fields that should not be (un)marshalled for the
	// documented MoneyMovement API schema.
	CategoryID            string    `json:"-"`
	CategoryName          string    `json:"-"`
	Date                  *api.Date `json:"-"`
	PayeeID               *string   `json:"-"`
	PayeeName             *string   `json:"-"`
	RecurringJobID        *string   `json:"-"`
	RecurringJobType      *string   `json:"-"`
	ScheduledFlag         bool      `json:"-"`
	Approved              bool      `json:"-"`
	FlagColor             *string   `json:"-"`
	TransferAccountID     *string   `json:"-"`
	TransferTransactionID *string   `json:"-"`
	MatchedTransactionID  *string   `json:"-"`
	ImportID              *string   `json:"-"`
	Type                  string    `json:"-"`
	Isrenamed             bool      `json:"-"`
}

// MoneyMovementGroup represents a group of money movements
type MoneyMovementGroup struct {
	ID                string    `json:"id"`
	GroupCreatedAt    time.Time `json:"group_created_at"`
	Month             *api.Date `json:"month"`
	Note              string    `json:"note"`
	PerformedByUserID string    `json:"performed_by_user_id"`
}
