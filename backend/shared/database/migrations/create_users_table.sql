-- Migration: 20240101000001_create_users_table.up.sql

-- Extension pour générer des UUIDs si nécessaire
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Création du type ENUM pour les rôles
DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('administrateur', 'enseignant', 'tuteur', 'eleve');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Création de la table users
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(50) UNIQUE NOT NULL,
                                     first_name VARCHAR(100) NOT NULL,
                                     last_name VARCHAR(100) NOT NULL,
                                     email VARCHAR(255) UNIQUE NOT NULL,
                                     password_hash TEXT NOT NULL,
                                     role user_role NOT NULL DEFAULT 'eleve',
                                     phone VARCHAR(15),
                                     is_active BOOLEAN NOT NULL DEFAULT true,
                                     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                     last_login TIMESTAMP WITH TIME ZONE
);

-- Index pour améliorer les performances
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);

-- Fonction pour mettre à jour automatiquement updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger pour mettre à jour updated_at automatiquement
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Insertion des utilisateurs par défaut
-- Mot de passe haché pour "admin" (utilise bcrypt avec cost 12)
-- Hash généré avec : golang.org/x/crypto/bcrypt.GenerateFromPassword([]byte("admin"), 12)
INSERT INTO users (username, first_name, last_name, email, password_hash, role, is_active)
VALUES
    (
        'admin',
        'Admin',
        'System',
        'admin@example.com',
        -- Hash bcrypt pour "admin" (cost 12)
        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/kYLs5JvXu',
        'administrateur',
        true
    ),
    (
        'user',
        'User',
        'Test',
        'user@example.com',
        -- Hash bcrypt pour "admin" (cost 12)
        '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj/kYLs5JvXu',
        'eleve',
        true
    )
ON CONFLICT (username) DO NOTHING;

-- Commentaires pour documenter la table
COMMENT ON TABLE users IS 'Table des utilisateurs du système';
COMMENT ON COLUMN users.password_hash IS 'Hash bcrypt du mot de passe';
COMMENT ON COLUMN users.role IS 'Rôle de l''utilisateur dans le système';
COMMENT ON COLUMN users.is_active IS 'Statut actif/inactif de l''utilisateur';
