package auth

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"

	"github.com/gin-gonic/gin"
)

const PathLogin = "/api/auth/login"

func validateLoginParam(req dtos.LoginRequest) string {
	if len(req.Username) < consts.MinUsernameLen || len(req.Username) > consts.MaxUsernameLen {
		return consts.MsgInvalidUsername
	}
	if len(req.Password) < consts.MinPasswordLen || len(req.Password) > consts.MaxPasswordLen {
		return consts.MsgInvalidPassword
	}
	return ""
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var request dtos.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidRequest})
		return
	}

	if msg := validateLoginParam(request); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: msg})
		return
	}

	user, token, err := h.service.Login(c.Request.Context(), request)
	if err != nil {
		handleError(c, "로그인 실패", err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(consts.AuthorizationCookieName, token, consts.CookieMaxAge, consts.CookiePath, "", false, true)
	c.JSON(http.StatusOK, dtos.LoginResponse{
		AuthMode: "header-and-cookie",
		Token:    token,
		User:     dtos.MakeUserResponse(user),
	})
}
