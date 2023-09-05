package controllers

import (
	helper "github.com/pedromessetti/auth_project/helpers"
	models "github.com/pedromessetti/auth_project/models"
	"github.com/pedromessetti/auth_project/database"

	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = "Invalid email or password"
		check = false
	}
	return check, msg
}

func handleError(c *gin.Context, status int, errMsg string) {
	c.JSON(status, gin.H{"error": errMsg})
}

// The function handles the signup process for a user, including validation, checking for
// existing users, generating tokens, and inserting the user into the database.
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			handleError(c, http.StatusBadRequest, err.Error())
			return
		}

		validationError := validate.Struct(user)
		if validationError != nil {
			handleError(c, http.StatusBadRequest, validationError.Error())
			return
		}

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			handleError(c, http.StatusInternalServerError, "Error occured while checking for the email")
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			handleError(c, http.StatusInternalServerError, "Error occured while checking for the phone number")
		}

		if count > 0 {
			handleError(c, http.StatusInternalServerError, "User already exists")
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		token, refreshToken, err := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Error occured while generating tokens")
			return
		}
		user.Token = &token
		user.Refresh_token = &refreshToken

		insertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			handleError(c, http.StatusInternalServerError, "Error occured while inserting user")
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, insertionNumber)
	}
}

// The function handles the login process by verifying the user's email and password, generating
// tokens, and updating the tokens in the database.
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		
		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			handleError(c, http.StatusBadRequest, err.Error())
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Invalid email or password")
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !passwordIsValid {
			handleError(c, http.StatusInternalServerError, msg)
			return
		}

		if foundUser.Email == nil {
			handleError(c, http.StatusInternalServerError, "Invalid email or password")
			return
		}

		token, refreshToken, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Error occured while generating tokens")
			return
		}

		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Error occured while updating tokens")
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

// Retrieves a list of users from the database and returns the result as a JSON response.
func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := helper.CheckUserType(c, "ADMIN")
		if err != nil {
			handleError(c, http.StatusUnauthorized, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex, err := strconv.Atoi(c.Query("startIndex"))
		if err != nil || startIndex < 0 {
			startIndex = 0
		}

		// The match stage is used to filter documents based on certain criteria.
		// In this case, the match stage is specifying an empty filter,
		// which means it will match all documents in the collection.
		matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
		// Groups the documents based on a specified key and performs an operation on each group.
		// In this case, the documents are grouped by the `_id` field with the value "null".
		groupStage := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "null"},
				{Key: "total_count", Value: bson.D{
					{Key: "$sum", Value: 1},
				}},
				{Key: "data", Value: bson.D{
					{Key: "$push", Value: "$$ROOT"},
				}},
			}},
		}	
		// Stage used to reshape the output of the previous stages.
		// In this case, it is used to project the desired fields
		// and modify the structure of the output.
		projectStage := bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "_id", Value: 0},
				{Key: "total_count", Value: 1},
				{Key: "user_items", Value: bson.D{
					{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}},
				}},
			}},
		}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			handleError(c, http.StatusInternalServerError, err.Error())
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			return
		}
		c.JSON(http.StatusOK, allUsers[0])
	}
}

// Retrieves a user from a database based on the user ID provided in the request parameter.
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			handleError(c, http.StatusUnauthorized, err.Error())
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			handleError(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
