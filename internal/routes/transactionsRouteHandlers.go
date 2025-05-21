package routes

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/akashsharma99/passbook-app/internal/initializers"
	"github.com/akashsharma99/passbook-app/internal/types"
	"github.com/akashsharma99/passbook-app/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5" // Added for pgx.ErrNoRows
	// pgxpool import removed if no longer directly needed after changing updatePassbookAndCreateTrx signature
	// database/sql import removed as sql.ErrNoRows is replaced by pgx.ErrNoRows
)

func CreateTransaction(ctx *gin.Context) {
	loggedInUserID := ctx.MustGet("user_id").(string)
	passbookID := ctx.Param("passbook_id")
	var transaction types.Transaction
	if err := ctx.ShouldBindJSON(&transaction); err != nil {
		setErrorResponse(ctx, 400, "Invalid request body")
		return
	}
	// input sanitization
	err := sanitizeTransactionRequest(&transaction)
	if err != nil {
		setErrorResponse(ctx, 400, err.Error())
		return
	}
	// return 403 if user_id in token does not match user_id in passbook or if passbook does not exist
	var pbid string
	err = initializers.DB.QueryRow(context.Background(), "SELECT passbook_id FROM passbook_app.passbooks WHERE passbook_id=$1 AND user_id=$2", passbookID, loggedInUserID).Scan(&pbid)
	if err != nil {
		setErrorResponse(ctx, 403, "Invalid passbook")
		return
	}
	// create a new transaction
	uid, uiderr := utils.GenerateUUID()
	if uiderr != nil {
		log.Println("Failed to generate transaction_id for user_id:", loggedInUserID)
		setErrorResponse(ctx, 500, "Failed to create transaction")
		return
	}
	timeNow := time.Now().UTC()
	tr := types.Transaction{
		TransactionID:   uid,
		Amount:          transaction.Amount,
		TransactionDate: transaction.TransactionDate,
		TransactionType: transaction.TransactionType,
		PartyName:       transaction.PartyName,
		Description:     transaction.Description,
		CreatedAt:       timeNow,
		UpdatedAt:       timeNow,
		Tags:            transaction.Tags,
		PassbookID:      passbookID,
		UserID:          loggedInUserID,
	}
	// update the passbook and create the transaction
	err = updatePassbookAndCreateTrx(initializers.DB, &tr)
	if err != nil {
		// if the new balance is less than 0, return an error
		if err.Error() == "insufficient balance" {
			setErrorResponse(ctx, 400, "Insufficient balance")
			return
		}
		setErrorResponse(ctx, 500, "Failed to create transaction")
		return
	}
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Transaction created successfully",
		"data": map[string]interface{}{
			"transaction": tr,
		},
	})

}

/*
Lock on the passbook before creating a transaction to update the total balance of the passbook depending on the transaction CREDIT or DEBIT.
Update the passbook's updated_at field and also disallow the transaction if the new balance is less than 0.
Create the transaction and commit the transaction.
*/
func updatePassbookAndCreateTrx(conn initializers.PgxPoolIface, tr *types.Transaction) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	// get the passbook details
	var passbook types.Passbook
	err = tx.QueryRow(context.Background(), "SELECT total_balance FROM passbook_app.passbooks WHERE passbook_id=$1 FOR UPDATE", tr.PassbookID).Scan(&passbook.TotalBalance)
	if err != nil {
		return err
	}
	// update the total balance of the passbook depending on the transaction type
	if tr.TransactionType == "CREDIT" {
		passbook.TotalBalance += tr.Amount
	} else {
		passbook.TotalBalance -= tr.Amount
	}
	// if the new balance is less than 0, return an error
	if passbook.TotalBalance < 0 {
		return errors.New("insufficient balance")
	}
	// update the passbook's updated_at and total_balance field
	passbook.UpdatedAt = tr.UpdatedAt
	_, err = tx.Exec(context.Background(), "UPDATE passbook_app.passbooks SET total_balance=$1, updated_at=$2 WHERE passbook_id=$3", passbook.TotalBalance, passbook.UpdatedAt, tr.PassbookID)
	if err != nil {
		return err
	}
	// create the transaction
	_, err = tx.Exec(context.Background(), "INSERT INTO passbook_app.transactions (transaction_id, amount, transaction_date, transaction_type, party_name, description, created_at, updated_at, tags, passbook_id, user_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", tr.TransactionID, tr.Amount, tr.TransactionDate, tr.TransactionType, tr.PartyName, tr.Description, tr.CreatedAt, tr.UpdatedAt, tr.Tags, tr.PassbookID, tr.UserID)
	if err != nil {
		return err
	}
	// commit the transaction
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func GetTransaction(ctx *gin.Context) {
	loggedInUserID := ctx.MustGet("user_id").(string)
	passbookID := ctx.Param("passbook_id")
	transactionID := ctx.Param("transaction_id")

	var transaction types.Transaction

	query := `
		SELECT 
			transaction_id, amount, transaction_date, transaction_type, 
			party_name, description, created_at, updated_at, tags, 
			passbook_id, user_id 
		FROM passbook_app.transactions 
		WHERE transaction_id=$1 AND passbook_id=$2 AND user_id=$3
	`

	err := initializers.DB.QueryRow(
		context.Background(),
		query,
		transactionID,
		passbookID,
		loggedInUserID,
	).Scan(
		&transaction.TransactionID,
		&transaction.Amount,
		&transaction.TransactionDate,
		&transaction.TransactionType,
		&transaction.PartyName,
		&transaction.Description,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.Tags,
		&transaction.PassbookID,
		&transaction.UserID,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // Changed to pgx.ErrNoRows
			setErrorResponse(ctx, 404, "Transaction not found or not authorized")
			return
		}
		log.Printf("Error fetching transaction %s for passbook %s: %v", transactionID, passbookID, err)
		setErrorResponse(ctx, 500, "Failed to retrieve transaction")
		return
	}

	ctx.JSON(200, gin.H{
		"status": "success",
		"data": map[string]interface{}{
			"transaction": transaction,
		},
	})
}

func sanitizeTransactionRequest(tr *types.Transaction) error {

	// transaction type should be part of slice ValidTransactionTypes
	(*tr).TransactionType = utils.TrimAndSanitizeStrict((*tr).TransactionType)
	if (*tr).TransactionType == "" || !utils.Contains(types.ValidTransactionTypes, (*tr).TransactionType) {
		return errors.New("invalid transaction type")
	}
	// amount should fit in DECIMAL(11,2)
	if (*tr).Amount <= 0 || (*tr).Amount > 999999999.99 {
		return errors.New("invalid amount")
	}
	// truncate total balance to 2 decimal places if more than 2 decimal digits
	(*tr).Amount = float64(int((*tr).Amount*100)) / 100

	// party name max 255 characters
	(*tr).PartyName = utils.TrimAndSanitizeStrict((*tr).PartyName)
	if (*tr).PartyName == "" || len((*tr).PartyName) > 255 {
		return errors.New("invalid party name")
	}
	// description
	(*tr).Description = utils.TrimAndSanitizeStrict((*tr).Description)
	// tags max 512 characters
	(*tr).Tags = utils.TrimAndSanitizeStrict((*tr).Tags)
	if len((*tr).Tags) > 512 {
		return errors.New("invalid tag length")
	}
	// transaction date should be a valid date and not empty
	if (*tr).TransactionDate.IsZero() {
		return errors.New("invalid transaction date")
	}

	return nil
}
