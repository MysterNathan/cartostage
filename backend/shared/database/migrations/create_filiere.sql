-- Création de la fonction update_updated_at_column si elle n'existe pas déjà
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

-- Création de la table filieres
CREATE TABLE IF NOT EXISTS filieres (
                                        id SERIAL PRIMARY KEY,
                                        code VARCHAR(50) NOT NULL UNIQUE,
    label VARCHAR(255) NOT NULL,
    color VARCHAR(7) NOT NULL CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

    -- Contraintes de validation
    CONSTRAINT chk_code_not_empty CHECK (LENGTH(TRIM(code)) > 0),
    CONSTRAINT chk_label_not_empty CHECK (LENGTH(TRIM(label)) > 0)
    );

-- Index pour les recherches fréquentes
CREATE INDEX IF NOT EXISTS idx_filieres_code ON filieres(code);
CREATE INDEX IF NOT EXISTS idx_filieres_label ON filieres(label);
CREATE INDEX IF NOT EXISTS idx_filieres_created_at ON filieres(created_at);

-- Trigger pour mettre à jour updated_at automatiquement
DROP TRIGGER IF EXISTS update_filieres_updated_at ON filieres;
CREATE TRIGGER update_filieres_updated_at
    BEFORE UPDATE ON filieres
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Ajout d'une clé étrangère sur la table stages pour référencer filieres.code
-- (optionnel, selon si vous voulez une contrainte référentielle)
-- ALTER TABLE stages
-- ADD CONSTRAINT fk_stages_filiere
-- FOREIGN KEY (filiere) REFERENCES filieres(code) ON UPDATE CASCADE;
