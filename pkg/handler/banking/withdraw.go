package banking

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
)

const PathWithdraw = "/api/banking/withdraw"

func validateWithdrawParam(amount, balance int64) string {
	if amount <= 0 {
		return consts.MsgInvalidAmount
	}
	if balance < amount {
		return consts.MsgInsufficientBalance
	}
	return ""
}

// POST /api/banking/withdraw
func (h *BankingHandler) Withdraw(c *gin.Context) {
	var request dtos.BalanceWithdrawRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidRequest})
		return
	}

	user := util.GetCurrentUser(c)
	if msg := validateWithdrawParam(request.Amount, user.Balance); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: msg})
		return
	}

	userResp, err := h.service.Withdraw(c.Request.Context(), user, request.Amount)
	if err != nil {
		handleError(c, "출금 실패", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		consts.KeyMessage: consts.MsgBalanceWithdraw,
		consts.KeyUser:    userResp,
		consts.KeyAmount:  request.Amount,
	})
}
