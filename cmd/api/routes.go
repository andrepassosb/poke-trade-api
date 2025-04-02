package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)
		v1.GET("users", app.getAllUsers)
		v1.GET("users/:id", app.getUserByID)
		
	}


	return g
}