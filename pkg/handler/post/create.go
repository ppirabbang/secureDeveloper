package post

import (
	"net/http"
	"strings"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
)

func validatePostParam(title, content string) string {
	if len(title) == 0 || len(title) > consts.MaxTitleLen {
		return consts.MsgInvalidTitle
	}
	if len(content) == 0 || len(content) > consts.MaxContentLen {
		return consts.MsgInvalidContent
	}
	return ""
}

// POST /api/posts
func (h *PostHandler) Create(c *gin.Context) {
	var request dtos.CreatePostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidRequest})
		return
	}

	title := strings.TrimSpace(request.Title)
	content := strings.TrimSpace(request.Content)
	if msg := validatePostParam(title, content); msg != "" {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: msg})
		return
	}

	user := util.GetCurrentUser(c)
	created, err := h.service.Create(c.Request.Context(), user.ID, title, content)
	if err != nil {
		handleError(c, "게시글 생성 실패", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		consts.KeyMessage: consts.MsgPostCreated,
		consts.KeyPost:    created,
	})
}
