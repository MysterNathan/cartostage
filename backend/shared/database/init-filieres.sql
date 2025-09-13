-- Création de la table filieres
CREATE TABLE IF NOT EXISTS filieres (
                                        id SERIAL PRIMARY KEY,
                                        code VARCHAR(10) NOT NULL UNIQUE,
    label VARCHAR(200) NOT NULL,
    color VARCHAR(7) NOT NULL DEFAULT '#3B82F6',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- Création d'un index sur le code pour des recherches rapides
CREATE INDEX IF NOT EXISTS idx_filieres_code ON filieres(code);

-- Fonction pour mettre à jour automatiquement updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger pour mettre à jour updated_at automatiquement
DROP TRIGGER IF EXISTS update_filieres_updated_at ON filieres;
CREATE TRIGGER update_filieres_updated_at
    BEFORE UPDATE ON filieres
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insertion des données initiales
INSERT INTO filieres (id, code, label, color) VALUES
                                                  (1, 'CCST', 'CCST - Cybersécurité, Informatique et réseaux, Électronique', '#3B82F6'),
                                                  (2, 'ASSP', 'ASSP - Accompagnement, Soins et Services à la Personne', '#10B981'),
                                                  (3, 'SN', 'SN - Systèmes Numériques', '#06B6D4'),
                                                  (4, 'SLAM', 'SLAM - Solutions Logicielles et Applications Métiers', '#8B5CF6'),
                                                  (5, 'TEST', 'test', '#a10c0c')
    ON CONFLICT (id) DO UPDATE SET
    code = EXCLUDED.code,
                            label = EXCLUDED.label,
                            color = EXCLUDED.color,
                            updated_at = CURRENT_TIMESTAMP;

-- Réinitialiser la séquence pour les prochains ID
SELECT setval('filieres_id_seq', (SELECT MAX(id) FROM filieres));

-- Vérification des données insérées
SELECT * FROM filieres ORDER BY id;
