package money_movement_test

import (
	"net/http"
	"testing"

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
        "month": "2024-01-01",
        "moved_at": "2024-01-15T10:00:00Z",
        "note": "Test note",
        "money_movement_group_id": "group-123",
        "performed_by_user_id": "user-456",
        "from_category_id": "cat-from",
        "to_category_id": "cat-to",
        "amount": -15000
      }
    ]
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		snapshot, err := client.MoneyMovement().GetMoneyMovements("plan-id-123", nil)
		assert.NoError(t, err)
		assert.NotNil(t, snapshot)
		assert.Len(t, snapshot.MoneyMovements, 1)

		expectedDate, _ := api.DateFromString("2024-01-01")
		assert.Equal(t, "mm-123", snapshot.MoneyMovements[0].ID)
		assert.Equal(t, expectedDate, snapshot.MoneyMovements[0].Month)
		assert.Equal(t, int64(-15000), snapshot.MoneyMovements[0].Amount)
	})

	t.Run(`success empty`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/money_movements"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movements": []
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		movements, err := client.MoneyMovement().GetMoneyMovements("plan-id-123", nil)
		assert.NoError(t, err)
		assert.Len(t, movements.MoneyMovements, 0)
	})

	t.Run(`success with filter`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/money_movements?last_knowledge_of_server=42"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movements": [],
    "server_knowledge": 42
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		filter := &api.Filter{LastKnowledgeOfServer: 42}
		snapshot, err := client.MoneyMovement().GetMoneyMovements("plan-id-123", filter)
		assert.NoError(t, err)
		assert.NotNil(t, snapshot)
		assert.Equal(t, int64(42), snapshot.ServerKnowledge)
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
        "month": "2024-01-01",
        "moved_at": "2024-01-15T10:00:00Z",
        "note": "Test note",
        "money_movement_group_id": "group-123",
        "performed_by_user_id": "user-456",
        "from_category_id": "cat-from",
        "to_category_id": "cat-to",
        "amount": -15000
      }
    ]
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		snapshot, err := client.MoneyMovement().GetMoneyMovementsByMonth("plan-id-123", "2024-01", nil)
		assert.NoError(t, err)
		assert.Len(t, snapshot.MoneyMovements, 1)
		assert.Equal(t, "mm-123", snapshot.MoneyMovements[0].ID)
	})

	t.Run(`success with filter`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/months/2024-01/money_movements?last_knowledge_of_server=99"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movements": [],
    "server_knowledge": 99
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		filter := &api.Filter{LastKnowledgeOfServer: 99}
		snapshot, err := client.MoneyMovement().GetMoneyMovementsByMonth("plan-id-123", "2024-01", filter)
		assert.NoError(t, err)
		assert.NotNil(t, snapshot)
		assert.Equal(t, int64(99), snapshot.ServerKnowledge)
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
        "group_created_at": "2024-01-10T12:00:00Z",
        "month": "2024-01-01",
        "note": "Test note",
        "performed_by_user_id": "user-456",
        "money_movements": [
          {
            "id": "mm-123",
            "month": "2024-01-01",
            "moved_at": "2024-01-15T10:00:00Z",
            "note": "Inner note",
            "money_movement_group_id": "group-123",
            "performed_by_user_id": "user-456",
            "from_category_id": "cat-from",
            "to_category_id": "cat-to",
            "amount": -15000
          }
        ]
      }
    ]
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		snapshot, err := client.MoneyMovement().GetMoneyMovementGroups("plan-id-123", nil)
		assert.NoError(t, err)
		assert.Len(t, snapshot.MoneyMovementGroups, 1)

		assert.Equal(t, "group-123", snapshot.MoneyMovementGroups[0].ID)
		assert.Equal(t, "Test note", snapshot.MoneyMovementGroups[0].Note)
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
		snapshot, err := client.MoneyMovement().GetMoneyMovementGroups("plan-id-123", nil)
		assert.NoError(t, err)
		assert.Len(t, snapshot.MoneyMovementGroups, 0)
	})

	t.Run(`success with filter`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/money_movement_groups?last_knowledge_of_server=7"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movement_groups": [],
    "server_knowledge": 7
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		filter := &api.Filter{LastKnowledgeOfServer: 7}
		snapshot, err := client.MoneyMovement().GetMoneyMovementGroups("plan-id-123", filter)
		assert.NoError(t, err)
		assert.NotNil(t, snapshot)
		assert.Equal(t, int64(7), snapshot.ServerKnowledge)
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
		snapshot, err := client.MoneyMovement().GetMoneyMovementGroupsByMonth("plan-id-123", "2024-01", nil)
		assert.NoError(t, err)
		assert.Len(t, snapshot.MoneyMovementGroups, 1)
		assert.Equal(t, "group-123", snapshot.MoneyMovementGroups[0].ID)
	})

	t.Run(`success with filter`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		url := "https://api.ynab.com/v1/plans/plan-id-123/months/2024-01/money_movement_groups?last_knowledge_of_server=8"
		httpmock.RegisterResponder(http.MethodGet, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(200, `{
  "data": {
    "money_movement_groups": [],
    "server_knowledge": 8
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		filter := &api.Filter{LastKnowledgeOfServer: 8}
		snapshot, err := client.MoneyMovement().GetMoneyMovementGroupsByMonth("plan-id-123", "2024-01", filter)
		assert.NoError(t, err)
		assert.NotNil(t, snapshot)
		assert.Equal(t, int64(8), snapshot.ServerKnowledge)
	})
}
