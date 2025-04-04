package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/andrepassosb/poke-trade-api/database"

	"github.com/gin-gonic/gin"
)

type addCardRequest struct {
	UserID   int    `json:"user_id"`
	CardID string `json:"card_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

func (app *application) updateCardList(c *gin.Context) {
	var request addCardRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)

	card := database.Card{
		UserID:   user.Id,
		CardID:   request.CardID,
		Quantity: request.Quantity,
	}

	err := app.models.Cards.InsertOrUpdate(card.UserID, card.CardID, card.Quantity)
	if err != nil {
		log.Printf("Error inserting card: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card updated successfully"})
}

func (app *application) getAllCards(c *gin.Context) {
	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		log.Println("Error converting id to integer:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cards, err := app.models.Cards.GetAll(intID)
	if err != nil {
		log.Printf("Error getting cards: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cards": cards})
}