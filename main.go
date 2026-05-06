package main

import (
	"my-project/config"
	"my-project/controller"
	"my-project/middlware"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func homeHandler(c *gin.Context) {

}

func main() {
	config.Connect()
	defer config.DB.Close()

	r := gin.Default()
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
