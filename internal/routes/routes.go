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
		//auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", LoginUser)
			auth.POST("/register", CreateUser)
			//auth.GET("/refresh", RefreshToken)
		}
		// users routes

		users := v1.Group("/users")
		{
			users.GET("/me", GetUser)
			// users.PATCH("/me", UpdateUser)
		}

		/*
			passbooks := v1.Group("/passbooks")
			{
				passbooks.POST("", CreatePassbook)                // creates a new passbook
				passbooks.GET("", GetPassbooks)                   // gets all passbooks for a user
				passbooks.GET("/:passbook_id", GetPassbook)       // gets a passbook by id
				passbooks.PATCH("/:passbook_id", UpdatePassbook)  // updates a passbook by id
				passbooks.DELETE("/:passbook_id", DeletePassbook) // deletes a passbook by id

				transactions := passbooks.Group("/:passbook_id/transactions")
				{
					transactions.GET("", GetTransactions)                      // gets all transactions for a passbook
					transactions.POST("", CreateTransaction)                   // creates a new transaction for a passbook
					transactions.GET("/:transaction_id", GetTransaction)       // gets a transaction by id
					transactions.PATCH("/:transaction_id", UpdateTransaction)  // updates a transaction by id
					transactions.DELETE("/:transaction_id", DeleteTransaction) // deletes a transaction by id
				}
			}
		*/
	}

	return router
}
