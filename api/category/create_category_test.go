package category_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"

	"github.com/geshas/ynab.go"
	"github.com/geshas/ynab.go/api/category"
)

func TestService_CreateCategory(t *testing.T) {
	t.Run(`success`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		note := "Test category note"
		payload := category.PayloadCreateCategory{
			Name:            "New Category",
			CategoryGroupID: "group-123",
			Note:            &note,
		}

		url := "https://api.ynab.com/v1/plans/plan-id-123/categories"
		httpmock.RegisterResponder(http.MethodPost, url,
			func(req *http.Request) (*http.Response, error) {
				resModel := struct {
					Category *category.PayloadCreateCategory `json:"category"`
				}{}
				err := json.NewDecoder(req.Body).Decode(&resModel)
				assert.NoError(t, err)
				assert.Equal(t, &payload, resModel.Category)

				res := httpmock.NewStringResponse(201, `{
  "data": {
    "category": {
      "id": "new-cat-456",
      "category_group_id": "group-123",
      "category_group_name": "Test Group",
      "name": "New Category",
      "hidden": false,
      "original_category_group_id": null,
      "note": "Test category note",
      "budgeted": 0,
      "activity": 0,
      "balance": 0,
      "deleted": false
    }
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		c, err := client.Category().CreateCategory("plan-id-123", payload)
		assert.NoError(t, err)

		expected := &category.Category{
			ID:                "new-cat-456",
			CategoryGroupID:   "group-123",
			CategoryGroupName: "Test Group",
			Name:              "New Category",
			Hidden:            false,
			Budgeted:          0,
			Activity:          0,
			Balance:           0,
			Deleted:           false,
			Note:              &note,
		}
		assert.Equal(t, expected, c)
	})

	t.Run(`error`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		payload := category.PayloadCreateCategory{
			Name:            "New Category",
			CategoryGroupID: "group-123",
		}

		url := "https://api.ynab.com/v1/plans/plan-id-123/categories"
		httpmock.RegisterResponder(http.MethodPost, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(400, `{
  "error": {
    "id": "400",
    "name": "Bad Request",
    "detail": "Invalid category data"
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		c, err := client.Category().CreateCategory("plan-id-123", payload)
		assert.Error(t, err)
		assert.Nil(t, c)
	})
}

func TestService_CreateCategoryGroup(t *testing.T) {
	t.Run(`success`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		payload := category.PayloadCreateCategoryGroup{
			Name: "New Category Group",
		}

		url := "https://api.ynab.com/v1/plans/plan-id-123/category_groups"
		httpmock.RegisterResponder(http.MethodPost, url,
			func(req *http.Request) (*http.Response, error) {
				resModel := struct {
					CategoryGroup *category.PayloadCreateCategoryGroup `json:"category_group"`
				}{}
				err := json.NewDecoder(req.Body).Decode(&resModel)
				assert.NoError(t, err)
				assert.Equal(t, &payload, resModel.CategoryGroup)

				res := httpmock.NewStringResponse(201, `{
  "data": {
    "category_group": {
      "id": "new-group-456",
      "name": "New Category Group",
      "hidden": false,
      "deleted": false
    }
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		g, err := client.Category().CreateCategoryGroup("plan-id-123", payload)
		assert.NoError(t, err)

		expected := &category.Group{
			ID:      "new-group-456",
			Name:    "New Category Group",
			Hidden:  false,
			Deleted: false,
		}
		assert.Equal(t, expected, g)
	})

	t.Run(`error`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		payload := category.PayloadCreateCategoryGroup{
			Name: "New Category Group",
		}

		url := "https://api.ynab.com/v1/plans/plan-id-123/category_groups"
		httpmock.RegisterResponder(http.MethodPost, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(400, `{
  "error": {
    "id": "400",
    "name": "Bad Request",
    "detail": "Invalid category group data"
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		g, err := client.Category().CreateCategoryGroup("plan-id-123", payload)
		assert.Error(t, err)
		assert.Nil(t, g)
	})
}

func TestService_UpdateCategoryGroup(t *testing.T) {
	t.Run(`success`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		newName := "Updated Category Group"
		payload := category.PayloadUpdateCategoryGroup{
			Name: newName,
		}

		url := "https://api.ynab.com/v1/plans/plan-id-123/category_groups/group-456"
		httpmock.RegisterResponder(http.MethodPatch, url,
			func(req *http.Request) (*http.Response, error) {
				resModel := struct {
					CategoryGroup *category.PayloadUpdateCategoryGroup `json:"category_group"`
				}{}
				err := json.NewDecoder(req.Body).Decode(&resModel)
				assert.NoError(t, err)
				assert.Equal(t, &payload, resModel.CategoryGroup)

				res := httpmock.NewStringResponse(200, `{
  "data": {
    "category_group": {
      "id": "group-456",
      "name": "Updated Category Group",
      "hidden": false,
      "deleted": false
    }
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		g, err := client.Category().UpdateCategoryGroup("plan-id-123", "group-456", payload)
		assert.NoError(t, err)

		expected := &category.Group{
			ID:      "group-456",
			Name:    "Updated Category Group",
			Hidden:  false,
			Deleted: false,
		}
		assert.Equal(t, expected, g)
	})

	t.Run(`error`, func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		payload := category.PayloadUpdateCategoryGroup{
			Name: "Updated Category Group",
		}

		url := "https://api.ynab.com/v1/plans/plan-id-123/category_groups/group-456"
		httpmock.RegisterResponder(http.MethodPatch, url,
			func(req *http.Request) (*http.Response, error) {
				res := httpmock.NewStringResponse(404, `{
  "error": {
    "id": "404",
    "name": "Not Found",
    "detail": "Category group not found"
  }
}`)
				return res, nil
			},
		)

		client := ynab.NewClient("")
		g, err := client.Category().UpdateCategoryGroup("plan-id-123", "group-456", payload)
		assert.Error(t, err)
		assert.Nil(t, g)
	})
}
