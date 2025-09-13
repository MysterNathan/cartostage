-- Création de la table stages
CREATE TABLE IF NOT EXISTS stages (
                                      id SERIAL PRIMARY KEY,
                                      poste VARCHAR(255) NOT NULL,
                                      adresse TEXT NOT NULL,
                                      lat DECIMAL(10, 8) NOT NULL,
                                      lng DECIMAL(11, 8) NOT NULL,
                                      places_disponibles INTEGER NOT NULL DEFAULT 0,
                                      enterprise VARCHAR(255) NOT NULL,
                                      filiere VARCHAR(100) NOT NULL,
                                      sector VARCHAR(100) NOT NULL,
                                      commune VARCHAR(255) NOT NULL,
                                      capacity_total INTEGER NOT NULL DEFAULT 0,
                                      capacity_filled INTEGER NOT NULL DEFAULT 0,
                                      period VARCHAR(100) NOT NULL,
                                      parcours VARCHAR(50) NOT NULL CHECK (parcours IN ('scolaire', 'apprentissage', 'mixte')),
                                      famille_metiers VARCHAR(255),
                                      niveau_scolaire VARCHAR(10) CHECK (niveau_scolaire IN ('2de', '1re', 'Tle')),
                                      created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
                                      updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Contraintes de validation
                                      CONSTRAINT chk_capacity_valid
                                          CHECK (capacity_filled <= capacity_total),
                                      CONSTRAINT chk_places_valid
                                          CHECK (places_disponibles >= 0),
                                      CONSTRAINT chk_coordinates_valid
                                          CHECK (lat BETWEEN -90 AND 90 AND lng BETWEEN -180 AND 180)
);

-- Index pour les recherches géospatiales et fréquentes
CREATE INDEX IF NOT EXISTS idx_stages_lat_lng ON stages(lat, lng);
CREATE INDEX IF NOT EXISTS idx_stages_commune ON stages(commune);
CREATE INDEX IF NOT EXISTS idx_stages_filiere ON stages(filiere);
CREATE INDEX IF NOT EXISTS idx_stages_sector ON stages(sector);
CREATE INDEX IF NOT EXISTS idx_stages_enterprise ON stages(enterprise);
CREATE INDEX IF NOT EXISTS idx_stages_period ON stages(period);
CREATE INDEX IF NOT EXISTS idx_stages_parcours ON stages(parcours);
CREATE INDEX IF NOT EXISTS idx_stages_niveau_scolaire ON stages(niveau_scolaire);
CREATE INDEX IF NOT EXISTS idx_stages_places_disponibles ON stages(places_disponibles);
CREATE INDEX IF NOT EXISTS idx_stages_created_at ON stages(created_at);

-- Index composite pour les recherches complexes
CREATE INDEX IF NOT EXISTS idx_stages_filiere_commune ON stages(filiere, commune);
CREATE INDEX IF NOT EXISTS idx_stages_sector_period ON stages(sector, period);
CREATE INDEX IF NOT EXISTS idx_stages_parcours_niveau ON stages(parcours, niveau_scolaire);

-- Trigger pour mettre à jour updated_at automatiquement
DROP TRIGGER IF EXISTS update_stages_updated_at ON stages;
CREATE TRIGGER update_stages_updated_at
    BEFORE UPDATE ON stages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
