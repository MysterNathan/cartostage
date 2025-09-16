-- Migration: 001_create_users_and_profiles_tables.sql

-- Table des utilisateurs principale
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     username VARCHAR(255) UNIQUE NOT NULL,
                                     email VARCHAR(255) UNIQUE,
                                     password_hash VARCHAR(255) NOT NULL,
                                     role VARCHAR(50) DEFAULT 'admin' CHECK (role IN ('student', 'teacher', 'enterprise', 'admin', 'tutor')),
                                     entity_type VARCHAR(50), -- Type d'entité: 'school', 'company', etc.
                                     entity_id INTEGER, -- ID de l'entité associée
                                     is_active BOOLEAN DEFAULT true,
                                     email_verified BOOLEAN DEFAULT false,
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     last_login TIMESTAMP NULL
);

-- Table des profils utilisateurs (informations détaillées)
CREATE TABLE IF NOT EXISTS user_profiles (
                                             user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                                             first_name VARCHAR(255),
                                             last_name VARCHAR(255),
                                             phone VARCHAR(20),
                                             poste VARCHAR(255),
                                             departement VARCHAR(255),
                                             is_active BOOLEAN DEFAULT true,
                                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index pour améliorer les performances
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_entity_type ON users(entity_type);
CREATE INDEX IF NOT EXISTS idx_users_entity_id ON users(entity_id);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_last_login ON users(last_login);
CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);

-- Trigger pour mettre à jour updated_at automatiquement
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers pour les deux tables
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_profiles_updated_at
    BEFORE UPDATE ON user_profiles
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Insérer des utilisateurs par défaut avec leurs profils
-- Admin principal (mot de passe: "password")
INSERT INTO users (username, email, password_hash, role, is_active, email_verified) VALUES
    ('admin', 'admin@stages.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', true, true)
ON CONFLICT (username) DO NOTHING;

-- Profil pour l'admin
INSERT INTO user_profiles (user_id, first_name, last_name, poste, departement)
SELECT id, 'Admin', 'Principal', 'Administrateur Système', 'IT'
FROM users WHERE username = 'admin'
ON CONFLICT (user_id) DO NOTHING;

