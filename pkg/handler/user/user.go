package user

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
)

const PathMe = "/api/me"

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET(PathMe, h.Me)
}

// GET /api/me
func (h *UserHandler) Me(c *gin.Context) {
	user := util.GetCurrentUser(c)
	c.JSON(http.StatusOK, gin.H{consts.KeyUser: dtos.MakeUserResponse(user)})
}
