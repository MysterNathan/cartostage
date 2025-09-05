-- Migration: create_stages_table.sql
-- Description: Création de la table stages pour le système de gestion des stages du lycée

CREATE TABLE IF NOT EXISTS stages (
                                      id SERIAL PRIMARY KEY,
                                      poste VARCHAR(255) NOT NULL,
    adresse TEXT NOT NULL,
    lat DECIMAL(10, 8) NOT NULL,
    lng DECIMAL(11, 8) NOT NULL,
    places_disponibles INTEGER DEFAULT 0 CHECK (places_disponibles >= 0),
    entreprise VARCHAR(255) NOT NULL,
    filiere VARCHAR(100) NOT NULL,
    sector VARCHAR(100),
    commune VARCHAR(100),
    capacity_total INTEGER NOT NULL DEFAULT 0 CHECK (capacity_total >= 0),
    capacity_filled INTEGER NOT NULL DEFAULT 0 CHECK (capacity_filled >= 0),
    period VARCHAR(100),
    parcours VARCHAR(20) NOT NULL CHECK (parcours IN ('scolaire', 'apprentissage', 'mixte')),
    famille_metiers VARCHAR(150),
    niveau_scolaire VARCHAR(10) NOT NULL CHECK (niveau_scolaire IN ('2de', '1re', 'Tle')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
                             );

-- Index pour améliorer les performances des requêtes fréquentes
CREATE INDEX IF NOT EXISTS idx_stages_filiere ON stages(filiere);
CREATE INDEX IF NOT EXISTS idx_stages_commune ON stages(commune);
CREATE INDEX IF NOT EXISTS idx_stages_sector ON stages(sector);
CREATE INDEX IF NOT EXISTS idx_stages_parcours ON stages(parcours);
CREATE INDEX IF NOT EXISTS idx_stages_famille_metiers ON stages(famille_metiers);
CREATE INDEX IF NOT EXISTS idx_stages_niveau_scolaire ON stages(niveau_scolaire);
CREATE INDEX IF NOT EXISTS idx_stages_lat_lng ON stages(lat, lng);

-- Contrainte pour s'assurer que capacity_filled <= capacity_total
ALTER TABLE stages ADD CONSTRAINT chk_capacity_consistency
    CHECK (capacity_filled <= capacity_total);

-- Trigger pour mettre à jour automatiquement updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_stages_updated_at
    BEFORE UPDATE ON stages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Commentaires sur la table et les colonnes
COMMENT ON TABLE stages IS 'Table des stages proposés par les entreprises partenaires du lycée';
COMMENT ON COLUMN stages.id IS 'Identifiant unique du stage';
COMMENT ON COLUMN stages.poste IS 'Intitulé du poste de stage';
COMMENT ON COLUMN stages.adresse IS 'Adresse complète de l''entreprise';
COMMENT ON COLUMN stages.lat IS 'Latitude GPS de l''entreprise';
COMMENT ON COLUMN stages.lng IS 'Longitude GPS de l''entreprise';
COMMENT ON COLUMN stages.places_disponibles IS 'Nombre de places disponibles actuellement';
COMMENT ON COLUMN stages.entreprise IS 'Nom de l''entreprise';
COMMENT ON COLUMN stages.filiere IS 'Filière de formation concernée';
COMMENT ON COLUMN stages.sector IS 'Secteur d''activité de l''entreprise';
COMMENT ON COLUMN stages.commune IS 'Commune de l''entreprise';
COMMENT ON COLUMN stages.capacity_total IS 'Nombre total de places pour ce stage';
COMMENT ON COLUMN stages.capacity_filled IS 'Nombre de places déjà pourvues';
COMMENT ON COLUMN stages.period IS 'Période du stage';
COMMENT ON COLUMN stages.parcours IS 'Type de parcours : scolaire, apprentissage ou mixte';
COMMENT ON COLUMN stages.famille_metiers IS 'Famille de métiers du stage';
COMMENT ON COLUMN stages.niveau_scolaire IS 'Niveau scolaire concerné : 2de, 1re ou Tle';
