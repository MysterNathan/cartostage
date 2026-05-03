-- Migration: 002_create_establishments.up.sql

CREATE TABLE establishments (
                                id SERIAL PRIMARY KEY,
                                name VARCHAR(255) NOT NULL,
                                address TEXT NOT NULL,
                                lat DECIMAL(10,8) NOT NULL,
                                lng DECIMAL(11,8) NOT NULL,
                                sector VARCHAR(100) NOT NULL,
                                size INTEGER CHECK (size > 0),
                                siret VARCHAR(14) UNIQUE,
                                email VARCHAR(320) CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
                                phone VARCHAR(20),
                                website VARCHAR(500) CHECK (website ~* '^https?://'),
                                description TEXT,
                                logo_url VARCHAR(500) CHECK (logo_url ~* '^https?://'),
                                created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

                                CONSTRAINT check_coordinates_lat CHECK (lat >= -90 AND lat <= 90),
                                CONSTRAINT check_coordinates_lng CHECK (lng >= -180 AND lng <= 180),
                                CONSTRAINT check_siret_format CHECK (siret ~* '^[0-9]{14}$')
);

-- Index pour améliorer les performances
CREATE INDEX idx_establishments_sector ON establishments(sector);
CREATE INDEX idx_establishments_size ON establishments(size);
CREATE INDEX idx_establishments_coordinates ON establishments(lat, lng);

-- Trigger pour updated_at (réutilise la fonction créée précédemment)
CREATE TRIGGER update_establishments_updated_at
    BEFORE UPDATE ON establishments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Données d'exemple
INSERT INTO establishments (
    name, address, lat, lng, sector, size, siret, email, phone,
    website, description, logo_url
) VALUES
      (
          'TechCorp Solutions',
          '42 Rue de la Tech, 75001 Paris',
          48.8566140,
          2.3522219,
          'Informatique',
          75,
          '12345678901234',
          'contact@techcorp-solutions.fr',
          '+33 1 42 00 00 00',
          'https://www.techcorp-solutions.fr',
          'Entreprise spécialisée dans le développement de solutions logicielles sur mesure pour les PME et grandes entreprises.',
          'https://www.techcorp-solutions.fr/logo.png'
      ),
      (
          'Digital Marketing Pro',
          '15 Avenue des Champs-Élysées, 75008 Paris',
          48.8698679,
          2.3075756,
          'Marketing',
          25,
          '98765432109876',
          'info@digital-marketing-pro.com',
          '+33 1 53 00 00 00',
          'https://www.digital-marketing-pro.com',
          'Agence de marketing digital offrant des services de SEO, publicité en ligne et stratégie digitale.',
          'https://www.digital-marketing-pro.com/assets/logo.svg'
      );update_updated_at_column
