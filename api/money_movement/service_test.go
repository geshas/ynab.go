package money_movement_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/geshas/ynab.go"
	"github.com/geshas/ynab.go/api"
)

func TestService_GetMoneyMovements(t *testing.T) {
	t.Run(`success`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/money_movements"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movements": [
      {
        "id": "mm-123",
        "category_id": "cat-456",
        "category_name": "Groceries",
        "date": "2024-01-15",
        "amount": -15000,
        "payee_id": "payee-789",
        "payee_name": "Grocery Store",
        "recurring_job_id": null,
        "recurring_job_type": null,
        "scheduled_flag": false,
        "approved": true,
        "flag_color": null,
        "transfer_account_id": null,
        "transfer_transaction_id": null,
        "matched_transaction_id": null,
        "import_id": null,
        "type": "withdrawal",
        "isrenamed": false
      }
    ],
    "server_knowledge": 0
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		movements, err := client.MoneyMovement().GetMoneyMovements("plan-id-123")
		assert.NoError(t, err)
		assert.NotNil(t, movements)
		assert.Len(t, movements.MoneyMovements, 1)

		expectedDate, _ := api.DateFromString("2024-01-15")
		assert.Equal(t, "mm-123", movements.MoneyMovements[0].ID)
		assert.Equal(t, "cat-456", movements.MoneyMovements[0].CategoryID)
		assert.Equal(t, "Groceries", movements.MoneyMovements[0].CategoryName)
		assert.Equal(t, &expectedDate, movements.MoneyMovements[0].Date)
		assert.Equal(t, int64(-15000), movements.MoneyMovements[0].Amount)
	})

	t.Run(`success empty`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/money_movements"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movements": [],
    "server_knowledge": 0
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		movements, err := client.MoneyMovement().GetMoneyMovements("plan-id-123")
		assert.NoError(t, err)
		assert.Len(t, movements.MoneyMovements, 0)
	})
}

func TestService_GetMoneyMovementsByMonth(t *testing.T) {
	t.Run(`success`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/months/2024-01/money_movements"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movements": [
      {
        "id": "mm-123",
        "category_id": "cat-456",
        "category_name": "Groceries",
        "date": "2024-01-15",
        "amount": -15000,
        "payee_id": "payee-789",
        "payee_name": "Grocery Store",
        "recurring_job_id": null,
        "recurring_job_type": null,
        "scheduled_flag": false,
        "approved": true,
        "flag_color": null,
        "transfer_account_id": null,
        "transfer_transaction_id": null,
        "matched_transaction_id": null,
        "import_id": null,
        "type": "withdrawal",
        "isrenamed": false
      }
    ]
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		movements, err := client.MoneyMovement().GetMoneyMovementsByMonth("plan-id-123", "2024-01")
		assert.NoError(t, err)
		assert.Len(t, movements.MoneyMovements, 1)
		assert.Equal(t, "mm-123", movements.MoneyMovements[0].ID)
	})
}

func TestService_GetMoneyMovementGroups(t *testing.T) {
	t.Run(`success`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/money_movement_groups"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movement_groups": [
      {
        "id": "group-123",
        "category_id": "cat-456",
        "category_name": "Groceries",
        "income": false,
        "goal_target": 50000,
        "goal_target_date": "2024-12-31T00:00:00Z",
        "group_created_at": "2024-01-10T12:00:00Z",
        "month": "2024-01-01",
        "note": "Groceries adjustments",
        "performed_by_user_id": "user-123",
        "money_movements": [
          {
            "id": "mm-123",
            "category_id": "cat-456",
            "category_name": "Groceries",
            "date": "2024-01-15",
            "amount": -15000,
            "payee_id": "payee-789",
            "payee_name": "Grocery Store",
            "recurring_job_id": null,
            "recurring_job_type": null,
            "scheduled_flag": false,
            "approved": true,
            "flag_color": null,
            "transfer_account_id": null,
            "transfer_transaction_id": null,
            "matched_transaction_id": null,
            "import_id": null,
            "type": "withdrawal",
            "isrenamed": false
          }
        ]
      }
    ],
    "server_knowledge": 0
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		groups, err := client.MoneyMovement().GetMoneyMovementGroups("plan-id-123")
		assert.NoError(t, err)
		assert.Len(t, groups.MoneyMovementGroups, 1)

		assert.Equal(t, "group-123", groups.MoneyMovementGroups[0].ID)
		assert.Equal(t, "cat-456", groups.MoneyMovementGroups[0].CategoryID)
		assert.Equal(t, "Groceries", groups.MoneyMovementGroups[0].CategoryName)
		assert.False(t, groups.MoneyMovementGroups[0].Income)
		assert.Equal(t, int64(50000), *groups.MoneyMovementGroups[0].GoalTarget)

		expectedDate, _ := time.Parse(time.RFC3339, "2024-12-31T00:00:00Z")
		assert.Equal(t, &expectedDate, groups.MoneyMovementGroups[0].GoalTargetDate)
	})

	t.Run(`success empty`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/money_movement_groups"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movement_groups": []
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		groups, err := client.MoneyMovement().GetMoneyMovementGroups("plan-id-123")
		assert.NoError(t, err)
		assert.Len(t, groups, 0)
	})
}

func TestService_GetMoneyMovementGroupsByMonth(t *testing.T) {
	t.Run(`success`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/months/2024-01/money_movement_groups"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movement_groups": [
      {
        "id": "group-123",
        "category_id": "cat-456",
        "category_name": "Groceries",
        "income": false,
        "goal_target": null,
        "goal_target_date": null,
        "goal_underfunded": null,
        "goal_overspent": null,
        "money_movements": []
      }
    ]
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		groups, err := client.MoneyMovement().GetMoneyMovementGroupsByMonth("plan-id-123", "2024-01")
		assert.NoError(t, err)
		assert.Len(t, groups, 1)
		assert.Equal(t, "group-123", groups.MoneyMovementGroups[0].ID)
	})
}
