# projet-forum
Configuration de la base de données PostgreSQL avec les tables nécessaires
Modèles pour les utilisateurs, fils de discussion et messages
Middleware d'authentification avec JWT
Contrôleurs pour :
Authentification (inscription/connexion)
Gestion des fils de discussion
Gestion des messages
Administration du forum
Routes API RESTful pour toutes les fonctionnalités
Pour utiliser l'API, voici les principales routes disponibles :
Authentification :
POST /api/auth/register - Inscription
POST /api/auth/login - Connexion
Fils de discussion :
POST /api/threads - Créer un fil
GET /api/threads/ - Voir un fil
GET /api/threads/tag - Lister les fils par tag
PUT /api/threads/update - Modifier un fil
DELETE /api/threads/delete - Supprimer un fil
Messages :
POST /api/messages - Créer un message
GET /api/messages/ - Voir les messages d'un fil
PUT /api/messages/update - Modifier un message
DELETE /api/messages/delete - Supprimer un message
POST /api/messages/like - Liker un message
POST /api/messages/dislike - Disliker un message
Administration :
POST /api/admin/ban - Bannir un utilisateur
POST /api/admin/unban - Débannir un utilisateur
PUT /api/admin/thread/status - Modifier le statut d'un fil
GET /api/admin/stats - Voir les statistiques du forum
Pour démarrer le projet :
Installer PostgreSQL et créer la base de données
Exécuter le script SQL dans database/init.sql
Installer les dépendances Go : go mod tidy
Lancer le serveur : go run main.go
Le serveur démarrera sur http://localhost:8080.