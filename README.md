
# 🎬 CinéForum

Plateforme de discussion et de partage de son interet sur le cinema

## 📁 Architecture du projet

```
cineforum/
├── public/                     # Fichiers statiques servis
│   ├── index.html             # Page principale
│   └── static/
│       ├── css/              # Styles compilés
│       └── js/               # Scripts compilés
├── src/                       # Code source
│   ├── css/                  # Styles SCSS modulaires
│   ├── js/                   # JavaScript modulaire
│   └── html/                 # Templates HTML
├── backend/                   # Backend (à développer)
├── tests/                    # Tests automatisés
└── docs/                     # Documentation
```

## 🚀 Installation et démarrage

### Prérequis
- Node.js (v16+)
- npm ou yarn

### Installation
```bash
# Cloner le projet
git clone https://github.com/votre-repo/cineforum.git
cd cineforum

# Installer les dépendances
npm install
```

### Développement
```bash
# Serveur de développement
npm run dev

# Compilation en mode watch
npm run watch

# Linting et formatage
npm run lint
npm run prettier
```

### Production
```bash
npm run build


npm run serve
```

## 📂 Structure des fichiers

### HTML (public/index.html)
Page principale nettoyée sans JavaScript intégré. Les scripts sont chargés séparément pour une meilleure organisation.

### JavaScript
- **main.js** : Point d'entrée, initialisation de l'app
- **auth.js** : Gestion de l'authentification (JWT, localStorage)
- **threads.js** : Gestion des discussions (CRUD, affichage)
- **utils.js** : Fonctions utilitaires (formatage, validation)

### CSS (à créer)
Structure modulaire recommandée :
- **base/** : Reset, variables, typographie
- **components/** : Boutons, cartes, formulaires
- **layouts/** : Header, footer, grilles

## 🔧 Fonctionnalités

### ✅ Implémentées
- ✅ Interface utilisateur complète
- ✅ Gestion de l'authentification (frontend)
- ✅ Affichage des discussions
- ✅ Recherche et filtres
- ✅ Modales de création
- ✅ Données d'exemple

### 🚧 À développer
- 🚧 Backend API (Node.js/Express ou Python/Django)
- 🚧 Base de données (PostgreSQL/MongoDB)
- 🚧 Authentification réelle
- 🚧 Upload d'images
- 🚧 Notifications en temps réel
- 🚧 API externe films (TMDB)

## 🎯 Prochaines étapes

### Phase 1: Backend
1. **API REST** avec authentification JWT
2. **Base de données** pour utilisateurs et discussions
3. **Upload d'images** pour les affiches
4. **Tests unitaires** et d'intégration

### Phase 2: Fonctionnalités avancées
1. **Notifications** en temps réel (WebSocket)
2. **Système de vote** et modération
3. **Intégration API TMDB** pour les métadonnées films
4. **Recommandations** personnalisées

### Phase 3: Optimisations
1. **PWA** (Progressive Web App)
2. **Cache** et performance
3. **SEO** et accessibilité
4. **Analytics** et monitoring

## 🛠 Technologies utilisées

### Frontend
- **HTML5** sémantique
- **CSS3** / **SASS** pour les styles
- **JavaScript ES6+** natif (pas de framework)
- **Webpack** pour le bundling

### Outils de développement
- **ESLint** pour la qualité du code
- **Prettier** pour le formatage
- **Jest** pour les tests
- **Live Server** pour le développement

### Backend (prévu)
- **Node.js** + **Express** ou **Python** + **Django**
- **PostgreSQL** ou **MongoDB**
- **JWT** pour l'authentification
- **Cloudinary** pour les images

## 📝 Scripts disponibles

| Script | Description |
|--------|-------------|
| `npm run dev` | Serveur de développement |
| `npm run build` | Build pour production |
| `npm run watch` | Compilation en mode watch |
| `npm run lint` | Vérification du code |
| `npm run test` | Tests automatisés |
| `npm run prettier` | Formatage du code |

## 🔗 API Backend (à implémenter)

### Endpoints prévus
```
GET    /api/threads          # Liste des discussions
POST   /api/threads          # Créer une discussion
GET    /api/threads/:id      # Détail d'une discussion
PUT    /api/threads/:id      # Modifier une discussion
DELETE /api/threads/:id      # Supprimer une discussion

POST   /api/auth/login       # Connexion
POST   /api/auth/register    # Inscription
POST   /api/auth/logout      # Déconnexion
GET    /api/auth/profile     # Profil utilisateur

GET    /api/movies/search    # Recherche de films (TMDB)
```

## 🎨 Design System

### Couleurs (à définir dans variables CSS)
```css
:root {
  --primary: #e50914;     /* Rouge Netflix */
  --secondary: #f5c518;   /* Jaune IMDB */
  --dark: #141414;        /* Noir Netflix */
  --light: #ffffff;
  --gray: #8c8c8c;
}
```

### Typographie
- **Titres** : Font moderne (Inter, Roboto)
- **Corps** : Font lisible
- **Code** : Monospace (Fira Code)

## 📱 Responsive Design

Le design s'adapte automatiquement :
- **Mobile** : < 768px
- **Tablet** : 768px - 1024px  
- **Desktop** : > 1024px

## 🤝 Contribution

1. **Fork** le projet
2. **Créer** une branche feature (`git checkout -b feature/AmazingFeature`)
3. **Commit** vos changements (`git commit -m 'Add AmazingFeature'`)
4. **Push** sur la branche (`git push origin feature/AmazingFeature`)
5. **Ouvrir** une Pull Request

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 👥 Équipe

- **Frontend** : Interface utilisateur moderne
- **Backend** : API robuste (à développer)
- **Design** : UX/UI cinéma-friendly
- **DevOps** : Déploiement et monitoring

---

**CinéForum** - *Partagez votre passion du cinéma* 🎬