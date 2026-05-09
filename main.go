package main

import (
	"my-project/config"
	"my-project/controller"
	"my-project/middlware"

	_ "my-project/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Mening To-Do API-im
// @version         1.0
// @description     Bu Go va Gin yordamida to-do API loyihasi.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Jamshid
// @contact.url    https://t.me/Xoliqov_jamshid

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in                         header
// @name                       Authorization

func main() {
	config.Connect()
	defer config.DB.Close()

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	v1 := r.Group("/api")
	{
		v1.GET("/", homeHandler)
		v1.POST("/login", controller.Login)
		v1.POST("/register", controller.Register)
	}
	protected := r.Group("/api")
	protected.Use(middlware.Authmidllware())
	{

		protected.GET("/todolarniolish", controller.GetallTodo)
		protected.GET("/todoniolish/:id", controller.GetbuyidTodo)
		protected.PUT("/todoniyangilash/:id", controller.Put_todo)
		protected.DELETE("/deletetodo/:id", controller.Delete_todo)
		protected.POST("/todoqushish", controller.Addtodo)

	}
	r.Run("0.0.0.0:8080")
}
func homeHandler(c *gin.Context) {

}
