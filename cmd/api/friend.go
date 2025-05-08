package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/andrepassosb/poke-trade-api/database"

	"github.com/gin-gonic/gin"
)

type addFriendRequest struct {
	FriendID int `json:"friend_id" binding:"required"`
}

func (app *application) addFriend(c *gin.Context) {
	var request addFriendRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)


	friend := database.Friend{
		FriendID:  request.FriendID,
	}

	err := app.models.Friends.Insert(user.Id, friend.FriendID)
	if err != nil {
		log.Printf("Error inserting friend: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend added successfully"})
}

func (app *application) getFriends(c *gin.Context) {
	user := app.GetUserFromContext(c)

	friends, err := app.models.Friends.GetAll(user.Id)
	if err != nil {
		log.Printf("Error getting friends: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

func (app *application) getFriendByID(c *gin.Context) {
	user := app.GetUserFromContext(c)

	id := c.Param("friendId")
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Error converting id to integer:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	friend, err := app.models.Friends.GetByID(intID, user.Id)
	if err != nil {
		log.Println("Error getting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if friend == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend not found"})
		return
	}


	c.JSON(http.StatusOK, database.User{
		Id:       friend.FriendID,
		Username: friend.Username, 
	})
	
}

func (app *application) deleteFriend(c *gin.Context) {
	id := c.Param("friendId")
	friendID, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Error converting id to integer:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user := app.GetUserFromContext(c)

	err = app.models.Friends.Delete(user.Id, friendID)
	if err != nil {
		log.Printf("Error deleting friend: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Friend deleted successfully"})
}