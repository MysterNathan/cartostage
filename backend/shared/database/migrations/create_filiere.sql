-- Migration: 006_create_filieres.up.sql

CREATE TABLE filieres (
                          id SERIAL PRIMARY KEY,
                          code VARCHAR(50) NOT NULL UNIQUE,
                          label VARCHAR(255) NOT NULL,
                          color VARCHAR(7) NOT NULL CHECK (color ~ '^#[0-9A-Fa-f]{6}$'),
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index pour performances
CREATE INDEX idx_filieres_code ON filieres(code);

-- Trigger pour updated_at
CREATE TRIGGER update_filieres_updated_at
    BEFORE UPDATE ON filieres
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Données d'exemple
INSERT INTO filieres (code, label, color) VALUES
                                              ('INFO', 'Informatique', '#3498db'),
                                              ('ELEC', 'Électronique', '#e74c3c'),
                                              ('MECA', 'Mécanique', '#f39c12'),
                                              ('COMM', 'Communication', '#9b59b6'),
                                              ('MARK', 'Marketing', '#e67e22'),
                                              ('RH', 'Ressources Humaines', '#2ecc71'),
                                              ('COMP', 'Comptabilité', '#34495e'),
                                              ('LOG', 'Logistique', '#16a085');
