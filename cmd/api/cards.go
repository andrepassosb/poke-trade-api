package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardUpdate struct {
	UserID   int    `json:"user_id"`
	CardID   string `json:"card_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required"`
}

// Inicia o worker de atualização dos cards
func (app *application) startCardUpdateWorker() {
	go func() {
		for update := range app.cardUpdateQueue {
			log.Printf("Processing card update: %+v", update)
			err := app.models.Cards.InsertOrUpdate(update.UserID, update.CardID, update.Quantity)
			if err != nil {
				log.Printf("Error processing card update: %v", err)
			}
		}
	}()
}

// Adiciona a atualização do card à fila de updates
func (app *application) updateCardList(c *gin.Context) {
	var request CardUpdate

	if err := c.ShouldBindJSON(&request); err != nil {
		if err.Error() == "EOF" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is empty"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := app.GetUserFromContext(c)

	update := CardUpdate{
		UserID:   user.Id,
		CardID:   request.CardID,
		Quantity: request.Quantity,
	}

	// Adiciona a atualização na fila
	app.cardUpdateQueue <- update

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

	c.JSON(http.StatusOK, gin.H{
		"user_id": intID,
		"cards":   cards,
	})
}
