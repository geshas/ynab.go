package plan

import (
	"fmt"

	"github.com/geshas/ynab.go/api"
)

// NewService facilitates the creation of a new plan service instance
func NewService(c api.ClientReader) *Service {
	return &Service{c}
}

// Service wraps YNAB plan API endpoints
type Service struct {
	c api.ClientReader
}

// GetPlans fetches the list of plans
// https://api.ynab.com/v1#/Plans/getPlans
func (s *Service) GetPlans() ([]*Summary, error) {
	return s.GetPlansWithAccounts(false)
}

// GetPlansWithAccounts fetches the list of plans
// with optional account information included
// https://api.ynab.com/v1#/Plans/getPlans
func (s *Service) GetPlansWithAccounts(includeAccounts bool) ([]*Summary, error) {
	resModel := struct {
		Data struct {
			Plans []*Summary `json:"plans"`
		} `json:"data"`
	}{}

	url := "/plans"
	if includeAccounts {
		url = "/plans?include_accounts=true"
	}

	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.Plans, nil
}

// GetPlan fetches a single plan with all related entities,
// effectively a full plan export with filtering capabilities
// https://api.ynab.com/v1#/Plans/getPlanById
func (s *Service) GetPlan(planID string, f *api.Filter) (*Snapshot, error) {
	resModel := struct {
		Data struct {
			Plan            *Plan  `json:"plan"`
			ServerKnowledge uint64 `json:"server_knowledge"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/plans/%s", planID)
	if f != nil {
		url = fmt.Sprintf("%s?%s", url, f.ToQuery())
	}

	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}

	return &Snapshot{
		Plan:            resModel.Data.Plan,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// GetLastUsedPlan fetches the last used plan with all related entities
// https://api.ynab.com/v1#/Plans/getPlanById
func (s *Service) GetLastUsedPlan(f *api.Filter) (*Snapshot, error) {
	const lastUsedPlanID = "last-used"
	return s.GetPlan(lastUsedPlanID, f)
}

// GetPlanSettings fetches a plan settings
// https://api.ynab.com/v1#/Plans/getPlanSettingsById
func (s *Service) GetPlanSettings(planID string) (*Settings, error) {
	resModel := struct {
		Data struct {
			Settings *Settings `json:"settings"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/plans/%s/settings", planID)
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}

	return resModel.Data.Settings, nil
}
