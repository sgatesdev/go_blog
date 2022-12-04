package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	user "samgates.io/server/database/models"
	token "samgates.io/server/utils"
)

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input LoginInput

	// validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// check to make sure not already in database
	userMongo, err := user.Find(input.Username)
	if userMongo.Username != "" || (err != nil && err.Error() != "mongo: no documents in result") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists or there was a database error. Please try again"})
		return
	}

	// try to add to database
	createError := user.Create(input.Username, input.Password)

	if createError != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to create user. Please try again"})
		return
	}

	// generate token, return it
	token, err := token.Generate(input.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func Login(c *gin.Context) {
	var input LoginInput

	// validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// authenticate user
	userMongo, err := user.Authenticate(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	// generate token, return it
	token, err := token.Generate(userMongo.Username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"token": "could not generate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func CurrentUser(c *gin.Context) {
	// verify that token can be used to accurately identify user
	username, err := token.ExtractTokenUsername(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// lookup user in DB
	mongoUser, err := user.Find(username)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to identify user with that token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": mongoUser.Username})
}
