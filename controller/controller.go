package controller

import (
	"database/sql"
	"my-project/config"
	"my-project/model"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
// @Summary      Tizimga kirish
// @Description  Username va password orqali JWT token oladi
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security      ApiKeyAuth
// @Param        input  body      model.User  true  "Login ma'lumotlari"
// @Success      200    {object}  map[string]string "token qaytadi"
// @Router       /login [post]
func Login(c *gin.Context) {
	var Input model.User
	var dbUser model.User

	if err := c.ShouldBindJSON(&Input); err != nil {
		c.JSON(400, gin.H{"error": "Malumot noto'h'ri"})
		return
	}
	err := config.DB.QueryRow("SELECT id,username,password FROM users WHERE username=?", Input.Username).Scan(&dbUser.Id, &dbUser.Username, &dbUser.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Foydalanuvchi topilmadi yoki parol xato"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(Input.Password))
	if err != nil {
		c.JSON(401, gin.H{"error": "Foydalanuvchi topilmadi yoki parol xato kiritildi"})
		return
	}
	var Secret = []byte(os.Getenv("jwtSecret"))
	createToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": dbUser.Id,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := createToken.SignedString(Secret)
	if err != nil {
		c.JSON(500, gin.H{"error": "token yaratishda muammo bor"})
		return
	}
	c.JSON(201, gin.H{
		"message": "xush kelipsiz",
		"token":   token,
	})

}

// Register godoc
// @Summary      Ro'yxatdan o'tish
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security      ApiKeyAuth
// @Param        user  body      model.User  true  "Yangi foydalanuvchi"
// @Success      201   {object}  map[string]string
// @Router       /register [post]
func Register(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "user qushish tizimi hosircha ishlamayapti keyinroq urnalab ko'ring"})
		return

	}
	hashpasword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "parolni saqlash qismida xatolik mavjud"})
		return
	}
	_, err = config.DB.Exec("INSERT INTO users(username,password) VALUES(?,?)", user.Username, string(hashpasword))
	if err != nil {
		c.JSON(400, gin.H{"error": "bu login band qayta urinib ko'ring"})
		return
	}
	c.JSON(201, gin.H{"message": "foydalanuvchi muaffaqiyatli yaratildi"})
}

// @Summary      Yangi todo qo'shish
// @Tags         todo
// @Accept       json
// @Produce      json
// @Security      ApiKeyAuth
// @Param        task  body      model.Todo  true  "Yangi vazifa ma'lumotlari"
// @Success      201   {object}  model.Todo
// @Router       /todoqushish [post]
func Addtodo(c *gin.Context) {
	id, _ := c.Get("userID")
	UserId := id.(int)
	var newTodo model.Todo

	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(400, gin.H{"error": "kiritma xato kiritildi"})
		return
	}
	result, err := config.DB.Exec("INSERT INTO todos(work,time,completed,user_id) VALUES(?,?,?,?)", newTodo.Work, newTodo.Time, newTodo.Completed, UserId)
	if err != nil {
		c.JSON(400, gin.H{"error": "malumot qabul qilinmadi"})
		return
	}
	lastId, _ := result.LastInsertId()
	newTodo.Id = int(lastId)
	c.JSON(201, newTodo)

}

// @Summary      Barcha todolarni olish
// @Description  Foydalanuvchiga tegishli barcha vazifalar ro'yxatini qaytaradi
// @Tags         todo
// @Accept       json
// @Produce      json
// @Security      ApiKeyAuth
// @Success      200  {array}   model.Todo
// @Router       /todolarniolish [get]
func GetallTodo(c *gin.Context) {
	id, _ := c.Get("userID")

	idInt := id.(int)
	rows, err := config.DB.Query("SELECT id,work,time,completed,user_id FROM todos WHERE user_id=?", idInt)
	if err != nil {
		c.JSON(400, gin.H{"error": "db dan dan xato qaytdi"})
		return
	}
	defer rows.Close()
	todo := []model.Todo{}

	for rows.Next() {
		var t model.Todo
		err = rows.Scan(&t.Id, &t.Work, &t.Time, &t.Completed, &t.User_id)
		if err != nil {
			c.JSON(400, gin.H{"error": "db dan malumotni olib bulmadi"})
			return
		}
		todo = append(todo, t)
	}

	c.JSON(200, todo)

}

// @Summary      Bitta todoni olish
// @Description  ID bo'yicha bitta vazifani batafsil ko'rish
// @Tags         todo
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  model.Todo
// @Security      ApiKeyAuth
// @Router       /todoniolish/{id} [get]
func GetbuyidTodo(c *gin.Context) {
	id := c.Param("id")
	userid, _ := c.Get("userID")
	InduserId := userid.(int)
	var t model.Todo
	err := config.DB.QueryRow("SELECT id,work,time,completed, user_id FROM todos WHERE id=? AND user_id=?", id, InduserId).Scan(&t.Id, &t.Work, &t.Time, &t.Completed, &t.User_id)
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

// Put_todo godoc
// @Summary      Todoni yangilash
// @Tags         todo
// @Param        id    path      int         true  "Todo ID"
// @Param        todo  body      model.Todo  true  "Yangi ma'lumotlar"
// @Success      200   {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /todoniyangilash/{id} [put]
func Put_todo(c *gin.Context) {

	id := c.Param("id")

	userid, _ := c.Get("userID")
	IntUserid := userid.(int)
	var newTodo model.Todo
	if err := c.ShouldBindJSON(&newTodo); err != nil {
		c.JSON(400, gin.H{"error": "malumot to'liq emas"})
		return
	}
	result, err := config.DB.Exec("UPDATE todos SET work=?,time=?,completed=? WHERE id=? AND user_id=?", newTodo.Work, newTodo.Time, newTodo.Completed, id, IntUserid)
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

// Delete_todo godoc
// @Summary      Todoni o'chirish
// @Tags         todo
// @Param        id   path      int  true  "Todo ID"
// @Success      200  {object}  map[string]string
// @Security      ApiKeyAuth
// @Router       /deletetodo/{id} [delete]
func Delete_todo(c *gin.Context) {
	id := c.Param("id")
	user_id, _ := c.Get("userID")
	Intid := user_id.(int)
	result, err := config.DB.Exec("DELETE FROM todos WHERE id=? AND  user_id=? ", id, Intid)
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
