package handler

import (
	"strings"

	"gosecureskeleton/pkg/ext/db/sqlite"
	"gosecureskeleton/pkg/handler/auth"
	"gosecureskeleton/pkg/handler/banking"
	"gosecureskeleton/pkg/handler/post"
	"gosecureskeleton/pkg/handler/user"
	"gosecureskeleton/pkg/middleware"
	"gosecureskeleton/pkg/service"
	"gosecureskeleton/pkg/session"

	"github.com/gin-gonic/gin"
)

func SetupRouter(store *sqlite.Store, sessions *session.Store) *gin.Engine {
	router := gin.New()
	router.Use(middleware.SetTraceID())
	router.Use(middleware.RequestLogger())
	router.Use(gin.Recovery())

	registerStaticRoutes(router)
	registerAPIRoutes(router, store, sessions)

	return router
}

func registerAPIRoutes(router *gin.Engine, store *sqlite.Store, sessions *session.Store) {
	authSvc := service.NewAuthService(store, sessions)
	authH := auth.NewAuthHandler(authSvc)
	authH.RegisterRoutes(&router.RouterGroup)

	protected := router.Group("")
	protected.Use(middleware.AuthRequired(sessions, store))

	userH := user.NewUserHandler()
	userH.RegisterRoutes(protected)

	bankingSvc := service.NewBankingService(store)
	bankingH := banking.NewBankingHandler(bankingSvc)
	bankingH.RegisterRoutes(protected)

	postSvc := service.NewPostService(store)
	postH := post.NewPostHandler(postSvc)
	postH.RegisterRoutes(protected)
}

// fe 페이지 캐싱으로 테스트에 혼동이 있어, 별도 처리없이 두시면 될 것 같습니다
func registerStaticRoutes(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/static/") || c.Request.URL.Path == "/" {
			c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}
		c.Next()
	})
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
}
