package middlware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authmidllware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokennikutish := c.GetHeader("Authorization")
		if tokennikutish == "" {
			c.JSON(400, gin.H{"error": "token topilmadi"})
			c.Abort()
			return
		}
		tokentozalash := strings.TrimPrefix(tokennikutish, "Bearer ")
		token, err := jwt.Parse(tokentozalash, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("jwtSecret")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "token xato yoki muddati tugagan"})
			c.Abort()
			return
		}
		c.Next()
	}

}
