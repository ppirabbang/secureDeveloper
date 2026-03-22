package auth

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/middleware"

	"github.com/gin-gonic/gin"
)

const PathLogout = "/api/auth/logout"

// POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	token := middleware.TokenFromRequest(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{consts.KeyMessage: consts.MsgMissingToken})
		return
	}

	if err := h.service.Logout(c.Request.Context(), token); err != nil {
		handleError(c, "로그아웃 실패", err)
		return
	}

	ClearAuthorizationCookie(c)
	c.JSON(http.StatusOK, gin.H{consts.KeyMessage: consts.MsgLogoutSuccess})
}
