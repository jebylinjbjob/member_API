package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"member_API/auth"
	"member_API/controllers"
	"member_API/graphql"
)

// SetupRouter registers API routes on the provided Gin engine.
func SetupRouter(r *gin.Engine) {
	r.GET("/Hello", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, RESTful API!"})
	})

	// 公開路由 - 不需要認證
	public := r.Group("/api/v1")
	{
		// 認證相關路由
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)
	}

	// GraphQL endpoint（可選擇是否需要認證）
	r.Any("/graphql", func(c *gin.Context) {
		h := graphql.GetHandler()
		if h == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "GraphQL handler not initialized"})
			return
		}
		h.ServeHTTP(c.Writer, c.Request)
	})

	// 受保護的路由 - 需要認證
	protected := r.Group("/api/v1")
	protected.Use(auth.AuthMiddleware()) // 添加認證中間件
	{
		protected.GET("/users", controllers.GetUsers)
		protected.GET("/user/:id", func(c *gin.Context) {
			controllers.GetUserByID(c)
		})
		protected.GET("/profile", controllers.GetProfile) // 獲取當前用戶信息
	}
}
