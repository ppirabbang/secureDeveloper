package auth

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
)

const PathRegister = "/api/auth/register"

func validateRegisterParam(req dtos.RegisterRequest) string {
	if len(req.Username) < consts.MinUsernameLen || len(req.Username) > consts.MaxUsernameLen {
		return consts.MsgInvalidUsername
	}
	if len(req.Name) < consts.MinNameLen || len(req.Name) > consts.MaxNameLen {
		return consts.MsgInvalidName
	}
	if len(req.Password) < consts.MinPasswordLen || len(req.Password) > consts.MaxPasswordLen {
		return consts.MsgInvalidPassword
	}
	if !util.IsValidEmail(req.Email) {
		return consts.MsgInvalidEmail
	}
	if !util.IsValidPhone(req.Phone) {
		return consts.MsgInvalidPhone
	}
	return ""
}

// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var request dtos.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidRequest})
		return
	}

	if msg := validateRegisterParam(request); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: msg})
		return
	}

	userResp, err := h.service.Register(c.Request.Context(), request)
	if err != nil {
		handleError(c, "회원가입 실패", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		consts.KeyMessage: consts.MsgRegisterSuccess,
		consts.KeyUser:    userResp,
	})
}
