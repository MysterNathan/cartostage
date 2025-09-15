-- Migration: 001_create_users_table.sql

CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) UNIQUE NOT NULL,
                                     email VARCHAR(255) UNIQUE,
                                     password_hash VARCHAR(255) NOT NULL,
                                     role VARCHAR(50) DEFAULT 'admin' CHECK (role IN ('student', 'teacher', 'enterprise', 'admin')),
                                     entity_id INTEGER, -- Référence à l'école/entreprise selon le rôle
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     last_login TIMESTAMP NULL
);

-- Index pour améliorer les performances
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_entity_id ON users(entity_id);
CREATE INDEX IF NOT EXISTS idx_users_last_login ON users(last_login);

-- Trigger pour mettre à jour updated_at automatiquement
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Insérer des utilisateurs par défaut
-- Admin principal (mot de passe: "password")
INSERT INTO users (username, email, password_hash, role) VALUES
    ('admin', 'admin@stages.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin')
ON CONFLICT (username) DO NOTHING;

-- Utilisateurs de test pour chaque rôle
INSERT INTO users (username, email, password_hash, role, entity_id) VALUES
                                                                        -- Étudiant (entity_id = école/établissement)
                                                                        ('etudiant1', 'etudiant1@ecole.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', 1),

                                                                        -- Professeur (entity_id = école/établissement)
                                                                        ('prof1', 'prof1@ecole.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'teacher', 1),

                                                                        -- Entreprise (entity_id = ID entreprise)
                                                                        ('entreprise1', 'contact@entreprise1.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'enterprise', 1)
ON CONFLICT (username) DO NOTHING;

-- Commentaires pour documenter la logique entity_id
COMMENT ON COLUMN users.entity_id IS 'ID de l''entité associée: école pour student/teacher, entreprise pour enterprise, NULL pour admin';
COMMENT ON COLUMN users.role IS 'Rôle utilisateur: student, teacher, enterprise, admin';
COMMENT ON TABLE users IS 'Table des utilisateurs avec gestion des rôles et entités associées';
