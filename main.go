package main

import (
	"my-project/config"
	"my-project/controller"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func homeHandler(c *gin.Context) {

}

func main() {
	config.Connect()
	defer config.DB.Close()

	r := gin.Default()

	r.GET("/", homeHandler)
	r.GET("/todolarniolish", controller.GetallTodo)
	r.GET("/todoniolish/:id", controller.GetbuyidTodo)
	r.PUT("/todoniyangilash/:id", controller.Put_todo)
	r.DELETE("/deletetodo/:id", controller.Delete_todo)
	r.POST("/todoqushish", controller.Addtodo)
	r.Run("0.0.0.0:8080")

}
