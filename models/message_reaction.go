package models

import (
	"projet-forum/config"
	"time"
)

type MessageReaction struct {
	ID           int       `json:"id"`
	MessageID    int       `json:"message_id"`
	UserID       int       `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
	CreatedAt    time.Time `json:"created_at"`
}

func (MessageReaction) TableName() string {
	return "message_reactions"
}

func (r *MessageReaction) CreateReaction() error {
	query := `
		INSERT INTO message_reactions (message_id, user_id, reaction_type, created_at)
		VALUES (?, ?, ?, ?)`

	result, err := config.DB.Exec(
		query,
		r.MessageID,
		r.UserID,
		r.ReactionType,
		time.Now(),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	r.ID = int(id)
	return nil
}

func GetMessageReaction(messageID, userID int) (*MessageReaction, error) {
	reaction := &MessageReaction{}
	query := `SELECT * FROM message_reactions WHERE message_id = ? AND user_id = ?`
	err := config.DB.QueryRow(query, messageID, userID).Scan(
		&reaction.ID,
		&reaction.MessageID,
		&reaction.UserID,
		&reaction.ReactionType,
		&reaction.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return reaction, nil
}

func DeleteMessageReaction(messageID, userID int) error {
	query := `DELETE FROM message_reactions WHERE message_id = ? AND user_id = ?`
	_, err := config.DB.Exec(query, messageID, userID)
	return err
}
