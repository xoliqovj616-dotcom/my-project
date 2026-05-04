package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Todo struct {
	Id        int    `json:"id"`
	Work      string `json:"work"`
	Time      string `json:"time"`
	Completed bool   `json:"completed"`
}

var db *sql.DB

func Init() {
	var err error
	db, err = sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("db ga ulanib bulmadi")
	}
	createTable := `CREATE TABLE IF NOT EXISTS todos(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	work TEXT NOT NULL,
    time  TEXT,
	completed INTEGER DEFAULT 0
	);`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("db ulandi")
}
func homeHandler(c *gin.Context) {

}
func Addtodo(c *gin.Context) {
	var newTodo Todo

	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(400, gin.H{"error": "kiritma xato kiritildi"})
		return
	}
	result, err := db.Exec("INSERT INTO todos(work,time,completed) VALUES(?,?,?)", newTodo.Work, newTodo.Time, newTodo.Completed)
	if err != nil {
		c.JSON(400, gin.H{"error": "malumot qabul qilinmadi"})
		return
	}
	lastId, _ := result.LastInsertId()
	newTodo.Id = int(lastId)
	c.JSON(201, newTodo)

}

func GetallTodo(c *gin.Context) {
	rows, err := db.Query("SELECT id,work,time,completed from todos")
	if err != nil {
		c.JSON(400, gin.H{"error": "db dan malumot topilmadi"})
		return
	}
	defer rows.Close()
	todo := []Todo{}

	for rows.Next() {
		var t Todo
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
	var t Todo
	err := db.QueryRow("SELECT id,work,time,completed FROM todos WHERE id=?", id).Scan(&t.Id, &t.Work, &t.Time, &t.Completed)
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

}
func Delete_todo(c *gin.Context) {

}

func main() {
	Init()
	defer db.Close()

	r := gin.Default()

	r.GET("/", homeHandler)
	r.GET("/todolarniolish", GetallTodo)
	r.GET("/todoniolish/:id", GetbuyidTodo)
	r.PUT("/todoniyangilash/:id", Put_todo)
	r.DELETE("/deletetodo/:id", Delete_todo)
	r.POST("/todoqushish", Addtodo)
	r.Run("0.0.0.0:8080")

}
