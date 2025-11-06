package service

import (
	"wallet/config"
	"wallet/models"
)

// TransactionServiceImpl implements transaction service interfaces
type TransactionServiceImpl struct{}

// NewTransactionService creates transaction service instance
func NewTransactionService() *TransactionServiceImpl {
	return &TransactionServiceImpl{}
}

// GetUserTransactions retrieves user's transaction history
func (s *TransactionServiceImpl) GetUserTransactions(userID, page, limit int) ([]models.Transaction, error) {
	offset := (page - 1) * limit
	var transactions []models.Transaction

	result := config.GetDB().Where(
		"from_user_id = ? OR to_user_id = ?", userID, userID,
	).Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions)

	if result.Error != nil {
		return nil, result.Error
	}

	return transactions, nil
}
