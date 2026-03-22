package auth

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/errors"
	"gosecureskeleton/pkg/service"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{service: svc}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST(PathRegister, h.Register)
	rg.POST(PathLogin, h.Login)
	rg.POST(PathLogout, h.Logout)
	rg.POST(PathWithdraw, h.Withdraw)
}

func ClearAuthorizationCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(consts.AuthorizationCookieName, "", -1, consts.CookiePath, "", false, true)
}

func handleError(c *gin.Context, msg string, err error) {
	if appErr, ok := errors.As(err); ok {
		c.JSON(appErr.Status, gin.H{consts.KeyMessage: appErr.Message})
		return
	}
	
	util.LogError(c.Request.Context(), msg, logrus.Fields{"error": err.Error()})
	c.JSON(http.StatusInternalServerError, gin.H{consts.KeyMessage: consts.MsgInternalServerError})
}
