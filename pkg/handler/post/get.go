package post

import (
	"net/http"
	"strconv"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"

	"github.com/gin-gonic/gin"
)

// GET /api/posts/:id
func (h *PostHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidPostID})
		return
	}

	found, err := h.service.Get(c.Request.Context(), uint(id))
	if err != nil {
		handleError(c, "게시글 조회 실패", err)
		return
	}

	c.JSON(http.StatusOK, dtos.PostResponse{Post: found})
}
