package controllers

import (
	"encoding/json"
	"net/http"
	"projet-forum/middleware"
	"projet-forum/models"
	"strconv"
)

// LikeMessage ajoute un like à un message
func LikeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Vérifier si l'utilisateur a déjà réagi à ce message
	reaction, err := models.GetMessageReaction(messageID, user.ID)
	if err == nil {
		if reaction.ReactionType == "like" {
			http.Error(w, "You have already liked this message", http.StatusBadRequest)
			return
		} else if reaction.ReactionType == "dislike" {
			// Supprimer le dislike avant d'ajouter le like
			if err := models.DeleteMessageReaction(messageID, user.ID); err != nil {
				http.Error(w, "Error removing previous reaction", http.StatusInternalServerError)
				return
			}
		}
	}

	// Ajouter le like
	message, err := models.GetMessageByID(messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := message.LikeMessage(); err != nil {
		http.Error(w, "Error liking message", http.StatusInternalServerError)
		return
	}

	// Enregistrer la réaction
	reaction = &models.MessageReaction{
		MessageID:    messageID,
		UserID:       user.ID,
		ReactionType: "like",
	}
	if err := reaction.CreateReaction(); err != nil {
		http.Error(w, "Error saving reaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// DislikeMessage ajoute un dislike à un message
func DislikeMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	messageID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	// Vérifier si l'utilisateur a déjà réagi à ce message
	reaction, err := models.GetMessageReaction(messageID, user.ID)
	if err == nil {
		if reaction.ReactionType == "dislike" {
			http.Error(w, "You have already disliked this message", http.StatusBadRequest)
			return
		} else if reaction.ReactionType == "like" {
			// Supprimer le like avant d'ajouter le dislike
			if err := models.DeleteMessageReaction(messageID, user.ID); err != nil {
				http.Error(w, "Error removing previous reaction", http.StatusInternalServerError)
				return
			}
		}
	}

	// Ajouter le dislike
	message, err := models.GetMessageByID(messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}

	if err := message.DislikeMessage(); err != nil {
		http.Error(w, "Error disliking message", http.StatusInternalServerError)
		return
	}

	// Enregistrer la réaction
	reaction = &models.MessageReaction{
		MessageID:    messageID,
		UserID:       user.ID,
		ReactionType: "dislike",
	}
	if err := reaction.CreateReaction(); err != nil {
		http.Error(w, "Error saving reaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
