package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"projet-forum/middleware"
	"projet-forum/models"
	"strconv"
)

type LikeController struct {
	DB *sql.DB
}

func (c *LikeController) LikeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	reaction, err := models.GetMessageReaction(c.DB, messageID, claims.UserID)
	if err == nil {
		if reaction.ReactionType == "like" {
			http.Error(w, "You have already liked this message", http.StatusBadRequest)
			return
		} else if reaction.ReactionType == "dislike" {
			if err := models.DeleteMessageReaction(c.DB, messageID, claims.UserID); err != nil {
				http.Error(w, "Error removing previous reaction", http.StatusInternalServerError)
				return
			}
		}
	}

	message, err := models.GetMessage(c.DB, messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := models.AddMessageReaction(c.DB, messageID, claims.UserID, "like"); err != nil {
		http.Error(w, "Error liking message", http.StatusInternalServerError)
		return
	}

	message.Likes++
	if err := message.UpdateMessage(c.DB); err != nil {
		http.Error(w, "Error updating message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func (c *LikeController) DislikeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	reaction, err := models.GetMessageReaction(c.DB, messageID, claims.UserID)
	if err == nil {
		if reaction.ReactionType == "dislike" {
			http.Error(w, "You have already disliked this message", http.StatusBadRequest)
			return
		} else if reaction.ReactionType == "like" {
			if err := models.DeleteMessageReaction(c.DB, messageID, claims.UserID); err != nil {
				http.Error(w, "Error removing previous reaction", http.StatusInternalServerError)
				return
			}
		}
	}

	message, err := models.GetMessage(c.DB, messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := models.AddMessageReaction(c.DB, messageID, claims.UserID, "dislike"); err != nil {
		http.Error(w, "Error disliking message", http.StatusInternalServerError)
		return
	}

	message.Dislikes++
	if err := message.UpdateMessage(c.DB); err != nil {
		http.Error(w, "Error updating message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
