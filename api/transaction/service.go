package transaction

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/geshas/ynab.go/api"
)

// NewService facilitates the creation of a new transaction service instance
func NewService(c api.ClientReaderWriter) *Service {
	return &Service{c}
}

// Service wraps YNAB transaction API endpoints
type Service struct {
	c api.ClientReaderWriter
}

// SearchResultSnapshot represents the result of a search with server knowledge
type SearchResultSnapshot struct {
	Transactions    []*Transaction
	ServerKnowledge uint64
}

// GetTransactions fetches the list of transactions from
// a budget with filtering capabilities
// https://api.ynab.com/v1#/Transactions/getTransactions
func (s *Service) GetTransactions(planID string, f *Filter) (*SearchResultSnapshot, error) {
	resModel := struct {
		Data struct {
			Transactions    []*Transaction `json:"transactions"`
			ServerKnowledge uint64         `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/transactions", url.PathEscape(planID))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}

	return &SearchResultSnapshot{
		Transactions:    resModel.Data.Transactions,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// GetTransaction fetches a specific transaction from a plan
// https://api.ynab.com/v1#/Transactions/getTransactionById
func (s *Service) GetTransaction(planID, transactionID string) (*Transaction, error) {
	resModel := struct {
		Data struct {
			Transaction *Transaction `json:"transaction"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/transactions/%s", url.PathEscape(planID), url.PathEscape(transactionID))
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.Transaction, nil
}

// CreateTransaction creates a new transaction for a plan
// https://api.ynab.com/v1#/Transactions/createTransaction
func (s *Service) CreateTransaction(planID string,
	p PayloadTransaction) (*OperationSummary, error) {

	return s.CreateTransactions(planID, []PayloadTransaction{p})
}

// CreateTransactions creates one or more new transactions for a plan
// https://api.ynab.com/v1#/Transactions/createTransaction
func (s *Service) CreateTransactions(planID string,
	p []PayloadTransaction) (*OperationSummary, error) {

	payload := struct {
		Transactions []PayloadTransaction `json:"transactions"`
	}{
		p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data *OperationSummary `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/transactions", url.PathEscape(planID))
	err = s.c.POST(reqURL, &resModel, buf)
	if err != nil {
		return nil, err
	}
	return resModel.Data, nil
}

// BulkCreateTransactions creates multiple transactions for a plan
// https://api.ynab.com/v1#/Transactions/bulkCreateTransactions
// Deprecated: This uses a legacy endpoint removed from the current OpenAPI spec.
// Use transaction.CreateTransactions instead.
func (s *Service) BulkCreateTransactions(planID string,
	ps []PayloadTransaction) (*Bulk, error) {

	payload := struct {
		Transactions []PayloadTransaction `json:"transactions"`
	}{
		ps,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			Bulk *Bulk `json:"bulk"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/transactions/bulk", url.PathEscape(planID))
	if err := s.c.POST(reqURL, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.Bulk, nil
}

// UpdateTransaction updates a whole transaction for a replacement
// https://api.ynab.com/v1#/Transactions/updateTransaction
func (s *Service) UpdateTransaction(planID, transactionID string,
	p PayloadTransaction) (*Transaction, error) {

	payload := struct {
		Transaction *PayloadTransaction `json:"transaction"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			Transaction *Transaction `json:"transaction"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/transactions/%s", url.PathEscape(planID), url.PathEscape(transactionID))
	if err := s.c.PUT(reqURL, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.Transaction, nil
}

// UpdateTransactions creates one or more new transactions for a plan
// https://api.ynab.com/v1#/Transactions/updateTransactions
func (s *Service) UpdateTransactions(planID string,
	p []PayloadTransaction) (*OperationSummary, error) {

	payload := struct {
		Transactions []PayloadTransaction `json:"transactions"`
	}{
		p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data *OperationSummary `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/transactions", url.PathEscape(planID))
	err = s.c.PATCH(reqURL, &resModel, buf)
	if err != nil {
		return nil, err
	}
	return resModel.Data, nil
}

// DeleteTransaction deletes a transaction from a plan
// https://api.ynab.com/v1#/Transactions/deleteTransaction
func (s *Service) DeleteTransaction(planID, transactionID string) (*Transaction, error) {
	resModel := struct {
		Data struct {
			Transaction *Transaction `json:"transaction"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/transactions/%s", url.PathEscape(planID), url.PathEscape(transactionID))
	err := s.c.DELETE(reqURL, &resModel)
	if err != nil {
		return nil, err
	}
	return resModel.Data.Transaction, nil
}

// GetTransactionsByAccount fetches the list of transactions of a specific account
// from a plan with filtering capabilities
// https://api.ynab.com/v1#/Transactions/getTransactionsByAccount
func (s *Service) GetTransactionsByAccount(planID, accountID string,
	f *Filter) (*SearchResultSnapshot, error) {

	resModel := struct {
		Data struct {
			Transactions    []*Transaction `json:"transactions"`
			ServerKnowledge uint64         `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/accounts/%s/transactions", url.PathEscape(planID), url.PathEscape(accountID))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}

	return &SearchResultSnapshot{
		Transactions:    resModel.Data.Transactions,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// GetTransactionsByMonth fetches the list of transactions for a specific month from a plan
// https://api.ynab.com/v1#/Transactions/getTransactionsByMonth
func (s *Service) GetTransactionsByMonth(planID, month string, f *Filter) (*SearchResultSnapshot, error) {
	resModel := struct {
		Data struct {
			Transactions    []*Transaction `json:"transactions"`
			ServerKnowledge uint64         `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/months/%s/transactions", url.PathEscape(planID), url.PathEscape(month))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}

	return &SearchResultSnapshot{
		Transactions:    resModel.Data.Transactions,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// GetTransactionsByCategory fetches the list of transactions of a specific category
// from a plan with filtering capabilities
// https://api.ynab.com/v1#/Transactions/getTransactionsByCategory
func (s *Service) GetTransactionsByCategory(planID, categoryID string,
	f *Filter) ([]*Hybrid, error) {
	snapshot, err := s.GetTransactionsByCategoryWithSnapshot(planID, categoryID, f)
	if err != nil {
		return nil, err
	}

	return snapshot.Transactions, nil
}

// GetTransactionsByCategoryWithSnapshot fetches the list of transactions of a specific category
// from a plan with filtering capabilities and server knowledge.
// https://api.ynab.com/v1#/Transactions/getTransactionsByCategory
func (s *Service) GetTransactionsByCategoryWithSnapshot(planID, categoryID string,
	f *Filter) (*HybridSearchResultSnapshot, error) {

	resModel := struct {
		Data struct {
			Transactions    []*Hybrid `json:"transactions"`
			ServerKnowledge uint64    `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/categories/%s/transactions", url.PathEscape(planID), url.PathEscape(categoryID))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}

	return &HybridSearchResultSnapshot{
		Transactions:    resModel.Data.Transactions,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// GetTransactionsByPayee fetches the list of transactions of a specific payee
// from a plan with filtering capabilities
// https://api.ynab.com/v1#/Transactions/getTransactionsByPayee
func (s *Service) GetTransactionsByPayee(planID, payeeID string,
	f *Filter) ([]*Hybrid, error) {
	snapshot, err := s.GetTransactionsByPayeeWithSnapshot(planID, payeeID, f)
	if err != nil {
		return nil, err
	}

	return snapshot.Transactions, nil
}

// GetTransactionsByPayeeWithSnapshot fetches the list of transactions of a specific payee
// from a plan with filtering capabilities and server knowledge.
// https://api.ynab.com/v1#/Transactions/getTransactionsByPayee
func (s *Service) GetTransactionsByPayeeWithSnapshot(planID, payeeID string,
	f *Filter) (*HybridSearchResultSnapshot, error) {

	resModel := struct {
		Data struct {
			Transactions    []*Hybrid `json:"transactions"`
			ServerKnowledge uint64    `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/payees/%s/transactions", url.PathEscape(planID), url.PathEscape(payeeID))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}

	return &HybridSearchResultSnapshot{
		Transactions:    resModel.Data.Transactions,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}, nil
}

// ScheduledSearchResultSnapshot represents the result of a scheduled transaction search with server knowledge
type ScheduledSearchResultSnapshot struct {
	ScheduledTransactions []*Scheduled
	ServerKnowledge       uint64
}

// HybridSearchResultSnapshot represents the result of a hybrid transaction search with server knowledge.
type HybridSearchResultSnapshot struct {
	Transactions    []*Hybrid
	ServerKnowledge uint64
}

// GetScheduledTransactions fetches the list of scheduled transactions from
// a budget with filtering capabilities
// https://api.ynab.com/v1#/Scheduled_Transactions/getScheduledTransactions
func (s *Service) GetScheduledTransactions(planID string, f *api.Filter) (*ScheduledSearchResultSnapshot, error) {
	resModel := struct {
		Data struct {
			ScheduledTransactions []*Scheduled `json:"scheduled_transactions"`
			ServerKnowledge       uint64       `json:"server_knowledge"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/scheduled_transactions", url.PathEscape(planID))
	if f != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, f.ToQuery())
	}

	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}

	return &ScheduledSearchResultSnapshot{
		ScheduledTransactions: resModel.Data.ScheduledTransactions,
		ServerKnowledge:       resModel.Data.ServerKnowledge,
	}, nil
}

// GetScheduledTransaction fetches a specific scheduled transaction from a plan
// https://api.ynab.com/v1#/Scheduled_Transactions/getScheduledTransactionById
func (s *Service) GetScheduledTransaction(planID, scheduledTransactionID string) (*Scheduled, error) {
	resModel := struct {
		Data struct {
			ScheduledTransaction *Scheduled `json:"scheduled_transaction"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/scheduled_transactions/%s", url.PathEscape(planID), url.PathEscape(scheduledTransactionID))
	if err := s.c.GET(reqURL, &resModel); err != nil {
		return nil, err
	}
	return resModel.Data.ScheduledTransaction, nil
}

// Filter represents the optional filter while fetching transactions
type Filter struct {
	Since                 *api.Date
	Type                  *Status
	LastKnowledgeOfServer *uint64
}

// ToQuery returns the filters as a HTTP query string
func (f *Filter) ToQuery() string {
	v := url.Values{}
	if f.Since != nil && !f.Since.IsZero() {
		v.Set("since_date", api.DateFormat(*f.Since))
	}
	if f.Type != nil {
		v.Set("type", string(*f.Type))
	}
	if f.LastKnowledgeOfServer != nil {
		v.Set("last_knowledge_of_server", fmt.Sprintf("%d", *f.LastKnowledgeOfServer))
	}
	return v.Encode()
}

// CreateScheduledTransaction creates a new scheduled transaction for a plan
// https://api.ynab.com/v1#/Scheduled_Transactions/createScheduledTransaction
func (s *Service) CreateScheduledTransaction(planID string, p PayloadScheduledTransaction) (*Scheduled, error) {
	payload := struct {
		ScheduledTransaction *PayloadScheduledTransaction `json:"scheduled_transaction"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			ScheduledTransaction *Scheduled `json:"scheduled_transaction"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/scheduled_transactions", url.PathEscape(planID))
	if err := s.c.POST(reqURL, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.ScheduledTransaction, nil
}

// UpdateScheduledTransaction updates a scheduled transaction for a plan
// https://api.ynab.com/v1#/Scheduled_Transactions/updateScheduledTransaction
func (s *Service) UpdateScheduledTransaction(planID, scheduledTransactionID string, p PayloadScheduledTransaction) (*Scheduled, error) {
	payload := struct {
		ScheduledTransaction *PayloadScheduledTransaction `json:"scheduled_transaction"`
	}{
		&p,
	}

	buf, err := json.Marshal(&payload)
	if err != nil {
		return nil, err
	}

	resModel := struct {
		Data struct {
			ScheduledTransaction *Scheduled `json:"scheduled_transaction"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/scheduled_transactions/%s", url.PathEscape(planID), url.PathEscape(scheduledTransactionID))
	if err := s.c.PUT(reqURL, &resModel, buf); err != nil {
		return nil, err
	}
	return resModel.Data.ScheduledTransaction, nil
}

// DeleteScheduledTransaction deletes a scheduled transaction from a plan
// https://api.ynab.com/v1#/Scheduled_Transactions/deleteScheduledTransaction
func (s *Service) DeleteScheduledTransaction(planID, scheduledTransactionID string) (*Scheduled, error) {
	resModel := struct {
		Data struct {
			ScheduledTransaction *Scheduled `json:"scheduled_transaction"`
		} `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/scheduled_transactions/%s", url.PathEscape(planID), url.PathEscape(scheduledTransactionID))
	err := s.c.DELETE(reqURL, &resModel)
	if err != nil {
		return nil, err
	}
	return resModel.Data.ScheduledTransaction, nil
}

// ImportTransactions imports available transactions from all linked accounts for a plan
// https://api.ynab.com/v1#/Transactions/importTransactions
func (s *Service) ImportTransactions(planID string) (*ImportResult, error) {
	resModel := struct {
		Data *ImportResult `json:"data"`
	}{}

	reqURL := fmt.Sprintf("/plans/%s/transactions/import", url.PathEscape(planID))
	if err := s.c.POST(reqURL, &resModel, nil); err != nil {
		return nil, err
	}
	return resModel.Data, nil
}