-- Utilisateurs étudiants
INSERT INTO users (username, email, password_hash, role, entity_type, entity_id, is_active, email_verified) VALUES
                                                                                                                ('etudiant1', 'etudiant1@ecole.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', 'school', 1, true, true),
                                                                                                                ('etudiant2', 'etudiant2@ecole.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', 'school', 1, true, false),
                                                                                                                ('etudiant3', 'etudiant3@ecole.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'student', 'school', 2, true, true)
ON CONFLICT (username) DO NOTHING;

-- Profils pour les étudiants
INSERT INTO user_profiles (user_id, first_name, last_name, phone)
SELECT id, 'Jean', 'Dupont', '+33123456789' FROM users WHERE username = 'etudiant1'
UNION ALL
SELECT id, 'Marie', 'Martin', '+33123456790' FROM users WHERE username = 'etudiant2'
UNION ALL
SELECT id, 'Pierre', 'Durand', '+33123456791' FROM users WHERE username = 'etudiant3'
ON CONFLICT (user_id) DO NOTHING;

-- Utilisateurs professeurs
INSERT INTO users (username, email, password_hash, role, entity_type, entity_id, is_active, email_verified) VALUES
                                                                                                                ('prof1', 'prof1@ecole.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'teacher', 'school', 1, true, true),
                                                                                                                ('prof2', 'prof2@ecole.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'teacher', 'school', 1, true, true),
                                                                                                                ('prof3', 'prof3@ecole.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'teacher', 'school', 2, true, true)
ON CONFLICT (username) DO NOTHING;

-- Profils pour les professeurs
INSERT INTO user_profiles (user_id, first_name, last_name, phone, poste, departement)
SELECT id, 'Robert', 'Moreau', '+33123456792', 'Professeur de Mathématiques', 'Sciences' FROM users WHERE username = 'prof1'
UNION ALL
SELECT id, 'Sophie', 'Bernard', '+33123456793', 'Professeure d''Informatique', 'Technique' FROM users WHERE username = 'prof2'
UNION ALL
SELECT id, 'Michel', 'Petit', '+33123456794', 'Professeur de Physique', 'Sciences' FROM users WHERE username = 'prof3'
ON CONFLICT (user_id) DO NOTHING;

-- Utilisateurs entreprises
INSERT INTO users (username, email, password_hash, role, entity_type, entity_id, is_active, email_verified) VALUES
                                                                                                                ('entreprise1', 'contact@entreprise1.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'enterprise', 'company', 1, true, true),
                                                                                                                ('entreprise2', 'contact@entreprise2.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'enterprise', 'company', 2, true, true),
                                                                                                                ('entreprise3', 'contact@entreprise3.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'enterprise', 'company', 3, false, true)
ON CONFLICT (username) DO NOTHING;

-- Profils pour les entreprises
INSERT INTO user_profiles (user_id, first_name, last_name, phone, poste, departement)
SELECT id, 'Catherine', 'Dubois', '+33123456795', 'Responsable RH', 'Ressources Humaines' FROM users WHERE username = 'entreprise1'
UNION ALL
SELECT id, 'François', 'Leroy', '+33123456796', 'Directeur Technique', 'IT' FROM users WHERE username = 'entreprise2'
UNION ALL
SELECT id, 'Isabelle', 'Roux', '+33123456797', 'Chargée de Recrutement', 'RH' FROM users WHERE username = 'entreprise3'
ON CONFLICT (user_id) DO NOTHING;

-- Utilisateurs tuteurs (nouveaux)
INSERT INTO users (username, email, password_hash, role, entity_type, entity_id, is_active, email_verified) VALUES
                                                                                                                ('tutor1', 'tutor1@entreprise1.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'tutor', 'company', 1, true, true),
                                                                                                                ('tutor2', 'tutor2@entreprise1.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'tutor', 'company', 1, true, true),
                                                                                                                ('tutor3', 'tutor3@entreprise2.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'tutor', 'company', 2, true, true),
                                                                                                                ('tutor4', 'tutor4@entreprise2.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'tutor', 'company', 2, true, false),
                                                                                                                ('tutor5', 'tutor5@entreprise3.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'tutor', 'company', 3, true, true)
ON CONFLICT (username) DO NOTHING;

-- Profils pour les tuteurs
INSERT INTO user_profiles (user_id, first_name, last_name, phone, poste, departement)
SELECT id, 'Alain', 'Garcia', '+33123456798', 'Ingénieur Senior', 'Développement' FROM users WHERE username = 'tutor1'
UNION ALL
SELECT id, 'Sylvie', 'Thomas', '+33123456799', 'Chef de Projet', 'IT' FROM users WHERE username = 'tutor2'
UNION ALL
SELECT id, 'Laurent', 'Martinez', '+33123456800', 'Architecte Logiciel', 'R&D' FROM users WHERE username = 'tutor3'
UNION ALL
SELECT id, 'Nathalie', 'Lopez', '+33123456801', 'Responsable Qualité', 'QA' FROM users WHERE username = 'tutor4'
UNION ALL
SELECT id, 'Thierry', 'Gonzalez', '+33123456802', 'Consultant Senior', 'Conseil' FROM users WHERE username = 'tutor5'
ON CONFLICT (user_id) DO NOTHING;

-- Commentaires pour documenter la structure
COMMENT ON COLUMN users.entity_type IS 'Type d''entité associée: school pour student/teacher, company pour enterprise/tutor, NULL pour admin';
COMMENT ON COLUMN users.entity_id IS 'ID de l''entité associée selon le type (école, entreprise, etc.)';
COMMENT ON COLUMN users.role IS 'Rôle utilisateur: student, teacher, enterprise, admin, tutor';
COMMENT ON COLUMN users.is_active IS 'Statut actif/inactif de l''utilisateur';
COMMENT ON COLUMN users.email_verified IS 'Email vérifié ou non';

COMMENT ON TABLE users IS 'Table des utilisateurs avec gestion des rôles et entités associées';
COMMENT ON TABLE user_profiles IS 'Profils détaillés des utilisateurs avec informations personnelles et professionnelles';

-- Afficher un résumé des utilisateurs créés
SELECT
    'Utilisateurs créés:' as info,
    role,
    entity_type,
    COUNT(*) as nombre
FROM users
GROUP BY role, entity_type
ORDER BY role, entity_type;
