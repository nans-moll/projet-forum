-- Création de la base de données
CREATE DATABASE forum_db;

-- Connexion à la base de données
\c forum_db;

-- Création de la table users
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    is_banned BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL,
    last_connection TIMESTAMP NOT NULL,
    profile_picture VARCHAR(255),
    biography TEXT,
    message_count INTEGER NOT NULL DEFAULT 0,
    thread_count INTEGER NOT NULL DEFAULT 0
);

-- Création de la table threads
CREATE TABLE threads (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    description TEXT NOT NULL,
    tags TEXT[] NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'open',
    visibility VARCHAR(20) NOT NULL DEFAULT 'public',
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    message_count INTEGER NOT NULL DEFAULT 0
);

-- Création de la table messages
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    thread_id INTEGER NOT NULL REFERENCES threads(id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    image_url VARCHAR(255),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    likes INTEGER NOT NULL DEFAULT 0,
    dislikes INTEGER NOT NULL DEFAULT 0
);

-- Création de la table friendships
CREATE TABLE friendships (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    UNIQUE(user_id, friend_id)
);

-- Création de la table message_reactions
CREATE TABLE message_reactions (
    id SERIAL PRIMARY KEY,
    message_id INTEGER NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    reaction_type VARCHAR(10) NOT NULL CHECK (reaction_type IN ('like', 'dislike')),
    created_at TIMESTAMP NOT NULL,
    UNIQUE(message_id, user_id)
);

-- Création des index
CREATE INDEX idx_threads_author ON threads(author_id);
CREATE INDEX idx_threads_status ON threads(status);
CREATE INDEX idx_threads_visibility ON threads(visibility);
CREATE INDEX idx_messages_thread ON messages(thread_id);
CREATE INDEX idx_messages_author ON messages(author_id);
CREATE INDEX idx_friendships_users ON friendships(user_id, friend_id);
CREATE INDEX idx_message_reactions_message ON message_reactions(message_id);
CREATE INDEX idx_message_reactions_user ON message_reactions(user_id);

-- Création d'un utilisateur admin par défaut
INSERT INTO users (username, email, password_hash, role, created_at, last_connection)
VALUES (
    'admin',
    'admin@forum.com',
    -- Mot de passe: Admin123!@#
    '8c6976e5b5410415bde908bd4dee15dfb167a9c873fc4bb8a81f6f2ab448a918',
    'admin',
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
); 