package customer_purchase

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/customer_purchase"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"github.com/gilabs/crm-healthcare/api/pkg/util/currency"
)

type Service struct {
	repo        interfaces.CustomerPurchaseHistoryRepository
	accountRepo interfaces.AccountRepository
}

func NewService(repo interfaces.CustomerPurchaseHistoryRepository, accountRepo interfaces.AccountRepository) *Service {
	return &Service{
		repo:        repo,
		accountRepo: accountRepo,
	}
}

func (s *Service) GetByAccountID(accountID string) ([]customer_purchase.CustomerPurchaseHistory, error) {
	return s.repo.FindByAccountID(accountID)
}

func (s *Service) GetAnalytics(accountID string) (map[string]interface{}, error) {
	return s.repo.GetAnalytics(accountID)
}

func (s *Service) GetProductAnalytics(accountID string) ([]customer_purchase.CustomerProductAnalytics, error) {
	analytics, err := s.repo.GetProductAnalytics(accountID)
	if err != nil {
		return nil, err
	}

	for i := range analytics {
		analytics[i].TotalAmountFormatted = currency.FormatCurrency(analytics[i].TotalAmountPurchased)
	}

	return analytics, nil
}

func (s *Service) GetSummary(accountID string) (*customer_purchase.CustomerPurchaseSummaryResponse, error) {
	summary, err := s.repo.GetSummary(accountID)
	if err != nil {
		return nil, err
	}

	summary.TotalAmountFormatted = currency.FormatCurrency(summary.TotalAmount)
	summary.AveragePurchaseAmountFormatted = currency.FormatCurrency(summary.AveragePurchaseAmount)

	return summary, nil
}

func (s *Service) RecordPurchase(history *customer_purchase.CustomerPurchaseHistory) error {
	return s.repo.Create(history)
}
