package banking

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/errors"
	"gosecureskeleton/pkg/service"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BankingHandler struct {
	service *service.BankingService
}

func NewBankingHandler(svc *service.BankingService) *BankingHandler {
	return &BankingHandler{service: svc}
}

func (h *BankingHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST(PathDeposit, h.Deposit)
	rg.POST(PathWithdraw, h.Withdraw)
	rg.POST(PathTransfer, h.Transfer)
}

func handleError(c *gin.Context, msg string, err error) {
	if appErr, ok := errors.As(err); ok {
		c.JSON(appErr.Status, gin.H{consts.KeyMessage: appErr.Message})
		return
	}
	util.LogError(c.Request.Context(), msg, logrus.Fields{"error": err.Error()})
	c.JSON(http.StatusInternalServerError, gin.H{consts.KeyMessage: consts.MsgInternalServerError})
}
