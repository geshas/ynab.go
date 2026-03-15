package money_movement

import (
	"fmt"

	"github.com/geshas/ynab.go/api"
)

// NewService facilitates the creation of a new money movement service instance
func NewService(c api.ClientReader) *Service {
	return &Service{c}
}

// Service wraps YNAB money movement API endpoints
type Service struct {
	c api.ClientReader
}

// MoneyMovementsSnapshot represents a set of money movements and the associated
// server knowledge value returned by the YNAB API.
type MoneyMovementsSnapshot struct {
	MoneyMovements  []*MoneyMovement `json:"money_movements"`
	ServerKnowledge int64            `json:"server_knowledge"`
}

// MoneyMovementGroupsSnapshot represents a set of money movement groups and the
// associated server knowledge value returned by the YNAB API.
type MoneyMovementGroupsSnapshot struct {
	MoneyMovementGroups []*MoneyMovementGroup `json:"money_movement_groups"`
	ServerKnowledge     int64                 `json:"server_knowledge"`
}

// GetMoneyMovements fetches all money movements for a plan
// https://api.ynab.com/v1#/Money%20Movements/getMoneyMovements
func (s *Service) GetMoneyMovements(planID string) (*MoneyMovementsSnapshot, error) {
	resModel := struct {
		Data struct {
			MoneyMovements  []*MoneyMovement `json:"money_movements"`
			ServerKnowledge int64            `json:"server_knowledge"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/plans/%s/money_movements", planID)
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}

	snapshot := &MoneyMovementsSnapshot{
		MoneyMovements:  resModel.Data.MoneyMovements,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}

	return snapshot, nil
}

// GetMoneyMovementsByMonth fetches money movements for a specific plan month
// https://api.ynab.com/v1#/Money%20Movements/getMoneyMovementsByMonth
func (s *Service) GetMoneyMovementsByMonth(planID string, month string) (*MoneyMovementsSnapshot, error) {
	resModel := struct {
		Data struct {
			MoneyMovements  []*MoneyMovement `json:"money_movements"`
			ServerKnowledge int64            `json:"server_knowledge"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/plans/%s/months/%s/money_movements", planID, month)
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}

	snapshot := &MoneyMovementsSnapshot{
		MoneyMovements:  resModel.Data.MoneyMovements,
		ServerKnowledge: resModel.Data.ServerKnowledge,
	}

	return snapshot, nil
}

// GetMoneyMovementGroups fetches all money movement groups for a plan
// https://api.ynab.com/v1#/Money%20Movements/getMoneyMovementGroups
func (s *Service) GetMoneyMovementGroups(planID string) (*MoneyMovementGroupsSnapshot, error) {
	resModel := struct {
		Data struct {
			MoneyMovementGroups []*MoneyMovementGroup `json:"money_movement_groups"`
			ServerKnowledge     int64                 `json:"server_knowledge"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/plans/%s/money_movement_groups", planID)
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}

	snapshot := &MoneyMovementGroupsSnapshot{
		MoneyMovementGroups: resModel.Data.MoneyMovementGroups,
		ServerKnowledge:     resModel.Data.ServerKnowledge,
	}

	return snapshot, nil
}

// GetMoneyMovementGroupsByMonth fetches money movement groups for a specific plan month
// https://api.ynab.com/v1#/Money%20Movements/getMoneyMovementGroupsByMonth
func (s *Service) GetMoneyMovementGroupsByMonth(planID string, month string) (*MoneyMovementGroupsSnapshot, error) {
	resModel := struct {
		Data struct {
			MoneyMovementGroups []*MoneyMovementGroup `json:"money_movement_groups"`
			ServerKnowledge     int64                 `json:"server_knowledge"`
		} `json:"data"`
	}{}

	url := fmt.Sprintf("/plans/%s/months/%s/money_movement_groups", planID, month)
	if err := s.c.GET(url, &resModel); err != nil {
		return nil, err
	}

	snapshot := &MoneyMovementGroupsSnapshot{
		MoneyMovementGroups: resModel.Data.MoneyMovementGroups,
		ServerKnowledge:     resModel.Data.ServerKnowledge,
	}

	return snapshot, nil
}
