# 🌐 Cine Forum Application

Une application web de forum complète développée en Go avec une interface web moderne et une API RESTful.

## 📋 Prérequis

- **Go 1.19 ou supérieur** - [Télécharger Go](https://golang.org/dl/)
- **MySQL 5.7 ou supérieur** - Via XAMPP recommandé
- **Git** - [Télécharger Git](https://git-scm.com/)
- **Navigateur web moderne** (Chrome, Firefox, Safari, Edge)

## 🏗️ Architecture du Projet

### Structure des fichiers
```
projet-forum/
├── controllers/           # Logique métier et contrôleurs
│   ├── auth_controller.go    # Authentification (login/register)
│   ├── user_controller.go    # Gestion des utilisateurs
│   ├── thread_controller.go  # Gestion des discussions
│   ├── stats_controller.go   # Statistiques du forum
│   └── admin_controller.go   # Administration
├── middleware/           # Middleware d'authentification
│   └── auth.go              # Vérification JWT et sessions
├── models/              # Modèles de données
├── routes/              # Configuration des routes
│   ├── api_routes.go        # Routes API REST
│   ├── auth_routes.go       # Routes d'authentification
│   └── admin_routes.go      # Routes d'administration
├── static/              # Fichiers statiques (CSS, JS, images)
│   ├── css/
│   │   └── styles.css       # Styles de l'interface
│   ├── js/
│   │   ├── api.js           # Client API JavaScript
│   │   ├── auth.js          # Gestion de l'authentification côté client
│   │   └── profile.js       # Interface utilisateur du profil
│   └── images/
│       └── default-avatar.png
├── templates/           # Templates HTML
│   ├── index.html           # Page d'accueil
│   ├── threads/
│   │   ├── show.html        # Affichage d'une discussion
│   │   ├── create.html      # Création de discussion
│   │   └── edit.html        # Édition de discussion
│   └── users/
│       └── profile.html     # Page de profil utilisateur
├── uploads/             # Dossier pour les fichiers uploadés
├── .env                 # Variables d'environnement
├── main.go             # Point d'entrée de l'application
├── go.mod              # Dépendances Go
├── go.sum              # Checksums des dépendances
├── schema.sql          # Schéma de la base de données
└── README.md           # Cette documentation
```

### Technologies utilisées

**Backend :**
- **Go (Golang)** - Langage principal
- **Gorilla Mux** - Routeur HTTP
- **MySQL** - Base de données
- **JWT** - Authentification par tokens
- **godotenv** - Gestion des variables d'environnement

**Frontend :**
- **HTML5** - Structure des pages
- **CSS3** - Styles et mise en page
- **JavaScript (Vanilla)** - Interactivité côté client
- **API REST** - Communication client-serveur

## 🚀 Installation et Configuration

### 1. Cloner le projet
```bash
git clone https://github.com/votre-username/projet-forum.git
cd projet-forum
```

### 2. Installer les dépendances Go
```bash
# Initialiser le module Go (si pas déjà fait)
go mod init projet-forum

# Télécharger les dépendances
go mod download

# Si certaines dépendances manquent, les installer :
go get github.com/gorilla/mux
go get github.com/go-sql-driver/mysql
go get github.com/joho/godotenv
go get github.com/dgrijalva/jwt-go
```

### 3. Configurer MySQL avec XAMPP

#### Installation XAMPP
1. **Télécharger XAMPP** depuis [https://www.apachefriends.org/](https://www.apachefriends.org/)
2. **Installer XAMPP** sur votre système
3. **Démarrer XAMPP Control Panel**

#### Configuration de la base de données
1. **Démarrer MySQL** dans XAMPP Control Panel
2. **Ouvrir phpMyAdmin** : [http://localhost/phpmyadmin](http://localhost/phpmyadmin)
3. **Créer la base de données** :
   ```sql
   CREATE DATABASE forum_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```
4. **Importer le schéma** :
   - Cliquer sur la base `forum_db`
   - Aller dans l'onglet "Importer"
   - Sélectionner le fichier `schema.sql`
   - Cliquer sur "Exécuter"

### 4. Configurer les variables d'environnement

Le fichier `.env` est déjà créé, mais vous devez ajuster certaines valeurs :

```bash
# Configuration de la base de données MySQL
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=                    # Mettez votre mot de passe MySQL ici
DB_NAME=forum_db

# Configuration du serveur
SERVER_PORT=8080

# Configuration JWT (changez cette clé en production)
JWT_SECRET=votre_clé_secrète_jwt_super_sécurisée_123456

# Configuration Email (optionnel)
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USER=
EMAIL_PASSWORD=
EMAIL_FROM=

# Mode de développement
APP_ENV=development
DEBUG=true
```

**⚠️ Important :** Remplissez `DB_PASSWORD` avec votre mot de passe MySQL (souvent vide par défaut avec XAMPP).

### 5. Créer les dossiers nécessaires
```bash
# Windows
mkdir uploads
mkdir static\images

# Linux/Mac
mkdir -p uploads
mkdir -p static/images
```

## 🎯 Démarrage de l'application

### 1. Vérifier que MySQL fonctionne
```bash
# Ouvrir XAMPP Control Panel
# Vérifier que MySQL est démarré (voyant vert)
```

### 2. Démarrer l'application
```bash
# Dans le dossier du projet
go run main.go
```

Vous devriez voir :
```
Serveur démarré sur http://localhost:8080
```

### 3. Accéder à l'application

**Interface Web :**
- **Page d'accueil** : [http://localhost:8080](http://localhost:8080)
- **Page de profil** : [http://localhost:8080/profile](http://localhost:8080/profile)
- **Test de profil** : [http://localhost:8080/profile-debug](http://localhost:8080/profile-debug)

**Interface d'authentification :**
- **Connexion** : [http://localhost:8080/auth/login](http://localhost:8080/auth/login)
- **Inscription** : [http://localhost:8080/auth/register](http://localhost:8080/auth/register)

## 🔑 Comptes de test

### Compte administrateur
- **Email** : admin@forum.com
- **Mot de passe** : Jesaispas01*

### Compte utilisateur
- **Email** : nans13@gmail.com
- **Mot de passe** : Eclipse1234@

## 📡 API Endpoints

### Routes publiques (Web)
- `GET /` - Page d'accueil
- `GET /profile` - Page de profil (nécessite connexion)
- `GET /threads` - Liste des discussions
- `GET /threads/show/{id}` - Affichage d'une discussion
- `GET /auth/login` - Page de connexion
- `GET /auth/register` - Page d'inscription

### Routes API REST

#### 🌍 Routes publiques
- `POST /api/register` - Inscription d'un nouvel utilisateur
- `POST /api/login` - Connexion d'un utilisateur
- `GET /api/threads` - Liste des fils de discussion
- `GET /api/threads/{id}` - Détails d'un fil de discussion
- `GET /api/threads/{id}/messages` - Messages d'un fil de discussion
- `GET /api/search` - Recherche de fils de discussion
- `GET /api/categories` - Liste des catégories
- `GET /api/categories/{id}` - Détails d'une catégorie

#### 🔒 Routes protégées (nécessite un token JWT)
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

#### 👑 Routes admin (nécessite un token JWT admin)
- `POST /api/admin/users/{id}/ban` - Bannir un utilisateur
- `POST /api/admin/users/{id}/unban` - Débannir un utilisateur
- `PUT /api/admin/threads/{id}/status` - Mise à jour du statut d'un fil
- `GET /api/admin/stats` - Statistiques du forum
- `POST /api/admin/categories` - Création d'une catégorie
- `PUT /api/admin/categories/{id}` - Mise à jour d'une catégorie
- `DELETE /api/admin/categories/{id}` - Suppression d'une catégorie

## 🔐 Authentification

L'application utilise **JWT (JSON Web Tokens)** pour l'authentification :

1. **Côté client** : Les tokens sont stockés dans `localStorage`
2. **Côté serveur** : Vérification via middleware pour les routes protégées

### Utilisation de l'API
```javascript
// Exemple d'appel API avec authentification
fetch('/api/users/me', {
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('jwt_token'),
    'Content-Type': 'application/json'
  }
})
```

## 🛠️ Développement

### Structure du code

**Contrôleurs** (`controllers/`) :
- Gèrent la logique métier
- Interagissent avec la base de données
- Retournent des réponses HTTP/JSON

**Middleware** (`middleware/`) :
- Vérification d'authentification
- Gestion des sessions
- Contrôle d'accès

**Routes** (`routes/`) :
- Configuration des endpoints
- Association routes ↔ contrôleurs
- Application des middlewares

**Frontend** (`static/`, `templates/`) :
- Interface utilisateur HTML/CSS/JS
- Communication avec l'API REST
- Gestion de l'état côté client

### Ajout de nouvelles fonctionnalités

1. **Nouveau contrôleur** → `controllers/nouveau_controller.go`
2. **Nouvelles routes** → `routes/nouvelles_routes.go`
3. **Nouveau template** → `templates/nouveau_template.html`
4. **Nouveaux styles** → `static/css/styles.css`

## 🐛 Débogage

### Logs utiles
```bash
# Démarrer avec logs détaillés
go run main.go
```

### Routes de debug
- `GET /test-profile` - Test d'accès au profil
- `GET /profile-debug` - Profil sans middleware d'authentification

### Problèmes courants

1. **Erreur 404 sur `/profile`** :
   - Vérifiez que vous êtes connecté
   - Testez avec `/profile-debug`

2. **Erreur de base de données** :
   - Vérifiez que MySQL est démarré dans XAMPP
   - Vérifiez les paramètres dans `.env`

3. **Token JWT invalide** :
   - Videz le localStorage du navigateur
   - Reconnectez-vous


## 🤝 Contribution

1. Fork le projet
2. Créez une branche pour votre fonctionnalité
3. Commitez vos changements
4. Poussez vers la branche
5. Ouvrez une Pull Request

---

🚀 **Bon développement !**