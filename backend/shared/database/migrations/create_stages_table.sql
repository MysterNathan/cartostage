-- Migration: 005_create_stages.up.sql

CREATE TABLE contents (
                          id SERIAL PRIMARY KEY,
                          content TEXT
);

-- Exemple de contenus
INSERT INTO contents (content) VALUES
                                   ('Développement d''une application web complète avec React et Node.js'),
                                   ('Campagne marketing digital multi-canaux'),
                                   ('Création d''une application mobile native'),
                                   ('Gestion du processus de recrutement'),
                                   ('Refonte de l''interface utilisateur du site web'),
                                   ('Analyse de données clients et reporting'),
                                   ('Animation des réseaux sociaux de l''entreprise');

CREATE TABLE stages (
                        id SERIAL PRIMARY KEY,
                        stage_offer_id INTEGER NOT NULL REFERENCES stage_offers(id) ON DELETE CASCADE,
                        student_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                        teacher_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                        tutor_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                        establishment_id INTEGER REFERENCES establishments(id) ON DELETE SET NULL,
                        content_id INTEGER REFERENCES contents(id) ON DELETE SET NULL,
                        status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'in_progress', 'completed', 'cancelled')),
                        start_date DATE NOT NULL,
                        end_date DATE NOT NULL,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

                        CONSTRAINT check_dates CHECK (end_date > start_date),
                        CONSTRAINT unique_student_stage UNIQUE(student_id, stage_offer_id, start_date)
);

-- Index pour améliorer les performances
CREATE INDEX idx_stages_stage_offer ON stages(stage_offer_id);
CREATE INDEX idx_stages_student ON stages(student_id);
CREATE INDEX idx_stages_teacher ON stages(teacher_id);
CREATE INDEX idx_stages_tutor ON stages(tutor_id);
CREATE INDEX idx_stages_establishment ON stages(establishment_id);
CREATE INDEX idx_stages_status ON stages(status);
CREATE INDEX idx_stages_dates ON stages(start_date, end_date);

-- Trigger pour updated_at
CREATE TRIGGER update_stages_updated_at
    BEFORE UPDATE ON stages
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Données d'exemple
INSERT INTO stages (
    stage_offer_id, student_id, teacher_id, tutor_id, establishment_id, content_id,
    status, start_date, end_date
) VALUES
      -- Offre 1 (TechCorp Solutions - Développeur Full Stack) : 3 étudiants avec le même tuteur (Sophie - id 17)
      (
          1, 2, 12, 17, 1, 1,
          'in_progress',
          '2024-02-01',
          '2024-07-31'
      ),
      (
          1, 3, 12, 17, 1, 1,
          'in_progress',
          '2024-02-01',
          '2024-07-31'
      ),
      (
          1, 4, 13, 17, 1, 1,
          'approved',
          '2024-02-15',
          '2024-08-15'
      ),

      -- Offre 2 (Digital Marketing Pro - Assistant Marketing) : 2 étudiants avec le même tuteur (François - id 18)
      (
          2, 5, 14, 18, 2, 2,
          'in_progress',
          '2024-03-01',
          '2024-06-30'
      ),
      (
          2, 6, 15, 18, 2, 2,
          'in_progress',
          '2024-03-01',
          '2024-06-30'
      ),

      -- Offre 3 (MobileTech Industries - Développeur Mobile) : 2 étudiants avec le même tuteur (Isabelle - id 19)
      (
          3, 7, 12, 19, 1, 3,
          'approved',
          '2024-01-15',
          '2024-06-15'
      ),
      (
          3, 8, 13, 19, 1, 3,
          'in_progress',
          '2024-01-15',
          '2024-06-15'
      ),

      -- Offre 4 (HR Consulting Group - Assistant RH) : 1 étudiant
      (
          4, 9, 14, 20, 2, 4,
          'in_progress',
          '2024-04-01',
          '2024-06-30'
      ),

      -- Offre 5 (Creative Studio - Designer UX/UI) : 1 étudiant
      (
          5, 10, 15, 17, 1, 5,
          'approved',
          '2024-05-01',
          '2024-08-31'
      ),

      -- Offre 6 (DataViz Solutions - Data Analyst) : 1 étudiant
      (
          6, 11, 16, 21, 1, 6,
          'in_progress',
          '2024-02-01',
          '2024-07-31'
      ),

      -- Offre 7 (Social Media Agency - Community Manager) : 2 étudiants avec le même tuteur (Antoine - id 20)
      (
          7, 2, 14, 20, 2, 7,
          'pending',
          '2024-06-01',
          '2024-08-31'
      ),
      (
          7, 3, 15, 20, 2, 7,
          'pending',
          '2024-06-01',
          '2024-08-31'
      );
