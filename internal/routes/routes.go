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
	}

	return router
}
