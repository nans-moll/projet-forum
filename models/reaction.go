package models

import (
	"database/sql"
	"time"
)

// MessageReaction représente une réaction (like/dislike) à un message
type MessageReaction struct {
	ID           int64     `json:"id"`
	MessageID    int64     `json:"message_id"`
	UserID       int64     `json:"user_id"`
	ReactionType string    `json:"reaction_type"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetMessageReaction récupère une réaction à un message
func GetMessageReaction(db *sql.DB, messageID, userID int64) (*MessageReaction, error) {
	query := `
		SELECT id, message_id, user_id, reaction_type, created_at
		FROM message_reactions
		WHERE message_id = ? AND user_id = ?
	`
	reaction := &MessageReaction{}
	err := db.QueryRow(query, messageID, userID).Scan(
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

// DeleteMessageReaction supprime une réaction à un message
func DeleteMessageReaction(db *sql.DB, messageID, userID int64) error {
	query := `
		DELETE FROM message_reactions
		WHERE message_id = ? AND user_id = ?
	`
	_, err := db.Exec(query, messageID, userID)
	return err
}

// AddMessageReaction ajoute une réaction à un message
func AddMessageReaction(db *sql.DB, messageID, userID int64, reactionType string) error {
	query := `
		INSERT INTO message_reactions (message_id, user_id, reaction_type, created_at)
		VALUES (?, ?, ?, CURRENT_TIMESTAMP)
	`
	_, err := db.Exec(query, messageID, userID, reactionType)
	return err
}

// UpdateMessageReactionCount met à jour le compteur de likes/dislikes d'un message
func UpdateMessageReactionCount(db *sql.DB, messageID int64) error {
	query := `
		UPDATE messages m
		SET 
			likes = (SELECT COUNT(*) FROM message_reactions WHERE message_id = m.id AND reaction_type = 'like'),
			dislikes = (SELECT COUNT(*) FROM message_reactions WHERE message_id = m.id AND reaction_type = 'dislike')
		WHERE id = ?
	`
	_, err := db.Exec(query, messageID)
	return err
}
