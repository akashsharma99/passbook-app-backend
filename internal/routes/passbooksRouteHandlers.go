package routes

import (
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/akashsharma99/passbook-app/internal/initializers"
	"github.com/akashsharma99/passbook-app/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/microcosm-cc/bluemonday"
)

func CreatePassbook(ctx *gin.Context) {
	loggedInUserID := ctx.MustGet("userId").(string)
	log.Println("Creating Passbook for user_id:", loggedInUserID)
	var passbook types.Passbook
	if err := ctx.ShouldBindJSON(&passbook); err != nil {
		setErrorResponse(ctx, 400, "Invalid request body")
		return
	}
	// input sanitization
	err := sanitizePassbookRequest(&passbook)
	if err != nil {
		setErrorResponse(ctx, 400, err.Error())
		return
	}
	// check if the passbook already exists for the user
	var existingId string
	err = initializers.DB.QueryRow(context.Background(), "SELECT passbook_id FROM passbook_app.passbooks WHERE user_id=$1 AND bank_name=$2 AND account_number=$3", loggedInUserID, passbook.BankName, passbook.AccountNumber).Scan(&existingId)
	if err != nil && err != pgx.ErrNoRows {
		log.Println("Failed to check if passbook exists for user_id:", loggedInUserID)
		setErrorResponse(ctx, 500, "Failed to create passbook")
		return
	}
	if existingId != "" {
		log.Println("Passbook already exists for user_id:", loggedInUserID)
		setErrorResponse(ctx, 400, "Account already exists")
		return
	}
	// passbook does not exist, create a new passbook
	uid, uiderr := uuid.NewRandom()
	if uiderr != nil {
		log.Println("Failed to generate passbook_id for user_id:", loggedInUserID)
		setErrorResponse(ctx, 500, "Failed to create passbook")
		return
	}
	pbook := types.Passbook{
		PassbookID:    uid.String(),
		BankName:      passbook.BankName,
		AccountNumber: passbook.AccountNumber,
		TotalBalance:  passbook.TotalBalance,
		Nickname:      passbook.Nickname,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	_, err2 := initializers.DB.Exec(context.Background(), "INSERT INTO passbook_app.passbooks (passbook_id, user_id, bank_name, account_number, total_balance, nickname, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		pbook.PassbookID, loggedInUserID, pbook.BankName, pbook.AccountNumber, pbook.TotalBalance, pbook.Nickname, pbook.CreatedAt, pbook.UpdatedAt)
	if err2 != nil {
		log.Println(err2)
		log.Println("Failed to create passbook for user_id:", loggedInUserID)
		setErrorResponse(ctx, 500, "Failed to create passbook")
		return
	}
	log.Println("Passbook created for user_id:", loggedInUserID)
	// return the passbook details along with generated passbook_id
	ctx.JSON(201, gin.H{
		"status": "success",
		"data": map[string]types.Passbook{
			"passbook": pbook,
		},
	})
}

func sanitizePassbookRequest(pb *types.Passbook) error {

	// bankname validations
	(*pb).BankName = trimAndSanitizeString((*pb).BankName)
	if (*pb).BankName == "" || len((*pb).BankName) > 255 {
		return errors.New("invalid bank name")
	}
	// account number validations
	(*pb).AccountNumber = trimAndSanitizeString((*pb).AccountNumber)
	if (*pb).AccountNumber == "" || len((*pb).AccountNumber) > 255 {
		return errors.New("invalid account number")
	}
	// total balance should fit in DECIMAL(10,2)
	if (*pb).TotalBalance < 0 || (*pb).TotalBalance > 999999999.99 {
		return errors.New("invalid total balance")
	}
	// nickname validations
	(*pb).Nickname = trimAndSanitizeString((*pb).Nickname)
	if (*pb).Nickname == "" || len((*pb).Nickname) > 255 {
		return errors.New("invalid nickname")
	}
	return nil
}
func trimAndSanitizeString(s string) string {
	p := bluemonday.StrictPolicy()
	s = strings.TrimSpace(s)
	s = p.Sanitize(s)
	return s
}
