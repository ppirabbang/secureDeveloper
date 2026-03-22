package post

import (
	"net/http"
	"strconv"
	"strings"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
)

// PUT /api/posts/:id
func (h *PostHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidPostID})
		return
	}

	var request dtos.UpdatePostRequest
	if err = c.ShouldBindJSON(&request); err != nil {
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
	updated, err := h.service.Update(c.Request.Context(), uint(id), user.ID, title, content)
	if err != nil {
		handleError(c, "게시글 수정 실패", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		consts.KeyMessage: consts.MsgPostUpdated,
		consts.KeyPost:    updated,
	})
}
