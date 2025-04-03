package api

import (
	"github.com/andrepassosb/poke-trade-api/database"

	"github.com/gin-gonic/gin"
)

func (app *application) GetUserFromContext(c *gin.Context) *database.User {
	contextUser, exist := c.Get("user")
	if !exist {
		return &database.User{}
	}
	user, ok := contextUser.(*database.User)
	if !ok {
		return &database.User{}
	}

	return user
}