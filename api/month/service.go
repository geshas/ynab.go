package month

import (
	"fmt"
	"net/url"

	"github.com/geshas/ynab.go/api"
)

// NewService facilitates the creation of a new month service instance
func NewService(c api.ClientReader) *Service {
	return &Service{c}
}

// Service wraps YNAB month API endpoints
type Service struct {
	c api.ClientReader
}

// GetMonths fetches the list of months from a plan
// https://api.ynab.com/v1#/Months/getBudgetMonths
func (s *Service) GetMonths(planID string, f *api.Filter) (*SearchResultSnapshot, error) {
	resModel := struct {
		Data struct {
			Months          []*Summary `json:"months"`
			ServerKnowledge uint64     `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/months", url.PathEscape(planID))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return &SearchResultSnapshot{
		Months:          resModel.Data.Months,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// GetMonth fetches a specific month from a plan
// https://api.ynab.com/v1#/Months/getBudgetMonth
func (s *Service) GetMonth(planID string, month api.Date) (*Month, error) {
	resModel := struct {
		Data struct {
			Month *Month `json:"month"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/months/%s",
		url.PathEscape(planID), url.PathEscape(api.DateFormat(month)))
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.Month, nil
}
