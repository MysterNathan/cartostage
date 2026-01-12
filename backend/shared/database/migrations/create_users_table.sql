-- Migration: 003_create_users.up.sql

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       first_name VARCHAR(100) NOT NULL,
                       last_name VARCHAR(100) NOT NULL,
                       email VARCHAR(320) UNIQUE NOT NULL CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
                       password_hash VARCHAR(255) NOT NULL,
                       role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'tutor', 'student', 'teacher')),
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
      -- Admin (1)
      (
          'admin_user',
          'Marie',
          'Dupont',
          'admin@techcorp-solutions.fr',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'admin',
          '+33 6 12 34 56 78',
          1
      ),

      -- Students (10)
      (
          'student_lucas',
          'Lucas',
          'Bernard',
          'lucas.bernard@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 11 22 33 44',
          NULL
      ),
      (
          'student_emma',
          'Emma',
          'Petit',
          'emma.petit@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 22 33 44 55',
          NULL
      ),
      (
          'student_theo',
          'Théo',
          'Dubois',
          'theo.dubois@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 33 44 55 66',
          NULL
      ),
      (
          'student_lea',
          'Léa',
          'Moreau',
          'lea.moreau@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 44 55 66 77',
          NULL
      ),
      (
          'student_hugo',
          'Hugo',
          'Simon',
          'hugo.simon@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 55 66 77 88',
          NULL
      ),
      (
          'student_chloe',
          'Chloé',
          'Laurent',
          'chloe.laurent@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 66 77 88 99',
          NULL
      ),
      (
          'student_arthur',
          'Arthur',
          'Lefevre',
          'arthur.lefevre@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 77 88 99 00',
          NULL
      ),
      (
          'student_camille',
          'Camille',
          'Roux',
          'camille.roux@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 88 99 00 11',
          NULL
      ),
      (
          'student_louis',
          'Louis',
          'Fournier',
          'louis.fournier@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 99 00 11 22',
          NULL
      ),
      (
          'student_alice',
          'Alice',
          'Girard',
          'alice.girard@eleve.edu',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'eleve',
          '+33 6 00 11 22 33',
          NULL
      ),

      -- Teachers (5)
      (
          'teacher_pierre',
          'Pierre',
          'Bonnet',
          'pierre.bonnet@digital-marketing-pro.com',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'teacher',
          '+33 6 10 20 30 40',
          2
      ),
      (
          'teacher_julie',
          'Julie',
          'Mercier',
          'julie.mercier@techcorp-solutions.fr',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'teacher',
          '+33 6 20 30 40 50',
          1
      ),
      (
          'teacher_marc',
          'Marc',
          'Vincent',
          'marc.vincent@digital-marketing-pro.com',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'teacher',
          '+33 6 30 40 50 60',
          2
      ),
      (
          'teacher_claire',
          'Claire',
          'Rousseau',
          'claire.rousseau@techcorp-solutions.fr',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'teacher',
          '+33 6 40 50 60 70',
          1
      ),
      (
          'teacher_nicolas',
          'Nicolas',
          'Blanc',
          'nicolas.blanc@digital-marketing-pro.com',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'teacher',
          '+33 6 50 60 70 80',
          2
      ),

      -- Tutors (5)
      (
          'tutor_sophie',
          'Sophie',
          'Garnier',
          'sophie.garnier@techcorp-solutions.fr',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'tutor',
          '+33 6 60 70 80 90',
          1
      ),
      (
          'tutor_francois',
          'François',
          'Faure',
          'francois.faure@digital-marketing-pro.com',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'tutor',
          '+33 6 70 80 90 01',
          2
      ),
      (
          'tutor_isabelle',
          'Isabelle',
          'Andre',
          'isabelle.andre@techcorp-solutions.fr',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'tutor',
          '+33 6 80 90 01 12',
          1
      ),
      (
          'tutor_antoine',
          'Antoine',
          'Lambert',
          'antoine.lambert@digital-marketing-pro.com',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'tutor',
          '+33 6 90 01 12 23',
          2
      ),
      (
          'tutor_valerie',
          'Valérie',
          'Fontaine',
          'valerie.fontaine@techcorp-solutions.fr',
          '$2y$10$1o9YjGNkF.RjDxbu0uAkKub/WR61BcyCupHgtX904FcPSpPk60Cei',
          'tutor',
          '+33 6 01 12 23 34',
          1
      );
