package post

import (
	"net/http"
	"strconv"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/util"

	"github.com/gin-gonic/gin"
)

// DELETE /api/posts/:id
func (h *PostHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{consts.KeyMessage: consts.MsgInvalidPostID})
		return
	}

	user := util.GetCurrentUser(c)
	if err = h.service.Delete(c.Request.Context(), uint(id), user.ID); err != nil {
		handleError(c, "게시글 삭제 실패", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{consts.KeyMessage: consts.MsgPostDeleted})
}
