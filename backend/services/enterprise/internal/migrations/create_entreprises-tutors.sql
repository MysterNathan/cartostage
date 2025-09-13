-- migrations/001_create_enterprises_and_tutors.sql

-- Extension pour UUID si nécessaire
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table entreprises
CREATE TABLE enterprises (
                             id SERIAL PRIMARY KEY,
                             nom VARCHAR(255) NOT NULL,
                             adresse TEXT,
                             secteur VARCHAR(255),
                             taille VARCHAR(20) CHECK (taille IN ('TPE', 'PME', 'ETI', 'GE')),
                             siret VARCHAR(14) UNIQUE,
                             email_contact VARCHAR(255),
                             telephone VARCHAR(20),
                             site_web VARCHAR(255),
                             description TEXT,
                             logo_url VARCHAR(500),
                             is_active BOOLEAN DEFAULT true,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table tuteurs
CREATE TABLE tutors (
                        id SERIAL PRIMARY KEY,
                        enterprise_id INTEGER NOT NULL REFERENCES enterprises(id) ON DELETE CASCADE,
                        prenom VARCHAR(100) NOT NULL,
                        nom VARCHAR(100) NOT NULL,
                        email VARCHAR(255) UNIQUE NOT NULL,
                        telephone VARCHAR(20),
                        poste VARCHAR(255),
                        departement VARCHAR(255),
                        is_active BOOLEAN DEFAULT true,
                        is_primary BOOLEAN DEFAULT false,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Contrainte : un seul tuteur principal par entreprise
                        CONSTRAINT unique_primary_tutor_per_enterprise
                            EXCLUDE (enterprise_id WITH =) WHERE (is_primary = true)
);

-- Mise à jour de la table stages
ALTER TABLE stages
    ADD COLUMN enterprise_id INTEGER REFERENCES enterprises(id),
ADD COLUMN tutor_id INTEGER REFERENCES tutors(id);

-- Mise à jour de la table users pour les nouveaux rôles
ALTER TABLE users
ALTER COLUMN role TYPE VARCHAR(50),
ADD CONSTRAINT check_role CHECK (role IN ('student', 'teacher', 'tutor', 'enterprise_admin', 'admin'));

-- Index pour optimiser les performances
CREATE INDEX idx_tutors_enterprise_id ON tutors(enterprise_id);
CREATE INDEX idx_tutors_email ON tutors(email);
CREATE INDEX idx_tutors_active ON tutors(is_active);
CREATE INDEX idx_enterprises_siret ON enterprises(siret);
CREATE INDEX idx_enterprises_secteur ON enterprises(secteur);
CREATE INDEX idx_stages_enterprise_tutor ON stages(enterprise_id, tutor_id);

-- Trigger pour mettre à jour updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_enterprises_updated_at BEFORE UPDATE ON enterprises FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tutors_updated_at BEFORE UPDATE ON tutors FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Données de test
INSERT INTO enterprises (nom, adresse, secteur, taille, siret, email_contact, telephone, description) VALUES
                                                                                                          ('TechCorp SARL', '123 Avenue des Technologies, 75001 Paris', 'Informatique', 'PME', '12345678901234', 'contact@techcorp.fr', '0123456789', 'Entreprise spécialisée en développement logiciel'),
                                                                                                          ('MecaSolutions', '45 Rue de l''Industrie, 69000 Lyon', 'Mécanique', 'ETI', '98765432109876', 'rh@mecasolutions.fr', '0478901234', 'Conception et fabrication de pièces mécaniques'),
                                                                                                          ('ElecPower', '78 Boulevard Électrique, 31000 Toulouse', 'Électronique', 'TPE', '11223344556677', 'jobs@elecpower.fr', '0561234567', 'Solutions en électronique industrielle');

INSERT INTO tutors (enterprise_id, prenom, nom, email, telephone, poste, departement, is_primary) VALUES
                                                                                                      (1, 'Marie', 'Dubois', 'marie.dubois@techcorp.fr', '0123456790', 'Responsable RH', 'Ressources Humaines', true),
                                                                                                      (1, 'Pierre', 'Martin', 'pierre.martin@techcorp.fr', '0123456791', 'Lead Developer', 'Développement', false),
                                                                                                      (2, 'Sophie', 'Bernard', 'sophie.bernard@mecasolutions.fr', '0478901235', 'Chef de projet', 'Production', true),
                                                                                                      (3, 'Jean', 'Moreau', 'jean.moreau@elecpower.fr', '0561234568', 'Directeur Technique', 'R&D', true);

INSERT INTO users (email, password_hash, role, entity_id) VALUES
                                                              ('marie.dubois@techcorp.fr', '$2a$10$example1', 'tutor', 1),
                                                              ('pierre.martin@techcorp.fr', '$2a$10$example2', 'tutor', 2),
                                                              ('sophie.bernard@mecasolutions.fr', '$2a$10$example3', 'tutor', 3),
                                                              ('jean.moreau@elecpower.fr', '$2a$10$example4', 'tutor', 4),
                                                              ('admin@techcorp.fr', '$2a$10$example5', 'enterprise_admin', 1),
                                                              ('admin@mecasolutions.fr', '$2a$10$example6', 'enterprise_admin', 2);
