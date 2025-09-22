-- Migration: 003_create_users.up.sql

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       first_name VARCHAR(100) NOT NULL,
                       last_name VARCHAR(100) NOT NULL,
                       email VARCHAR(320) UNIQUE NOT NULL CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'moderator', 'user', 'student', 'teacher')),
                       phone VARCHAR(20),
                       establishment_id INTEGER REFERENCES establishments(id) ON DELETE SET NULL,
                       is_active BOOLEAN NOT NULL DEFAULT TRUE,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       last_login TIMESTAMP WITH TIME ZONE
);

-- Index pour améliorer les performances
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_establishment ON users(establishment_id);
CREATE INDEX idx_users_active ON users(is_active);

-- Trigger pour updated_at
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Données d'exemple
INSERT INTO users (
    username, first_name, last_name, email, password_hash, role, phone, establishment_id
) VALUES
      (
          'admin_user',
          'Marie',
          'Dupont',
          'admin@techcorp-solutions.fr',
          '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj4A9LQOS5Wi', -- hash de "password123"
          'admin',
          '+33 6 12 34 56 78',
          1
      ),
      (
          'student_john',
          'John',
          'Martin',
          'john.martin@student.edu',
          '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj4A9LQOS5Wi', -- hash de "password123"
          'student',
          '+33 6 98 65 43 21',
          NULL
      ),
      (
          'teacher_sophie',
          'Sophie',
          'Lefebvre',
          'sophie.lefebvre@digital-marketing-pro.com',
          '$2a$12$LQv3c1yqBWVHxkd0LHAkCOYz6TtxMQJqhN8/LewdBPj4A9LQOS5Wi', -- hash de "password123"
          'teacher',
          '+33 6 55 44 33 22',
          2
      );
