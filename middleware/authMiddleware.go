package middleware

import (
	"net/http"

	helper "github.com/pedromessetti/auth_project/helpers"

	"github.com/gin-gonic/gin"
)

// Middleware function that validates a token in the request header and sets user information in the context if the token is valid.
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized [missing token]"})
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unauthorized [invalid token]"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("user_type", claims.User_type)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
