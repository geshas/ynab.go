package category

import (
	"encoding/json"
	"fmt"

	"github.com/geshas/ynab.go/api"
)

const currentMonthID = "current"

// NewService facilitates the creation of a new category service instance
func NewService(c api.ClientReaderWriter) *Service {
	return &Service{c}
}

// Service wraps YNAB category API endpoints
type Service struct {
	c api.ClientReaderWriter
}

// GetCategories fetches the list of category groups for a budget
// https://api.ynab.com/v1#/Categories/getCategories
func (s *Service) GetCategories(budgetID string, f *api.Filter) (*SearchResultSnapshot, error) {
	resModel := struct {
		Data struct {
			CategoryGroups  []*GroupWithCategories `json:"category_groups"`
			ServerKnowledge uint64                 `json:"server_knowledge"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/budgets/%s/categories", budgetID)
	if f != nil {
		url = fmt.Sprintf("%s?%s", url, f.ToQuery())
	}
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}

	return &SearchResultSnapshot{
		GroupWithCategories: resModel.Data.CategoryGroups,
		ServerKnowledge:     resModel.Data.ServerKnowledge,
	}, nil
}

// GetCategory fetches a specific category from a budget
// https://api.ynab.com/v1#/Categories/getCategoryById
func (s *Service) GetCategory(budgetID, categoryID string) (*Category, error) {
	resModel := struct {
		Data struct {
			Category *Category `json:"category"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/budgets/%s/categories/%s", budgetID, categoryID)
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.Category, nil
}

// GetCategoryForMonth fetches a specific category from a budget month
// https://api.ynab.com/v1#/Categories/getMonthCategoryById
func (s *Service) GetCategoryForMonth(budgetID, categoryID string,
	month api.Date) (*Category, error) {

	return s.getCategoryForMonth(budgetID, categoryID, api.DateFormat(month))
}

// GetCategoryForCurrentMonth fetches a specific category from the current budget month
// https://api.ynab.com/v1#/Categories/getMonthCategoryById
func (s *Service) GetCategoryForCurrentMonth(budgetID, categoryID string) (*Category, error) {
	return s.getCategoryForMonth(budgetID, categoryID, currentMonthID)
}

func (s *Service) getCategoryForMonth(budgetID, categoryID, month string) (*Category, error) {
	resModel := struct {
		Data struct {
			Category *Category `json:"category"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/budgets/%s/months/%s/categories/%s", budgetID, month, categoryID)
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.Category, nil
}

// UpdateCategoryForMonth updates a category for a month
// https://api.ynab.com/v1#/Categories/updateMonthCategory
func (s *Service) UpdateCategoryForMonth(budgetID, categoryID string, month api.Date,
	p PayloadMonthCategory) (*Category, error) {

	return s.updateCategoryForMonth(budgetID, categoryID, api.DateFormat(month), p)
}

// UpdateCategoryForCurrentMonth updates a category for the current month
// https://api.ynab.com/v1#/Categories/updateMonthCategory
func (s *Service) UpdateCategoryForCurrentMonth(budgetID, categoryID string,
	p PayloadMonthCategory) (*Category, error) {

	return s.updateCategoryForMonth(budgetID, categoryID, currentMonthID, p)
}

func (s *Service) updateCategoryForMonth(budgetID, categoryID, month string,
	p PayloadMonthCategory) (*Category, error) {

	payload := struct {
		Category *PayloadMonthCategory `json:"category"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			Category *Category `json:"category"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/budgets/%s/months/%s/categories/%s", budgetID,
		month, categoryID)

	if err := s.c.PATCH(url, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.Category, nil
}

// UpdateCategory updates a category
// https://api.ynab.com/v1#/Categories/updateCategory
func (s *Service) UpdateCategory(budgetID, categoryID string, p PayloadCategory) (*Category, error) {
	payload := struct {
		Category *PayloadCategory `json:"category"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			Category *Category `json:"category"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/budgets/%s/categories/%s", budgetID, categoryID)

	if err := s.c.PATCH(url, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.Category, nil
}

// CreateCategory creates a new category
// https://api.ynab.com/v1#/Categories/createCategory
func (s *Service) CreateCategory(planID string, p PayloadCreateCategory) (*Category, error) {
	payload := struct {
		Category *PayloadCreateCategory `json:"category"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			Category *Category `json:"category"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/plans/%s/categories", planID)

	if err := s.c.POST(url, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.Category, nil
}

// CreateCategoryGroup creates a new category group
// https://api.ynab.com/v1#/Categories/createCategoryGroup
func (s *Service) CreateCategoryGroup(planID string, p PayloadCreateCategoryGroup) (*Group, error) {
	payload := struct {
		CategoryGroup *PayloadCreateCategoryGroup `json:"category_group"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			CategoryGroup *Group `json:"category_group"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/plans/%s/category_groups", planID)

	if err := s.c.POST(url, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.CategoryGroup, nil
}

// UpdateCategoryGroup updates a category group
// https://api.ynab.com/v1#/Categories/updateCategoryGroup
func (s *Service) UpdateCategoryGroup(planID, categoryGroupID string, p PayloadUpdateCategoryGroup) (*Group, error) {
	payload := struct {
		CategoryGroup *PayloadUpdateCategoryGroup `json:"category_group"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			CategoryGroup *Group `json:"category_group"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/plans/%s/category_groups/%s", planID, categoryGroupID)

	if err := s.c.PATCH(url, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.CategoryGroup, nil
}
