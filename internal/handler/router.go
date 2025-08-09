package handler

import (
	"log"
	"net/http"

	"github.com/L0Qqi/wallet-api/internal/model"
	"github.com/L0Qqi/wallet-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupRouter(walletService *service.WalletService) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/wallet", func(c *gin.Context) {
			var req model.WalletOperationRequest

			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			if req.Amount <= 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be > 0"})
				return
			}

			err := walletService.Operate(c.Request.Context(), req)
			if err != nil {
				log.Printf("Operate error: %v", err)

				if err.Error() == "not enouth money" {
					c.JSON(http.StatusBadRequest, gin.H{"error": "cannot do this operation, not enough money"})
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "operation received"})
		})

		v1.GET("/wallets/:id", func(c *gin.Context) {
			id, err := uuid.Parse(c.Param("id"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
				return
			}

			balance, err := walletService.GetBalance(c.Request.Context(), id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{"walletId": id, "balance": balance})
		})
	}

	return r
}
