-- Migration: 005_create_stages.up.sql
CREATE TABLE contents (
                          id SERIAL PRIMARY KEY,
                          content TEXT
);

-- Exemple minimal de données
INSERT INTO contents (content) VALUES
                                   ('Exemple de contenu 1'),
                                   ('Exemple de contenu 2');

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

                        CONSTRAINT check_dates CHECK (end_date > start_date)
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
    stage_offer_id, student_id, teacher_id, tutor_id, establishment_id,
    status, start_date, end_date
) VALUES
      (
          1, 2, 3, 1, 1,
          'approved',
          '2024-03-01',
          '2024-05-31'
      ),
      (
          1, 2, 3, NULL, 2,
          'in_progress',
          '2024-02-15',
          '2024-04-15'
      );

