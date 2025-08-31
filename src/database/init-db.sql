-- Création de la base de données si elle n'existe pas
-- (optionnel car elle est créée par les variables d'environnement Docker)

-- Extension pour les UUID (optionnel, au cas où tu voudrais utiliser des UUID plus tard)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Suppression de la table si elle existe (pour réinitialisation)
DROP TABLE IF EXISTS stages CASCADE;

-- Création de la table stages
CREATE TABLE stages (
    id SERIAL PRIMARY KEY,
    poste VARCHAR(255) NOT NULL,
    adresse TEXT NOT NULL,
    lat DECIMAL(10, 6) NOT NULL,
    lng DECIMAL(10, 6) NOT NULL,
    places_disponibles INTEGER NOT NULL DEFAULT 0,
    entreprise VARCHAR(255) NOT NULL,
    filiere VARCHAR(50) NOT NULL,
    sector VARCHAR(255) NOT NULL,
    commune VARCHAR(100) NOT NULL,
    capacity_total INTEGER NOT NULL DEFAULT 0,
    capacity_filled INTEGER NOT NULL DEFAULT 0,
    period VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index sur les colonnes souvent utilisées pour les recherches
CREATE INDEX idx_stages_filiere ON stages(filiere);
CREATE INDEX idx_stages_commune ON stages(commune);
CREATE INDEX idx_stages_entreprise ON stages(entreprise);
CREATE INDEX idx_stages_sector ON stages(sector);
CREATE INDEX idx_stages_places_disponibles ON stages(places_disponibles);

-- Fonction pour mettre à jour automatiquement updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger pour mettre à jour automatiquement updated_at
CREATE TRIGGER update_stages_updated_at 
    BEFORE UPDATE ON stages 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Insertion des données initiales
INSERT INTO stages (
    id, poste, adresse, lat, lng, places_disponibles, 
    entreprise, filiere, sector, commune, capacity_total, capacity_filled, period
) VALUES 
(1, 'Distributeur technique (élec/CVC/quincaillerie)', 'Le Tampon (adresse à préciser)', -21.288, 55.522, 2, 'TechDis Réunion Sud', 'CCST', 'Distributeur technique (élec/CVC/quincaillerie)', 'Le Tampon', 3, 1, 'Oct–Nov; Mar–Avr'),
(2, 'Intégrateur solaire/PV - ADV/gestion offres', 'Saint-Pierre (adresse à préciser)', -21.338, 55.485, 1, 'SOLARPRO Océan Indien', 'CCST', 'Intégrateur solaire/PV - ADV/gestion offres', 'Saint-Pierre', 2, 1, 'Oct–Nov; Mar–Avr'),
(3, 'Réseaux & Télécoms - B2B/CRM/SAV', 'Saint-Louis (adresse à préciser)', -21.279, 55.415, 1, 'TelcoPro Solutions', 'CCST', 'Réseaux & Télécoms - B2B/CRM/SAV', 'Saint-Louis', 3, 2, 'Oct–Nov; Mar–Avr'),
(4, 'Quincaillerie professionnelle - B2B', 'L''Étang-Salé (adresse à préciser)', -21.25, 55.35, 2, 'RéuniQuinca Pro', 'CCST', 'Quincaillerie professionnelle - B2B', 'L''Étang-Salé', 2, 0, 'Oct–Nov; Mar–Avr'),
(5, 'Outillage & EPI - Showroom pro', 'Les Avirons (adresse à préciser)', -21.203, 55.336, 1, 'Outillage Industrie Sud', 'CCST', 'Outillage & EPI - Showroom pro', 'Les Avirons', 2, 1, 'Oct–Nov; Mar–Avr'),
(6, 'Matériaux/consommables BTP - ADV', 'Saint-Pierre (adresse à préciser)', -21.345, 55.47, 1, 'BTP Grossiste Sud', 'CCST', 'Matériaux/consommables BTP - ADV', 'Saint-Pierre', 3, 2, 'Oct–Nov; Mar–Avr'),
(7, 'Fournitures industrielles - ADV/SAV', 'Le Tampon (adresse à préciser)', -21.276, 55.51, 2, 'Fournitures Techniques Océan', 'CCST', 'Fournitures industrielles - ADV/SAV', 'Le Tampon', 2, 0, 'Oct–Nov; Mar–Avr'),
(8, 'Sûreté/vidéo/IP - support/CRM', 'Saint-Louis (adresse à préciser)', -21.27, 55.402, 1, 'DataSec Intégration', 'CCST', 'Sûreté/vidéo/IP - support/CRM', 'Saint-Louis', 2, 1, 'Oct–Nov; Mar–Avr'),
(9, 'CVC/Climatisation - SAV/planning', 'L''Étang-Salé (adresse à préciser)', -21.262, 55.338, 2, 'Clim''Expert Réunion', 'CCST', 'CVC/Climatisation - SAV/planning', 'L''Étang-Salé', 3, 1, 'Oct–Nov; Mar–Avr'),
(10, 'Location matériel événementiel - devis/retours', 'Saint-Pierre (adresse à préciser)', -21.335, 55.49, 1, 'MobilEvent Location', 'CCST', 'Location matériel événementiel - devis/retours', 'Saint-Pierre', 2, 1, 'Oct–Nov; Mar–Avr'),
(11, 'EHPAD - hôtellerie, vie sociale', 'Le Tampon (adresse à préciser)', -21.289, 55.509, 1, 'EHPAD Les Flamboyants', 'ASSP', 'EHPAD - hôtellerie, vie sociale', 'Le Tampon', 3, 2, 'Oct–Nov; Mar–Avr'),
(12, 'Résidence autonomie - animation/logistique', 'Saint-Pierre (adresse à préciser)', -21.349, 55.482, 1, 'Résidence Autonomie Océane', 'ASSP', 'Résidence autonomie - animation/logistique', 'Saint-Pierre', 2, 1, 'Oct–Nov; Mar–Avr'),
(13, 'Service soins/infirmiers à domicile - planning', 'Saint-Joseph (adresse à préciser)', -21.372, 55.624, 2, 'SSIAD Sud Bien-Être', 'ASSP', 'Service soins/infirmiers à domicile - planning', 'Saint-Joseph', 2, 0, 'Oct–Nov; Mar–Avr'),
(14, 'Service aide à domicile - accompagnement', 'Petite-Île (adresse à préciser)', -21.34, 55.558, 1, 'SAAD Aide & Sourire', 'ASSP', 'Service aide à domicile - accompagnement', 'Petite-Île', 2, 1, 'Oct–Nov; Mar–Avr'),
(15, 'Maison d''accueil spécialisée - activités', 'Saint-Louis (adresse à préciser)', -21.268, 55.418, 1, 'MAS Les Horizons', 'ASSP', 'Maison d''accueil spécialisée - activités', 'Saint-Louis', 3, 2, 'Oct–Nov; Mar–Avr'),
(16, 'Institut médico-éducatif - ateliers encadrés', 'Les Avirons (adresse à préciser)', -21.215, 55.349, 1, 'IME Les Colibris', 'ASSP', 'Institut médico-éducatif - ateliers encadrés', 'Les Avirons', 2, 1, 'Oct–Nov; Mar–Avr'),
(17, 'Multi-accueil petite enfance - routines', 'Le Tampon (adresse à préciser)', -21.292, 55.523, 2, 'Crèche P''tits Dodos', 'ASSP', 'Multi-accueil petite enfance - routines', 'Le Tampon', 2, 0, 'Oct–Nov; Mar–Avr'),
(18, 'Clinique - hôtellerie/services logistiques', 'Saint-Pierre (adresse à préciser)', -21.332, 55.468, 1, 'Clinique du Littoral', 'ASSP', 'Clinique - hôtellerie/services logistiques', 'Saint-Pierre', 3, 2, 'Oct–Nov; Mar–Avr'),
(19, 'CCAS - action sociale / accueil', 'Le Tampon (adresse à préciser)', -21.277, 55.532, 1, 'CCAS du Tampon', 'ASSP', 'CCAS - action sociale / accueil', 'Le Tampon', 1, 0, 'Oct–Nov; Mar–Avr'),
(20, 'Centre hospitalier Sud Réunion - logistique/hôtellerie', 'Saint-Pierre (adresse à préciser)', -21.345, 55.46, 1, 'CHSR - Logistique', 'ASSP', 'Centre hospitalier Sud Réunion - logistique/hôtellerie', 'Saint-Pierre', 3, 2, 'Oct–Nov; Mar–Avr');

-- Réinitialiser la séquence pour que les prochains ID commencent à 21
SELECT setval('stages_id_seq', 20);

-- Création de vues utiles
CREATE OR REPLACE VIEW stages_summary AS
SELECT 
    filiere,
    COUNT(*) as total_stages,
    SUM(capacity_total) as total_capacity,
    SUM(capacity_filled) as total_filled,
    SUM(places_disponibles) as total_available
FROM stages 
GROUP BY filiere;

CREATE OR REPLACE VIEW stages_by_commune AS
SELECT 
    commune,
    COUNT(*) as total_stages,
    SUM(places_disponibles) as total_available
FROM stages 
GROUP BY commune
ORDER BY total_stages DESC;

-- Requêtes utiles pour vérifier les données
-- SELECT * FROM stages_summary;
-- SELECT * FROM stages_by_commune;
-- SELECT COUNT(*) as total_stages FROM stages;
-- SELECT filiere, COUNT(*) as count FROM stages GROUP BY filiere;

COMMIT;
