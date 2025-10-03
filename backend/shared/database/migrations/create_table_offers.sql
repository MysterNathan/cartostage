-- Migration: 001_create_stage_offers.up.sql

CREATE TABLE stage_offers (
                              id SERIAL PRIMARY KEY,
                              position VARCHAR(255) NOT NULL,
                              address TEXT NOT NULL,
                              lat DECIMAL(10,8) NOT NULL,
                              lng DECIMAL(11,8) NOT NULL,
                              enterprise VARCHAR(255) NOT NULL,
                              sector VARCHAR(100) NOT NULL,
                              capacity_total INTEGER NOT NULL DEFAULT 1,
                              capacity_filled INTEGER NOT NULL DEFAULT 0,
                              period VARCHAR(100) NOT NULL,
                              course VARCHAR(255) NOT NULL,
                              job_family VARCHAR(100) NOT NULL,
                              scolar_level VARCHAR(50) NOT NULL,
                              created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

                              CONSTRAINT check_capacity CHECK (capacity_filled <= capacity_total),
                              CONSTRAINT check_coordinates_lat CHECK (lat >= -90 AND lat <= 90),
                              CONSTRAINT check_coordinates_lng CHECK (lng >= -180 AND lng <= 180)
);

-- Fonction pour mettre à jour automatiquement updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger pour updated_at
CREATE TRIGGER update_stage_offers_updated_at
    BEFORE UPDATE ON stage_offers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Données d'exemple
INSERT INTO stage_offers (
    position, address, lat, lng, enterprise, sector,
    capacity_total, capacity_filled, period, course, job_family, scolar_level
) VALUES
      (
          'Développeur Full Stack',
          '42 Rue de la Tech, 75001 Paris',
          48.8566140,
          2.3522219,
          'TechCorp Solutions',
          'Informatique',
          3,
          3,
          '6 mois - Février à Juillet 2024',
          'Master Informatique',
          'Développement Web',
          'BAC+5'
      ),
      (
          'Assistant Marketing Digital',
          '15 Avenue des Champs-Élysées, 75008 Paris',
          48.8698679,
          2.3075756,
          'Digital Marketing Pro',
          'Marketing',
          2,
          2,
          '4 mois - Mars à Juin 2024',
          'BTS Communication',
          'Marketing Digital',
          'BAC+2'
      ),
      (
          'Développeur Mobile',
          '28 Rue du Commerce, 69002 Lyon',
          45.7578137,
          4.8320114,
          'MobileTech Industries',
          'Informatique',
          2,
          2,
          '5 mois - Janvier à Mai 2024',
          'Licence Professionnelle',
          'Développement Mobile',
          'BAC+3'
      ),
      (
          'Assistant RH',
          '10 Boulevard Haussmann, 75009 Paris',
          48.8738857,
          2.3324122,
          'HR Consulting Group',
          'Ressources Humaines',
          1,
          1,
          '3 mois - Avril à Juin 2024',
          'Master RH',
          'Gestion RH',
          'BAC+5'
      ),
      (
          'Designer UX/UI',
          '55 Rue de la République, 13001 Marseille',
          43.2961743,
          5.3699525,
          'Creative Studio',
          'Design',
          2,
          1,
          '4 mois - Mai à Août 2024',
          'Bachelor Design',
          'Design Numérique',
          'BAC+3'
      ),
      (
          'Data Analyst',
          '88 Avenue Carnot, 33000 Bordeaux',
          44.8404400,
          -0.5805000,
          'DataViz Solutions',
          'Informatique',
          1,
          1,
          '6 mois - Février à Juillet 2024',
          'Master Data Science',
          'Analyse de Données',
          'BAC+5'
      ),
      (
          'Community Manager',
          '23 Rue Victor Hugo, 31000 Toulouse',
          43.6044622,
          1.4442469,
          'Social Media Agency',
          'Communication',
          3,
          2,
          '3 mois - Juin à Août 2024',
          'Licence Communication',
          'Community Management',
          'BAC+3'
      );
