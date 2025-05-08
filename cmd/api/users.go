package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/andrepassosb/poke-trade-api/database"
	"github.com/gin-gonic/gin"
)

type UsersResponse struct {
	Users []database.User `json:"users"`
}

type UserResponse struct {
	User database.User `json:"user"`
}

func (app *application) getAllUsers(c *gin.Context) {
	var usersResponse UsersResponse

	users, err := app.models.Users.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	usersResponse.Users = make([]database.User, len(users))
	for i, user := range users {
		usersResponse.Users[i] = *user
	}

	c.JSON(http.StatusOK, usersResponse)
}

func (app *application) getUserByID(c *gin.Context) {
	var userResponse UserResponse

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Error converting id to integer:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := app.models.Users.GetByID(intID)
	if err != nil {
		log.Println("Error getting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	userResponse.User = *user

	c.JSON(http.StatusOK, userResponse)
}