package helper

import (
	"errors"
	"github.com/gin-gonic/gin"
)

func CheckUserType(c *gin.Context, role string) (err error) {
	// Get the user type from the jwt
	userType := c.GetString("user_type")

	err = nil

	// If the user type does not match the user type from the jwt
	if userType != role {
		err = errors.New("Unauthorized")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	// Get the user type from the jwt
	userType := c.MustGet("userType").(string)

	// Get the user id from the jwt
	uid := c.MustGet("uid").(string)

	err = nil

	// If the user type is not "admin" and the user id does not match the user id from the jwt
	if userType == "USER" && uid != userId {
		err = errors.New("Unauthorized")
		return err
	}

	err = CheckUserType(c, userType)

	return err
}