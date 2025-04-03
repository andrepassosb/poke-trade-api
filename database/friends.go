package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type FriendModel struct {
	DB *sql.DB
}

type Friend struct {
    FriendID int  	`json:"id"`
    Username string	`json:"username"`
}



type FriendList struct {
	Friends []Friend `json:"friends"`
}

func (m *FriendModel) Insert(userID, friendID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	queryCheck := "SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)"
	errCheck := m.DB.QueryRowContext(ctx, queryCheck, friendID).Scan(&exists)
	if errCheck != nil {
		return errCheck
	}
	if !exists {
		return fmt.Errorf("Friend ID %d does not exist", friendID)
	}

	
	var friendshipExists bool
	queryFriendship := "SELECT EXISTS(SELECT 1 FROM friendships WHERE user_id = ? AND friend_id = ?)"
	errFriendship := m.DB.QueryRowContext(ctx, queryFriendship, userID, friendID).Scan(&friendshipExists)
	if errFriendship != nil {
		return errFriendship
	}
	if friendshipExists {
		return fmt.Errorf("friendship between user %d and friend %d already exists", userID, friendID)
	}

	if userID == friendID {
		return fmt.Errorf("user cannot be friends with themselves")
	}

	query := "INSERT INTO friendships (user_id, friend_id) VALUES (?, ?)"
	_, err := m.DB.ExecContext(ctx, query, userID, friendID)
	return err
}

func (m *FriendModel) GetAll(userID int) ([]*Friend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	fmt.Println("Fetching friends for user ID:", userID)

	query := "SELECT f.friend_id, u.username FROM friendships f JOIN users u ON f.friend_id = u.id WHERE f.user_id = ?"
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*Friend
	for rows.Next() {
		var friend Friend
		if err := rows.Scan( &friend.FriendID, &friend.Username); err != nil {
			return nil, err
		}
		friends = append(friends, &friend)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}

func (m *FriendModel) Delete(userID, friendID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM friendships WHERE user_id = ? AND friend_id = ?"
	_, err := m.DB.ExecContext(ctx, query, userID, friendID)
	if err != nil {
		return err
	}
	return nil
}

func (m *FriendModel) GetByID(friendID, userID int) (*Friend, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var friend Friend

	query := `
			SELECT u.id, u.username 
			FROM users u
			JOIN friendships f ON u.id = f.friend_id
			WHERE f.user_id = ? AND f.friend_id = ?;
		`

	err := m.DB.QueryRowContext(ctx, query, userID,friendID).Scan(&friend.FriendID, &friend.Username)


	if err != nil {
		fmt.Println("Error fetching friend:", err)
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &Friend{
		FriendID: friend.FriendID,
		Username: friend.Username,
	}, nil
}
