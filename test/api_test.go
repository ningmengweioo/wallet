package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"wallet/config"
	"wallet/router"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestAPI tests API functionality
func TestAPI(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Use path relative to project root
	os.Setenv("CONFIG_PATH", "../config/config.yaml")
	fmt.Printf("Set CONFIG_PATH to: %s\n", os.Getenv("CONFIG_PATH"))

	// Initialize config
	err := config.InitConfig()
	if err != nil {
		t.Fatalf("Failed to initialize config: %v", err)
	}

	// Get config
	cfg := config.GetConf()
	if cfg == nil {
		t.Fatalf("Failed to get config")
	}

	// Initialize database connection
	_, err = config.InitDB(cfg)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Create router
	r := router.SetupRouter()

	// Test user IDs (auto-increment, will be automatically assigned during testing)
	var userID1, userID2 int

	// Test 1: User Registration
	t.Run("RegisterUser", func(t *testing.T) {
		// 准备请求体（移除ID字段，因为现在是自增长的）
		registerReq := map[string]string{
			"username": "Test User",
			"email":    "test@example.com",
		}
		body, _ := json.Marshal(registerReq)

		// Create request
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Record response
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		// Verify response
		assert.Equal(t, http.StatusCreated, w.Code)

		// 验证响应体
		var response struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, response.Code)

		// 提取创建的第二个用户ID
		if userData, ok := response.Data["user"].(map[string]interface{}); ok {
			if idFloat, ok := userData["id"].(float64); ok {
				userID2 = int(idFloat)
				fmt.Printf("Created second user with ID: %d\n", userID2)
			}
		} // Extract created second user ID
		assert.Contains(t, response.Message, "created successfully")

		// Extract created user ID
		if userData, ok := response.Data["user"].(map[string]interface{}); ok {
			if idFloat, ok := userData["id"].(float64); ok {
				userID1 = int(idFloat)
				fmt.Printf("Created user with ID: %d\n", userID1)
			}
		}
	})

	// Test 2: Duplicate User Registration (should fail)
	t.Run("RegisterDuplicateUser", func(t *testing.T) {
		registerReq := map[string]string{
			"username": "Test User Duplicate",
			"email":    "test@example.com", // 使用相同邮箱
		}
		body, _ := json.Marshal(registerReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test 3: Get User Information
	t.Run("GetUser", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/users/%d", userID1), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	// Test 4: Deposit Operation
	t.Run("Deposit", func(t *testing.T) {
		depositReq := map[string]interface{}{
			"amount":      100.00,
			"description": "Initial deposit",
		}
		body, _ := json.Marshal(depositReq)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/wallets/%d/deposit", userID1), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		balance, ok := response.Data["balance"].(float64)
		assert.True(t, ok)
		assert.GreaterOrEqual(t, balance, 100.00)
	})

	// Test 5: Get Balance
	t.Run("GetBalance", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/wallets/%d/balance", userID1), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
		_, ok := response.Data["balance"]
		assert.True(t, ok)
	})

	// Test 6: Withdrawal Operation
	t.Run("Withdraw", func(t *testing.T) {
		withdrawReq := map[string]interface{}{
			"amount":      30.00,
			"description": "Withdrawal test",
		}
		body, _ := json.Marshal(withdrawReq)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/wallets/%d/withdraw", userID1), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	// Test 7: Withdrawal with Insufficient Balance (should fail)
	t.Run("WithdrawInsufficientBalance", func(t *testing.T) {
		withdrawReq := map[string]interface{}{
			"amount":      1000.00, // Exceeding current balance
			"description": "Large withdrawal",
		}
		body, _ := json.Marshal(withdrawReq)

		req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/wallets/%d/withdraw", userID1), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test 8: Register Second User for Transfer Test
	t.Run("RegisterSecondUser", func(t *testing.T) {
		registerReq := map[string]string{
			"username": "Second User",
			"email":    "second@example.com",
		}
		body, _ := json.Marshal(registerReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	// Test 9: Transfer Operation
	t.Run("Transfer", func(t *testing.T) {
		transferReq := map[string]interface{}{
			"from_user_id": fmt.Sprintf("%d", userID1),
			"to_user_id":   fmt.Sprintf("%d", userID2),
			"amount":       20.00,
			"description":  "Test transfer",
		}
		body, _ := json.Marshal(transferReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets/transfer", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Code    int                    `json:"code"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	// Test 10: Self Transfer (should fail)
	t.Run("TransferToSelf", func(t *testing.T) {
		transferReq := map[string]interface{}{
			"from_user_id": fmt.Sprintf("%d", userID1),
			"to_user_id":   fmt.Sprintf("%d", userID1),
			"amount":       10.00,
			"description":  "Self transfer",
		}
		body, _ := json.Marshal(transferReq)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets/transfer", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Test 11: Get User Transactions
	t.Run("GetUserTransactions", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/transactions/%d?page=1&limit=10", userID1), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Data    interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	// Test 12: Get All Users
	t.Run("GetAllUsers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Data    interface{} `json:"data"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	// Test 13: Health Check
	t.Run("HealthCheck", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ok", response["status"])
	})
}
