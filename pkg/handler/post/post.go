package post

import (
	"net/http"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/errors"
	"gosecureskeleton/pkg/service"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	PathPosts    = "/api/posts"
	PathPostByID = "/api/posts/:id"
)

type PostHandler struct {
	service *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{service: svc}
}

func (h *PostHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET(PathPosts, h.List)
	rg.POST(PathPosts, h.Create)
	rg.GET(PathPostByID, h.Get)
	rg.PUT(PathPostByID, h.Update)
	rg.DELETE(PathPostByID, h.Delete)
}

func handleError(c *gin.Context, msg string, err error) {
	if appErr, ok := errors.As(err); ok {
		c.JSON(appErr.Status, gin.H{consts.KeyMessage: appErr.Message})
		return
	}
	util.LogError(c.Request.Context(), msg, logrus.Fields{"error": err.Error()})
	c.JSON(http.StatusInternalServerError, gin.H{consts.KeyMessage: consts.MsgInternalServerError})
}
