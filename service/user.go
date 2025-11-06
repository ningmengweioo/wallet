package service

import (
	"errors"

	"wallet/config"
	"wallet/models"
)

// UserServiceImpl implements user service interfaces
type UserServiceImpl struct{}

// NewUserService creates user service instance
func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{}
}

// RegisterUser registers a new user with a wallet
func (s *UserServiceImpl) RegisterUser(username, email string) (*models.Users, *models.Wallets, error) {
	tx := config.GetDB().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := &models.Users{
		Username: username,
		Email:    email,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// Create wallet
	wallet := &models.Wallets{
		UserID:  user.ID, // Use auto-increment ID
		Balance: 0,
	}

	if err := tx.Create(wallet).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	return user, wallet, nil
}

// GetUserByID retrieves user information by ID
func (s *UserServiceImpl) GetUserByID(id int) (*models.Users, error) {
	var user models.Users
	if result := config.GetDB().Preload("Wallet").Where("id = ?", id).First(&user); result.Error != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// GetAllUsers retrieves all users information
func (s *UserServiceImpl) GetAllUsers() ([]models.Users, error) {
	var users []models.Users
	if result := config.GetDB().Preload("Wallet").Find(&users); result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// UserExistsByEmail checks if user exists by email
func (s *UserServiceImpl) UserExistsByEmail(email string) bool {
	var existingUser models.Users
	return config.GetDB().Where("email = ?", email).First(&existingUser).Error == nil
}
