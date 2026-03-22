package auth

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/middleware"

	"github.com/gin-gonic/gin"
)

const PathWithdraw = "/api/auth/withdraw"

func validateWithdrawAccountParam(req dtos.WithdrawAccountRequest) string {
	if len(req.Password) < consts.MinPasswordLen || len(req.Password) > consts.MaxPasswordLen {
		return consts.MsgInvalidPassword
	}
	return ""
}

// POST /api/auth/withdraw
func (h *AuthHandler) Withdraw(c *gin.Context) {
	var request dtos.WithdrawAccountRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidRequest})
		return
	}

	if msg := validateWithdrawAccountParam(request); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: msg})
		return
	}

	token := middleware.TokenFromRequest(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{consts.KeyMessage: consts.MsgMissingToken})
		return
	}

	userResp, err := h.service.WithdrawAccount(c.Request.Context(), token, request.Password)
	if err != nil {
		handleError(c, "회원탈퇴 실패", err)
		return
	}

	ClearAuthorizationCookie(c)
	c.JSON(http.StatusOK, gin.H{
		consts.KeyMessage: consts.MsgWithdrawSuccess,
		consts.KeyUser:    userResp,
	})
}
