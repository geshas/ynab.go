package plan

import (
	"fmt"
	"net/url"

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
	result, err := s.GetPlansWithAccountsDetailed(false)
	if err != nil {
		return nil, err
	}

	return result.Plans, nil
}

// GetPlansWithAccounts fetches the list of plans
// with optional account information included
// https://api.ynab.com/v1#/Plans/getPlans
func (s *Service) GetPlansWithAccounts(includeAccounts bool) ([]*Summary, error) {
	result, err := s.GetPlansWithAccountsDetailed(includeAccounts)
	if err != nil {
		return nil, err
	}

	return result.Plans, nil
}

// GetPlansDetailed fetches the list of plans and optional default plan.
// https://api.ynab.com/v1#/Plans/getPlans
func (s *Service) GetPlansDetailed() (*PlansResult, error) {
	return s.GetPlansWithAccountsDetailed(false)
}

// GetPlansWithAccountsDetailed fetches the list of plans, optional default plan,
// and optional account information included.
// https://api.ynab.com/v1#/Plans/getPlans
func (s *Service) GetPlansWithAccountsDetailed(includeAccounts bool) (*PlansResult, error) {
	resModel := struct {
		Data struct {
			Plans       []*Summary `json:"plans"`
			DefaultPlan *Summary   `json:"default_plan"`
		} `json:"data"`
	}{}

	reqURL := "/plans"
	if includeAccounts {
		reqURL = "/plans?include_accounts=true"
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}

	return &PlansResult{
		Plans:       resModel.Data.Plans,
		DefaultPlan: resModel.Data.DefaultPlan,
	}, nil
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

	reqURL := fmt.Sprintf("/plans/%s", url.PathEscape(planID))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
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

	reqURL := fmt.Sprintf("/plans/%s/settings", url.PathEscape(planID))
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}

	return resModel.Data.Settings, nil
}
