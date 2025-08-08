package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/wallet", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "operation received"})
		})

		v1.GET("/wallets/:id", func(c *gin.Context) {
			id := c.Param("id")
			c.JSON(http.StatusOK, gin.H{"walletId": id, "balance": 0})
		})
	}

	return r
}
