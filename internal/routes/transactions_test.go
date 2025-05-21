package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/akashsharma99/passbook-app/internal/initializers"
	"github.com/akashsharma99/passbook-app/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

// normalizeSQL replaces all whitespace with a single space and trims for pgxmock
func normalizeSQL(sql string) string {
	s := strings.TrimSpace(sql)
	return regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
}

func TestGetTransaction(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Initialize pgxmock pool with QueryMatcherRegexp option
	mockDB, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("Failed to create mock pool: %v", err)
	}
	defer mockDB.Close()

	originalDB := initializers.DB
	initializers.DB = mockDB // This assignment should now work
	defer func() { initializers.DB = originalDB }()

	// Expected SQL query from GetTransaction handler (normalized)
	// Using pgxmock.QueryMatcherRegexp for more robust matching.
	expectedSQL := `^SELECT transaction_id, amount, transaction_date, transaction_type, party_name, description, created_at, updated_at, tags, passbook_id, user_id FROM passbook_app.transactions WHERE transaction_id=\$1 AND passbook_id=\$2 AND user_id=\$3$`

	testUserID := "test-user-id"
	testPassbookID := "test-passbook-id"
	testTransactionID := "test-transaction-id"

	// Middleware to inject user_id into context for testing
	authTestMiddleware := func() gin.HandlerFunc {
		return func(c *gin.Context) {
			c.Set("user_id", testUserID)
			c.Next()
		}
	}
	
	router := gin.New() // Use gin.New() instead of NewRouter() to isolate with test middleware
	v1 := router.Group("/v1")
	passbooks := v1.Group("/passbooks", authTestMiddleware()) // Apply middleware here
	{
		transactions := passbooks.Group("/:passbook_id/transactions")
		{
			transactions.GET("/:transaction_id", GetTransaction)
		}
	}
	// Use pgxmock.QueryMatcherRegexp for all mock expectations
	mockDB.MatchExpectationsInOrder(false) // Allow out-of-order expectation matching if necessary for parallel tests later, though not strictly needed here.


	t.Run("Successful fetch", func(t *testing.T) {
		sampleTime := time.Now().UTC().Truncate(time.Microsecond) 
		expectedTransaction := types.Transaction{
			TransactionID:   testTransactionID,
			Amount:          100.50,
			TransactionDate: sampleTime,
			TransactionType: "CREDIT",
			PartyName:       "Test Party",
			Description:     "Test Description",
			CreatedAt:       sampleTime,
			UpdatedAt:       sampleTime,
			Tags:            "test,fetch",
			PassbookID:      testPassbookID,
			UserID:          testUserID,
		}

		rows := pgxmock.NewRows([]string{
			"transaction_id", "amount", "transaction_date", "transaction_type",
			"party_name", "description", "created_at", "updated_at", "tags",
			"passbook_id", "user_id",
		}).AddRow(
			expectedTransaction.TransactionID,
			expectedTransaction.Amount,
			expectedTransaction.TransactionDate,
			expectedTransaction.TransactionType,
			expectedTransaction.PartyName,
			expectedTransaction.Description,
			expectedTransaction.CreatedAt,
			expectedTransaction.UpdatedAt,
			expectedTransaction.Tags,
			expectedTransaction.PassbookID,
			expectedTransaction.UserID,
		)

		mockDB.ExpectQuery(expectedSQL).
			WithArgs(testTransactionID, testPassbookID, testUserID).
			WillReturnRows(rows)

		w := httptest.NewRecorder()
		reqURL := fmt.Sprintf("/v1/passbooks/%s/transactions/%s", testPassbookID, testTransactionID)
		req, _ := http.NewRequest("GET", reqURL, nil)
		
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		assert.Equal(t, "success", responseBody["status"])
		
		data, ok := responseBody["data"].(map[string]interface{})
		assert.True(t, ok, "data field is not a map")

		transactionData, ok := data["transaction"].(map[string]interface{})
		assert.True(t, ok, "transaction field is not a map")
		
		assert.Equal(t, expectedTransaction.TransactionID, transactionData["transaction_id"])
		assert.Equal(t, expectedTransaction.Amount, transactionData["amount"])
		assert.Equal(t, expectedTransaction.TransactionType, transactionData["transaction_type"])
		assert.Equal(t, expectedTransaction.PartyName, transactionData["party_name"])
		// Time needs special handling for comparison due to potential time zone/format issues in JSON
		parsedTransactionDate, _ := time.Parse(time.RFC3339Nano, transactionData["transaction_date"].(string))
		assert.True(t, expectedTransaction.TransactionDate.Equal(parsedTransactionDate))


		assert.NoError(t, mockDB.ExpectationsWereMet(), "pgxmock expectations not met")
	})

	t.Run("Transaction not found", func(t *testing.T) {
		mockDB.ExpectQuery(expectedSQL).
			WithArgs(testTransactionID, testPassbookID, testUserID).
			WillReturnError(pgx.ErrNoRows)

		w := httptest.NewRecorder()
		reqURL := fmt.Sprintf("/v1/passbooks/%s/transactions/%s", testPassbookID, testTransactionID)
		req, _ := http.NewRequest("GET", reqURL, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "Transaction not found or not authorized", responseBody["message"])
		
		assert.NoError(t, mockDB.ExpectationsWereMet(), "pgxmock expectations not met")
	})
	
	t.Run("Database error", func(t *testing.T) {
		mockDB.ExpectQuery(expectedSQL).
			WithArgs(testTransactionID, testPassbookID, testUserID).
			WillReturnError(fmt.Errorf("some db error")) // Generic DB error

		w := httptest.NewRecorder()
		reqURL := fmt.Sprintf("/v1/passbooks/%s/transactions/%s", testPassbookID, testTransactionID)
		req, _ := http.NewRequest("GET", reqURL, nil)

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var responseBody map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		assert.Equal(t, "error", responseBody["status"])
		assert.Equal(t, "Failed to retrieve transaction", responseBody["message"])

		assert.NoError(t, mockDB.ExpectationsWereMet(), "pgxmock expectations not met")
	})
}

// This function is defined in other route handler files,
// but not exported, so we define a local version for testing.
func testSetErrorResponse(ctx *gin.Context, code int, message string) { // Renamed
	ctx.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}
