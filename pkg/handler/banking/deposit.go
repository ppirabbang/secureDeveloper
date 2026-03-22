package banking

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
)

const PathDeposit = "/api/banking/deposit"

func validateDepositParam(amount int64) string {
	if amount <= 0 {
		return consts.MsgInvalidAmount
	}
	return ""
}

// POST /api/banking/deposit
func (h *BankingHandler) Deposit(c *gin.Context) {
	var request dtos.DepositRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidRequest})
		return
	}

	if msg := validateDepositParam(request.Amount); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: msg})
		return
	}

	user := util.GetCurrentUser(c)
	userResp, err := h.service.Deposit(c.Request.Context(), user, request.Amount)
	if err != nil {
		handleError(c, "입금 실패", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		consts.KeyMessage: consts.MsgDepositSuccess,
		consts.KeyUser:    userResp,
		consts.KeyAmount:  request.Amount,
	})
}
