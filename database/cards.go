package database

import (
	"context"
	"database/sql"
	"time"
)


type CardModel struct {
	DB *sql.DB
}

type Card struct {
	CardID  string	`json:"card_id"`
	UserID  int 	`json:"user_id"`
	Quantity int 	`json:"quantity"`
}

type CardList struct {
	Cards []Card `json:"cards"`
}

func (m *CardModel) InsertOrUpdate(userID int, cardID string, quantity int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Tenta atualizar a quantidade se o registro já existir
	updateQuery := "UPDATE user_cards SET quantity = quantity + ? WHERE user_id = ? AND card_id = ?"
	res, err := m.DB.ExecContext(ctx, updateQuery, quantity, userID, cardID)
	if err != nil {
		return err
	}

	// Verifica quantas linhas foram afetadas pelo UPDATE
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Se nenhuma linha foi atualizada, significa que não existe, então insere um novo
	if rowsAffected == 0 {
		insertQuery := "INSERT INTO user_cards (user_id, card_id, quantity) VALUES (?, ?, ?)"
		_, err = m.DB.ExecContext(ctx, insertQuery, userID, cardID, quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *CardModel) GetAll(userID int) ([]Card, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT card_id, quantity FROM user_cards WHERE user_id = ?"
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		if err := rows.Scan(&card.CardID, &card.Quantity); err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cards, nil
}