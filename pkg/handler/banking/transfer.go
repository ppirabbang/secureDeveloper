package banking

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
)

const PathTransfer = "/api/banking/transfer"

func validateTransferParam(senderUsername, toUsername string, amount, balance int64) string {
	if amount <= 0 {
		return consts.MsgInvalidAmount
	}
	if senderUsername == toUsername {
		return consts.MsgCannotTransferSelf
	}
	if balance < amount {
		return consts.MsgInsufficientBalance
	}
	return ""
}

// POST /api/banking/transfer
func (h *BankingHandler) Transfer(c *gin.Context) {
	var request dtos.TransferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidRequest})
		return
	}

	user := util.GetCurrentUser(c)
	if msg := validateTransferParam(user.Username, request.ToUsername, request.Amount, user.Balance); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: msg})
		return
	}

	userResp, err := h.service.Transfer(c.Request.Context(), user, request.ToUsername, request.Amount)
	if err != nil {
		handleError(c, "이체 실패", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		consts.KeyMessage: consts.MsgTransferSuccess,
		consts.KeyUser:    userResp,
		consts.KeyTarget:  request.ToUsername,
		consts.KeyAmount:  request.Amount,
	})
}
