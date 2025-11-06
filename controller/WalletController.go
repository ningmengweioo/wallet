package controller

import (
	"strconv"

	"wallet/service"
	"wallet/utils"

	"github.com/gin-gonic/gin"
)

// RegisterUser registers a new user
func RegisterUser(c *gin.Context) {
	type RegisterRequest struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters: "+err.Error())
		return
	}

	// Use service layer for user registration
	userService := NewUserService()
	if userService.UserExistsByEmail(req.Email) {
		utils.BadRequest(c, "User already exists")
		return
	}

	user, wallet, err := userService.RegisterUser(req.Username, req.Email)
	if err != nil {
		utils.InternalError(c, "Failed to create user")
		return
	}

	utils.Created(c, gin.H{"user": user, "wallet": wallet})
}

// GetBalance retrieves user's wallet balance
func GetBalance(c *gin.Context) {
	// Convert path parameter from string to int
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}

	// Use service layer to get balance
	walletService := NewWalletService()
	balance, err := walletService.GetBalance(userID)
	if err != nil {
		utils.NotFound(c, "Wallet not found")
		return
	}

	utils.Success(c, gin.H{"balance": balance})
}

// Deposit adds funds to user's wallet
func Deposit(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}

	type DepositRequest struct {
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	var req DepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	// Use service layer for deposit operation
	walletService := NewWalletService()
	balance, err := walletService.Deposit(userID, req.Amount, req.Description)
	if err != nil {
		if err.Error() == "wallet not found" {
			utils.NotFound(c, "Wallet not found")
		} else {
			utils.InternalError(c, "Failed to deposit")
		}
		return
	}

	utils.Success(c, gin.H{"balance": balance})
}

// Withdraw removes funds from user's wallet
func Withdraw(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}

	type WithdrawRequest struct {
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	var req WithdrawRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	// Use service layer for withdrawal operation
	walletService := NewWalletService()
	balance, err := walletService.Withdraw(userID, req.Amount, req.Description)
	if err != nil {
		switch err.Error() {
		case "wallet not found":
			utils.NotFound(c, "Wallet not found")
		case "insufficient balance":
			utils.BadRequest(c, "Insufficient balance")
		default:
			utils.InternalError(c, "Failed to withdraw")
		}
		return
	}

	utils.Success(c, gin.H{"balance": balance})
}

// Transfer moves funds between user wallets
func Transfer(c *gin.Context) {
	type TransferRequest struct {
		FromUserID  string  `json:"from_user_id" binding:"required"`
		ToUserID    string  `json:"to_user_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required,gt=0"`
		Description string  `json:"description"`
	}

	var req TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request parameters")
		return
	}

	if req.FromUserID == req.ToUserID {
		utils.BadRequest(c, "Cannot transfer to self")
		return
	}

	// Convert user IDs from string to int
	fromUserID, err := strconv.Atoi(req.FromUserID)
	if err != nil {
		utils.BadRequest(c, "Invalid sender user ID format")
		return
	}

	toUserID, err := strconv.Atoi(req.ToUserID)
	if err != nil {
		utils.BadRequest(c, "Invalid recipient user ID format")
		return
	}

	// Use service layer for transfer operation
	walletService := NewWalletService()
	fromBalance, toBalance, err := walletService.Transfer(fromUserID, toUserID, req.Amount, req.Description)
	if err != nil {
		switch err.Error() {
		case "sender wallet not found":
			utils.NotFound(c, "Sender wallet not found")
		case "recipient wallet not found":
			utils.NotFound(c, "Recipient wallet not found")
		case "insufficient balance":
			utils.BadRequest(c, "Insufficient balance")
		default:
			utils.InternalError(c, "Failed to transfer")
		}
		return
	}

	utils.Success(c, gin.H{
		"from_balance": fromBalance,
		"to_balance":   toBalance,
	})
}

// GetUser retrieves user information
func GetUser(c *gin.Context) {
	userID := c.Param("id")

	// Convert ID to integer for query
	idInt, err := strconv.Atoi(userID)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}

	// Use service layer to get user information
	userService := NewUserService()
	user, err := userService.GetUserByID(idInt)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.Success(c, user)
}

// GetAllUsers retrieves all users information
func GetAllUsers(c *gin.Context) {
	// Use service layer to get all users
	userService := NewUserService()
	users, err := userService.GetAllUsers()
	if err != nil {
		utils.InternalError(c, "Failed to fetch users")
		return
	}

	utils.Success(c, users)
}

// GetUserTransactions retrieves user's transaction history
func GetUserTransactions(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID format")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Use service layer to get transaction records
	transactionService := NewTransactionService()
	transactions, err := transactionService.GetUserTransactions(userID, page, limit)
	if err != nil {
		utils.InternalError(c, "Failed to fetch transactions")
		return
	}

	utils.Success(c, transactions)
}

// NewUserService creates user service instance
func NewUserService() *service.UserServiceImpl {
	return service.NewUserService()
}

// NewWalletService creates wallet service instance
func NewWalletService() *service.WalletServiceImpl {
	return service.NewWalletService()
}

// NewTransactionService creates transaction service instance
func NewTransactionService() *service.TransactionServiceImpl {
	return service.NewTransactionService()
}
