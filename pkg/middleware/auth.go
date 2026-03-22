package middleware

import (
	"net/http"
	"strings"

	"gosecureskeleton/pkg/consts"
	"gosecureskeleton/pkg/errors"
	"gosecureskeleton/pkg/ext/db/sqlite"
	"gosecureskeleton/pkg/session"

	"github.com/gin-gonic/gin"
)

// AuthRequired 인증이 필요한 api 들을 위한 미들웨어로 header 토큰으로 인증 완료 후 사용자 정보를 context 에 저장한다.
func AuthRequired(sessions *session.Store, store *sqlite.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := TokenFromRequest(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": consts.MsgMissingToken})
			return
		}

		userID, ok := sessions.Lookup(token)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": errors.ErrInvalidToken.Error()})
			return
		}

		user, found, err := store.FindUserByID(c.Request.Context(), userID)
		if err != nil || !found {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": errors.ErrInvalidToken.Error()})
			return
		}

		c.Set(consts.ContextUserKey, user)
		c.Next()
	}
}

func TokenFromRequest(c *gin.Context) string {
	headerValue := strings.TrimSpace(c.GetHeader(consts.AuthorizationHeader))
	if headerValue != "" {
		return headerValue
	}

	cookieValue, err := c.Cookie(consts.AuthorizationCookieName)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(cookieValue)
}
