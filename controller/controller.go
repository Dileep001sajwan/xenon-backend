package controller

import (
	"context"
	"net/http"
	"time"
	"xenon-backend/database"
	"xenon-backend/heleprs"
	"xenon-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		
		var existingUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)

		
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
			return
		} else if err != mongo.ErrNoDocuments {
			
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking existing user"})
			return
		}

		hashedPassword, err := heleprs.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
			return
		}
		user.Password = hashedPassword

		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"inserted_id": result.InsertedID})
	}
}

// Login API
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.LoginData

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var existingUser models.User
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser)

		if err == mongo.ErrNoDocuments {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
			return
		} else if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking user"})
			return
		}

		ok, check := heleprs.VerifyPassword(existingUser.Password, user.Password)

		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": check})
			return
		}

		token, err := heleprs.GenerateAccessToken(existingUser.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while generating access token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"access_token": token})
	}
}
