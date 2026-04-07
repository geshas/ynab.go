package account

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/geshas/ynab.go/api"
)

// NewService facilitates the creation of a new account service instance
func NewService(c api.ClientReaderWriter) *Service {
	return &Service{c}
}

// Service wraps YNAB account API endpoints
type Service struct {
	c api.ClientReaderWriter
}

var errUnsupportedAccountTypeForCreate = errors.New("unsupported account type for create")

// GetAccounts fetches the list of accounts from a plan
// https://api.ynab.com/v1#/Accounts/getAccounts
func (s *Service) GetAccounts(planID string, f *api.Filter) (*SearchResultSnapshot, error) {
	resModel := struct {
		Data struct {
			Accounts        []*Account `json:"accounts"`
			ServerKnowledge uint64     `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/accounts", url.PathEscape(planID))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}

	return &SearchResultSnapshot{
		Accounts:        resModel.Data.Accounts,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// GetAccount fetches a specific account from a plan
// https://api.ynab.com/v1#/Accounts/getAccountById
func (s *Service) GetAccount(planID, accountID string) (*Account, error) {
	resModel := struct {
		Data struct {
			Account *Account `json:"account"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/accounts/%s", url.PathEscape(planID), url.PathEscape(accountID))
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.Account, nil
}

// CreateAccount creates a new account in a plan
// https://api.ynab.com/v1#/Accounts/createAccount
func (s *Service) CreateAccount(planID string, p PayloadAccount) (*Account, error) {
	if err := validateCreateAccountType(p.Type); err != nil {
		return nil, err
	}

	payload := struct {
		Account *PayloadAccount `json:"account"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			Account *Account `json:"account"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/accounts", url.PathEscape(planID))
	if err := s.c.POST(reqURL, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.Account, nil
}

func validateCreateAccountType(accountType Type) error {
	switch accountType {
	case TypeChecking, TypeSavings, TypeCash, TypeCreditCard, TypeOtherAsset, TypeOtherLiability:
		return nil
	default:
		return fmt.Errorf("%w: %q", errUnsupportedAccountTypeForCreate, accountType)
	}
}
