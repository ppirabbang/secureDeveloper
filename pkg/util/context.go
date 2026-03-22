package util

import (
	"github.com/gin-gonic/gin"
	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/dtos"
)

func GetCurrentUser(c *gin.Context) dtos.User {
	val, _ := c.Get(consts.ContextUserKey)
	user, _ := val.(dtos.User)
	return user
}
