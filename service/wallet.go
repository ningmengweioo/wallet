package service

import (
	"errors"

	"wallet/config"
	"wallet/models"
)

// WalletServiceImpl implements wallet service interfaces
type WalletServiceImpl struct{}

// NewWalletService creates wallet service instance
func NewWalletService() *WalletServiceImpl {
	return &WalletServiceImpl{}
}

// GetBalance retrieves wallet balance
func (s *WalletServiceImpl) GetBalance(userID int) (float64, error) {
	var wallet models.Wallets
	if result := config.GetDB().Where("user_id = ?", userID).First(&wallet); result.Error != nil {
		return 0, errors.New("wallet not found")
	}

	return wallet.Balance, nil
}

// Deposit adds funds to wallet
func (s *WalletServiceImpl) Deposit(userID int, amount float64, description string) (float64, error) {
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var wallet models.Wallets
	if result := tx.Where("user_id = ?", userID).First(&wallet); result.Error != nil {
		tx.Rollback()
		return 0, errors.New("wallet not found")
	}

	wallet.Balance += amount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	transaction := models.Transaction{
		Type:        "deposit",
		ToUserID:    userID,
		Amount:      amount,
		Description: description,
		Status:      "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}

// Withdraw removes funds from wallet
func (s *WalletServiceImpl) Withdraw(userID int, amount float64, description string) (float64, error) {
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var wallet models.Wallets
	if result := tx.Where("user_id = ?", userID).First(&wallet); result.Error != nil {
		tx.Rollback()
		return 0, errors.New("wallet not found")
	}

	if wallet.Balance < amount {
		tx.Rollback()
		return 0, errors.New("insufficient balance")
	}

	wallet.Balance -= amount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	transaction := models.Transaction{
		Type:        "withdraw",
		FromUserID:  userID,
		Amount:      amount,
		Description: description,
		Status:      "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return wallet.Balance, nil
}

// Transfer moves funds between wallets
func (s *WalletServiceImpl) Transfer(fromUserID, toUserID int, amount float64, description string) (float64, float64, error) {
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var fromWallet, toWallet models.Wallets
	if result := tx.Where("user_id = ?", fromUserID).First(&fromWallet); result.Error != nil {
		tx.Rollback()
		return 0, 0, errors.New("sender wallet not found")
	}

	if result := tx.Where("user_id = ?", toUserID).First(&toWallet); result.Error != nil {
		tx.Rollback()
		return 0, 0, errors.New("recipient wallet not found")
	}

	if fromWallet.Balance < amount {
		tx.Rollback()
		return 0, 0, errors.New("insufficient balance")
	}

	fromWallet.Balance -= amount
	toWallet.Balance += amount

	if err := tx.Save(&fromWallet).Error; err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	if err := tx.Save(&toWallet).Error; err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	transaction := models.Transaction{
		Type:        "transfer",
		FromUserID:  fromUserID,
		ToUserID:    toUserID,
		Amount:      amount,
		Description: description,
		Status:      "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return 0, 0, err
	}

	return fromWallet.Balance, toWallet.Balance, nil
}
