package routes

import "github.com/gin-gonic/gin"

// create a router using gin and return it
func NewRouter() *gin.Engine {
	router := gin.Default()
	// add routes for v1 of api
	v1 := router.Group("/v1")
	{
		//health check route which return 200 OK
		v1.GET(("/health"), func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "OK",
			})
		})
		// users routes
		v1.POST("/users", CreateUser)
		v1.POST("/users/login", LoginUser)
		v1.GET("/users/me", GetUser)
		v1.PATCH("/users/me", ResetPassword)
		// passbooks routes
		v1.POST("/passbooks", CreatePassbook)                // creates a new passbook
		v1.GET("/passbooks", GetPassbooks)                   // gets all passbooks for a user
		v1.GET("/passbooks/:passbook_id", GetPassbook)       // gets a passbook by id
		v1.PATCH("/passbooks/:passbook_id", UpdatePassbook)  // updates a passbook by id
		v1.DELETE("/passbooks/:passbook_id", DeletePassbook) // deletes a passbook by id
		// transactions routes
		v1.GET("/passbooks/:passbook_id/transactions", GetTransactions)                      // get all transactions for a passbook
		v1.POST("/passbooks/:passbook_id/transactions", CreateTransaction)                   // create a new transaction for a passbook
		v1.GET("/passbooks/:passbook_id/transactions/:transaction_id", GetTransaction)       // get a transaction by id
		v1.PATCH("/passbooks/:passbook_id/transactions/:transaction_id", UpdateTransaction)  // update a transaction by id
		v1.DELETE("/passbooks/:passbook_id/transactions/:transaction_id", DeleteTransaction) // delete a transaction by id
	}

	return router
}
