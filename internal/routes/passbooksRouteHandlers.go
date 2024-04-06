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
		log.Println("Request sanitization and validation failed")
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
		UserID:        loggedInUserID,
		BankName:      passbook.BankName,
		AccountNumber: passbook.AccountNumber,
		TotalBalance:  passbook.TotalBalance,
		Nickname:      passbook.Nickname,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	_, err2 := initializers.DB.Exec(context.Background(), "INSERT INTO passbook_app.passbooks (passbook_id, user_id, bank_name, account_number, total_balance, nickname, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		pbook.PassbookID, pbook.UserID, pbook.BankName, pbook.AccountNumber, pbook.TotalBalance, pbook.Nickname, pbook.CreatedAt, pbook.UpdatedAt)
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
	// truncate total balance to 2 decimal places if more than 2 decimal digits
	(*pb).TotalBalance = float64(int((*pb).TotalBalance*100)) / 100

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

func GetPassbooks(ctx *gin.Context) {
	loggedInUserID := ctx.MustGet("userId").(string)
	log.Println("Getting Passbooks for user_id:", loggedInUserID)
	rows, err := initializers.DB.Query(context.Background(), "SELECT passbook_id, user_id, bank_name, account_number, total_balance, nickname, created_at, updated_at FROM passbook_app.passbooks WHERE user_id=$1", loggedInUserID)
	if err != nil {
		log.Println("Failed to get passbooks for user_id:", loggedInUserID)
		setErrorResponse(ctx, 500, "Failed to get passbooks")
		return
	}
	defer rows.Close()
	passbooks := make([]types.Passbook, 0)
	for rows.Next() {
		var p types.Passbook
		err := rows.Scan(&p.PassbookID, &p.UserID, &p.BankName, &p.AccountNumber, &p.TotalBalance, &p.Nickname, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			log.Println("Failed to get passbooks for user_id:", loggedInUserID)
			setErrorResponse(ctx, 500, "Failed to get passbooks")
			return
		}
		passbooks = append(passbooks, p)
	}
	log.Println("Passbooks fetched for user_id:", loggedInUserID)
	ctx.JSON(200, gin.H{
		"status": "success",
		"data": map[string][]types.Passbook{
			"passbooks": passbooks,
		},
	})
}

func GetPassbook(ctx *gin.Context) {
	loggedInUserID := ctx.MustGet("userId").(string)
	passbookID := ctx.Param("passbook_id")
	log.Println("Getting Passbook for user_id:", loggedInUserID, "passbook_id:", passbookID)
	row := initializers.DB.QueryRow(context.Background(), "SELECT passbook_id, user_id, bank_name, account_number, total_balance, nickname, created_at, updated_at FROM passbook_app.passbooks WHERE user_id=$1 AND passbook_id=$2", loggedInUserID, passbookID)
	var p types.Passbook
	err := row.Scan(&p.PassbookID, &p.UserID, &p.BankName, &p.AccountNumber, &p.TotalBalance, &p.Nickname, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Passbook not found for user_id:", loggedInUserID, "passbook_id:", passbookID)
			setErrorResponse(ctx, 404, "Passbook not found")
			return
		}
		log.Println("Failed to get passbook for user_id:", loggedInUserID, "passbook_id:", passbookID)
		setErrorResponse(ctx, 500, "Failed to get passbook")
		return
	}
	log.Println("Passbook fetched for user_id:", loggedInUserID, "passbook_id:", passbookID)
	ctx.JSON(200, gin.H{
		"status": "success",
		"data": map[string]types.Passbook{
			"passbook": p,
		},
	})
}
func UpdatePassbook(ctx *gin.Context) {
	loggedInUserID := ctx.MustGet("userId").(string)
	passbookID := ctx.Param("passbook_id")
	log.Println("Updating Passbook for user_id:", loggedInUserID, "passbook_id:", passbookID)
	var passbook types.Passbook
	if err := ctx.ShouldBindJSON(&passbook); err != nil {
		setErrorResponse(ctx, 400, "Invalid request body")
		return
	}
	// if user tries to update some other users passbook return 400
	if passbook.UserID != "" && passbook.UserID != loggedInUserID {
		log.Println("User_id in request body does not match with the logged in user_id")
		setErrorResponse(ctx, 403, "Passbook not owned by the logged in user")
		return
	} else if passbook.UserID == "" {
		log.Println("User_id not present in request body")
		setErrorResponse(ctx, 400, "User_id not present in request body")
		return
	}
	// input sanitization
	err := sanitizePassbookRequest(&passbook)
	if err != nil {
		log.Println("Request sanitization and validation failed")
		setErrorResponse(ctx, 400, err.Error())
		return
	}
	// check if the passbook exists for the user
	var existingId string
	err = initializers.DB.QueryRow(context.Background(), "SELECT passbook_id FROM passbook_app.passbooks WHERE user_id=$1 AND passbook_id=$2", loggedInUserID, passbookID).Scan(&existingId)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Passbook not found for user_id:", loggedInUserID, "passbook_id:", passbookID)
			setErrorResponse(ctx, 404, "Passbook not found")
			return
		}
		log.Println("Failed to check if passbook exists for user_id:", loggedInUserID)
		setErrorResponse(ctx, 500, "Failed to update passbook")
		return
	}
	passbook.UpdatedAt = time.Now().UTC()
	// passbook exists, update the passbook
	_, err2 := initializers.DB.Exec(context.Background(), "UPDATE passbook_app.passbooks SET bank_name=$1, account_number=$2, total_balance=$3, nickname=$4, updated_at=$5 WHERE user_id=$6 AND passbook_id=$7",
		passbook.BankName, passbook.AccountNumber, passbook.TotalBalance, passbook.Nickname, passbook.UpdatedAt, loggedInUserID, passbookID)
	if err2 != nil {
		log.Println(err2)
		log.Println("Failed to update passbook for user_id:", loggedInUserID, "passbook_id:", passbookID)
		setErrorResponse(ctx, 500, "Failed to update passbook")
		return
	}
	log.Println("Passbook updated for user_id:", loggedInUserID, "passbook_id:", passbookID)
	// return the passbook details along with generated passbook_id
	ctx.JSON(200, gin.H{
		"status": "success",
		"data": map[string]types.Passbook{
			"passbook": passbook,
		},
		"message": "Passbook updated successfully",
	})
}

func DeletePassbook(ctx *gin.Context) {
	loggedInUserID := ctx.MustGet("userId").(string)
	passbookID := ctx.Param("passbook_id")
	log.Println("Reqeust to delete passbook with id : ", passbookID)
	_, err := initializers.DB.Exec(context.Background(), "DELETE passbook_app.passbooks where user_id=$1 and passbook_id=$2", loggedInUserID, passbookID)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Println("Passbook not found for user_id:", loggedInUserID, "passbook_id:", passbookID)
			setErrorResponse(ctx, 404, "Passbook not found")
			return
		}
		log.Println("Failed to delete passbook for user_id:", loggedInUserID, "passbook_id:", passbookID)
		setErrorResponse(ctx, 500, "Failed to delete passbook")
		return
	}
	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": "Passbook deleted successfully",
	})
}
