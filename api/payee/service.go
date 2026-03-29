package payee

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/geshas/ynab.go/api"
)

// NewService facilitates the creation of a new payee service instance
func NewService(c api.ClientReaderWriter) *Service {
	return &Service{c}
}

// Service wraps YNAB payee API endpoints
type Service struct {
	c api.ClientReaderWriter
}

// GetPayees fetches the list of payees from a plan
// https://api.ynab.com/v1#/Payees/getPayees
func (s *Service) GetPayees(planID string, f *api.Filter) (*SearchResultSnapshot, error) {
	resModel := struct {
		Data struct {
			Payees          []*Payee `json:"payees"`
			ServerKnowledge uint64   `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/payees", url.PathEscape(planID))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return &SearchResultSnapshot{
		Payees:          resModel.Data.Payees,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// CreatePayee creates a new payee for a plan
// https://api.ynab.com/v1#/Payees/createPayee
func (s *Service) CreatePayee(planID string, p PayloadPayee) (*Payee, error) {
	payload := struct {
		Payee *PayloadPayee `json:"payee"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			Payee *Payee `json:"payee"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/payees", url.PathEscape(planID))
	if err := s.c.POST(reqURL, &resModel, buf); err != nil {
		return nil, err
	}

	return resModel.Data.Payee, nil
}

// GetPayee fetches a specific payee from a plan
// https://api.ynab.com/v1#/Payees/getPayeeById
func (s *Service) GetPayee(planID, payeeID string) (*Payee, error) {
	resModel := struct {
		Data struct {
			Payee *Payee `json:"payee"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/payees/%s", url.PathEscape(planID), url.PathEscape(payeeID))
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.Payee, nil
}

// GetPayeeLocations fetches the list of payee locations from a plan
// https://api.ynab.com/v1#/Payee_Locations/getPayeeLocations
func (s *Service) GetPayeeLocations(planID string) ([]*Location, error) {
	resModel := struct {
		Data struct {
			PayeeLocations []*Location `json:"payee_locations"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/payee_locations", url.PathEscape(planID))
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.PayeeLocations, nil
}

// GetPayeeLocation fetches a specific payee location from a plan
// https://api.ynab.com/v1#/Payee_Locations/getPayeeLocationById
func (s *Service) GetPayeeLocation(planID, payeeLocationID string) (*Location, error) {
	resModel := struct {
		Data struct {
			PayeeLocation *Location `json:"payee_location"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/payee_locations/%s", url.PathEscape(planID), url.PathEscape(payeeLocationID))
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.PayeeLocation, nil
}

// GetPayeeLocationsByPayee fetches the list of locations of a specific payee from a plan
// https://api.ynab.com/v1#/Payee_Locations/getPayeeLocationsByPayee
func (s *Service) GetPayeeLocationsByPayee(planID, payeeID string) ([]*Location, error) {
	resModel := struct {
		Data struct {
			PayeeLocations []*Location `json:"payee_locations"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/payees/%s/payee_locations", url.PathEscape(planID), url.PathEscape(payeeID))
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.PayeeLocations, nil
}

// UpdatePayee updates a payee for a plan
// https://api.ynab.com/v1#/Payees/updatePayee
func (s *Service) UpdatePayee(planID, payeeID string, p PayloadPayee) (*Payee, error) {
	payload := struct {
		Payee *PayloadPayee `json:"payee"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			Payee           *Payee `json:"payee"`
			ServerKnowledge uint64 `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/payees/%s", url.PathEscape(planID), url.PathEscape(payeeID))
	if err := s.c.PATCH(reqURL, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.Payee, nil
}
