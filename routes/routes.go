package routes

import (
	"gopost-be/controllers"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Enable CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Adjust to match your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := r.Group("/article")
	{
		api.POST("", controllers.CreateArticle)
		api.GET("/list/:limit/:offset", controllers.GetArticles)
		api.GET("/:id", controllers.GetArticleByID)
		api.PUT("/:id", controllers.UpdateArticle)
		api.DELETE("/:id", controllers.DeleteArticle)
	}

	return r
}
