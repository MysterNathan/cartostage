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
          2,
          1,
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
          3,
          0,
          '4 mois - Mars à Juin 2024',
          'BTS Communication',
          'Marketing Digital',
          'BAC+2'
      );
