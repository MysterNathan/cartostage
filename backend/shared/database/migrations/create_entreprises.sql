-- Création de la fonction trigger (une seule fois)
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Création de la table enterprises
CREATE TABLE IF NOT EXISTS enterprises (
                                           id SERIAL PRIMARY KEY,
                                           nom VARCHAR(255) NOT NULL,
                                           adresse TEXT NOT NULL,
                                           secteur VARCHAR(255) NOT NULL,
                                           taille VARCHAR(50),
                                           siret VARCHAR(14) UNIQUE,
                                           email_contact VARCHAR(255),
                                           telephone VARCHAR(20),
                                           site_web VARCHAR(255),
                                           description TEXT,
                                           logo_url VARCHAR(500),
                                           created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Index pour les recherches fréquentes
CREATE INDEX IF NOT EXISTS idx_enterprises_nom ON enterprises(nom);
CREATE INDEX IF NOT EXISTS idx_enterprises_secteur ON enterprises(secteur);
CREATE INDEX IF NOT EXISTS idx_enterprises_siret ON enterprises(siret);
CREATE INDEX IF NOT EXISTS idx_enterprises_created_at ON enterprises(created_at);

-- Trigger pour mettre à jour updated_at automatiquement
DROP TRIGGER IF EXISTS update_enterprises_updated_at ON enterprises;
CREATE TRIGGER update_enterprises_updated_at
    BEFORE UPDATE ON enterprises
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
