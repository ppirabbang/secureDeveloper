package post

import (
	"net/http"

	"gosecureskeleton/pkg/dtos"

	"github.com/gin-gonic/gin"
)

// GET /api/posts
func (h *PostHandler) List(c *gin.Context) {
	posts, err := h.service.List(c.Request.Context())
	if err != nil {
		handleError(c, "게시글 목록 조회 실패", err)
		return
	}

	c.JSON(http.StatusOK, dtos.PostListResponse{Posts: posts})
}
