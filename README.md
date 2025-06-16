# Forum API

Une API RESTful pour un forum de discussion développée en Go.

## Prérequis

- Go 1.16 ou supérieur
- MySQL 5.7 ou supérieur (via XAMPP recommandé)
- Git

## Installation

1. Cloner le dépôt :
```bash
git clone https://github.com/votre-username/projet-forum.git
cd projet-forum
```

2. Installer les dépendances :
```bash
go mod download
```

3. Configurer la base de données :
- Démarrer XAMPP et activer MySQL
- Importer le fichier `schema.sql` dans phpMyAdmin ou exécuter :
```bash
mysql -u root < schema.sql
```

4. Configurer les variables d'environnement (optionnel) :
```bash
# Windows (PowerShell)
$env:DATABASE_URL="root:@tcp(localhost:3306)/forum_db?parseTime=true"
$env:PORT="8080"

# Linux/Mac
export DATABASE_URL="root:@tcp(localhost:3306)/forum_db?parseTime=true"
export PORT="8080"
```

## Démarrage

```bash
go run main.go
```

Le serveur démarrera sur le port 8080 par défaut.

## API Endpoints

### Routes publiques

- `POST /api/register` - Inscription d'un nouvel utilisateur
- `POST /api/login` - Connexion d'un utilisateur
- `GET /api/threads` - Liste des fils de discussion
- `GET /api/threads/{id}` - Détails d'un fil de discussion
- `GET /api/threads/{id}/messages` - Messages d'un fil de discussion
- `GET /api/search` - Recherche de fils de discussion
- `GET /api/categories` - Liste des catégories
- `GET /api/categories/{id}` - Détails d'une catégorie

### Routes protégées (nécessite un token JWT)

- `GET /api/users/me` - Informations de l'utilisateur connecté
- `PUT /api/users/me` - Mise à jour du profil
- `PUT /api/users/me/password` - Changement de mot de passe
- `POST /api/threads` - Création d'un fil de discussion
- `PUT /api/threads/{id}` - Mise à jour d'un fil de discussion
- `DELETE /api/threads/{id}` - Suppression d'un fil de discussion
- `POST /api/threads/{id}/messages` - Création d'un message
- `PUT /api/messages/{id}` - Mise à jour d'un message
- `DELETE /api/messages/{id}` - Suppression d'un message
- `POST /api/messages/{id}/like` - Like d'un message
- `POST /api/messages/{id}/dislike` - Dislike d'un message

### Routes admin (nécessite un token JWT admin)

- `POST /api/admin/users/{id}/ban` - Bannir un utilisateur
- `POST /api/admin/users/{id}/unban` - Débannir un utilisateur
- `PUT /api/admin/threads/{id}/status` - Mise à jour du statut d'un fil
- `GET /api/admin/stats` - Statistiques du forum
- `POST /api/admin/categories` - Création d'une catégorie
- `PUT /api/admin/categories/{id}` - Mise à jour d'une catégorie
- `DELETE /api/admin/categories/{id}` - Suppression d'une catégorie

## Structure du projet

```
.
├── controllers/     # Contrôleurs de l'API
├── middleware/      # Middleware (auth, admin)
├── models/         # Modèles de données
├── main.go         # Point d'entrée
├── schema.sql      # Schéma de la base de données
└── README.md       # Documentation
```

## Authentification

L'API utilise JWT (JSON Web Tokens) pour l'authentification. Pour accéder aux routes protégées, incluez le token dans l'en-tête de la requête :

```
Authorization: Bearer <votre_token>
```

## Compte admin par défaut

- Email : admin@forum.com
- Mot de passe : admin123

## Licence

MIT