-- Création de la table tutors
CREATE TABLE IF NOT EXISTS tutors (
                                      id SERIAL PRIMARY KEY,
                                      enterprise_id INTEGER NOT NULL,
                                      prenom VARCHAR(100) NOT NULL,
                                      nom VARCHAR(100) NOT NULL,
                                      email VARCHAR(255) NOT NULL,
                                      telephone VARCHAR(20),
                                      poste VARCHAR(255),
                                      departement VARCHAR(255),
                                      is_active BOOLEAN DEFAULT true,
                                      created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                      updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Contrainte de clé étrangère
                                      CONSTRAINT fk_tutors_enterprise
                                          FOREIGN KEY (enterprise_id)
                                              REFERENCES enterprises(id)
                                              ON DELETE CASCADE
);

-- Index pour les recherches fréquentes
CREATE INDEX IF NOT EXISTS idx_tutors_enterprise_id ON tutors(enterprise_id);
CREATE INDEX IF NOT EXISTS idx_tutors_email ON tutors(email);
CREATE INDEX IF NOT EXISTS idx_tutors_nom_prenom ON tutors(nom, prenom);
CREATE INDEX IF NOT EXISTS idx_tutors_is_active ON tutors(is_active);
CREATE INDEX IF NOT EXISTS idx_tutors_created_at ON tutors(created_at);

-- Contrainte d'unicité pour éviter les doublons email par entreprise
CREATE UNIQUE INDEX IF NOT EXISTS idx_tutors_unique_email_enterprise
    ON tutors(email, enterprise_id);

-- Trigger pour mettre à jour updated_at automatiquement
DROP TRIGGER IF EXISTS update_tutors_updated_at ON tutors;
CREATE TRIGGER update_tutors_updated_at
    BEFORE UPDATE ON tutors
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
