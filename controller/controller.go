package controller

import (
	"database/sql"
	"my-project/config"
	"my-project/model"

	"github.com/gin-gonic/gin"
)

func Addtodo(c *gin.Context) {
	var newTodo model.Todo

	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(400, gin.H{"error": "kiritma xato kiritildi"})
		return
	}
	result, err := config.DB.Exec("INSERT INTO todos(work,time,completed) VALUES(?,?,?)", newTodo.Work, newTodo.Time, newTodo.Completed)
	if err != nil {
		c.JSON(400, gin.H{"error": "malumot qabul qilinmadi"})
		return
	}
	lastId, _ := result.LastInsertId()
	newTodo.Id = int(lastId)
	c.JSON(201, newTodo)

}
func GetallTodo(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id,work,time,completed from todos")
	if err != nil {
		c.JSON(400, gin.H{"error": "db dan dan xato qaytdi"})
		return
	}
	defer rows.Close()
	todo := []model.Todo{}

	for rows.Next() {
		var t model.Todo
		err = rows.Scan(&t.Id, &t.Work, &t.Time, &t.Completed)
		if err = rows.Err(); err != nil {
			c.JSON(400, gin.H{"error": "db dan malumotni olib bulmadi"})
			return
		}
		todo = append(todo, t)
	}

	c.JSON(200, todo)

}
func GetbuyidTodo(c *gin.Context) {
	id := c.Param("id")
	var t model.Todo
	err := config.DB.QueryRow("SELECT id,work,time,completed FROM todos WHERE id=?", id).Scan(&t.Id, &t.Work, &t.Time, &t.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(400, gin.H{"error": "bunday id lik malumot yoq"})
			return
		}
		c.JSON(400, gin.H{"error": ""})
		return
	}
	c.JSON(200, t)
}
func Put_todo(c *gin.Context) {
	id := c.Param("id")
	var newTodo model.Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(400, gin.H{"error": "malumot to'liq emas"})
		return
	}
	result, err := config.DB.Exec("UPDATE todos SET work=?,time=?,completed=? WHERE id=?", newTodo.Work, newTodo.Time, newTodo.Completed, id)
	if err != nil {
		c.JSON(500, gin.H{"error": "db tizimida da xatolik bor "})
		return
	}
	count, err := result.RowsAffected()
	if count == 0 {
		c.JSON(404, gin.H{"error": "Yangilanish uchun munday id topilmadi"})
		return
	}
	c.JSON(200, gin.H{"message": "todo o'zgartirish muaffaqiyatli bajarildi"})
}
func Delete_todo(c *gin.Context) {
	id := c.Param("id")
	result, err := config.DB.Exec("DELETE FROM todos WHERE id=?", id)
	if err != nil {
		c.JSON(400, gin.H{"error": "db tizimi vaqtinchaloik ishlamayapti"})
		return
	}
	count, err := result.RowsAffected()
	if count == 0 {
		c.JSON(404, gin.H{"error": "bunday id mavjud emas"})
		return
	}
	c.JSON(200, gin.H{"message": "malumot o'chirildi"})
}
