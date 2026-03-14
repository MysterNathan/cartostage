/*
*************************************************************************
                                TYPES
*************************************************************************
 */
CREATE TYPE user_type AS ENUM ('TEACHER', 'STUDENT', 'TUTOR');
/*
 TABLES
 */
CREATE TABLE IF NOT EXISTS form (
                                    id SERIAL PRIMARY KEY NOT NULL,
                                    stage_id INTEGER NOT NULL REFERENCES stages(id) ON DELETE CASCADE,
                                    status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
                                    content JSONB,
                                    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                    completed_at TIMESTAMP,

                                    CONSTRAINT fk_stage FOREIGN KEY (stage_id) REFERENCES stages(id)
);

CREATE TABLE IF NOT EXISTS form_section (
                                    id SERIAL PRIMARY KEY,
                                    form_id INTEGER NOT NULL REFERENCES form(id) ON DELETE CASCADE,
                                    section_type user_type NOT NULL ,
                                    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
                                    status VARCHAR(50) NOT NULL DEFAULT 'CREATED',
                                    content JSONB,
                                    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                    completed_at TIMESTAMP,

                                    CONSTRAINT fk_form FOREIGN KEY (form_id) REFERENCES form(id),
                                    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id)
);
/*
*************************************************************************
                                INDEX
*************************************************************************
 */
-- Index pour optimiser les requêtes fréquentes
CREATE INDEX idx_form_stage_id ON form(stage_id);
CREATE INDEX idx_form_status ON form(status);
CREATE INDEX idx_form_created_at ON form(created_at);
CREATE INDEX idx_form_content ON form USING GIN (content);

-- Index pour optimiser les requêtes fréquentes
CREATE INDEX idx_form_stage_id ON form_section(form_id);
CREATE INDEX idx_form_student_id ON form_section(user_id);
CREATE INDEX idx_form_status ON form_section(status);
CREATE INDEX idx_form_created_at ON form_section(created_at);
CREATE INDEX idx_form_content ON form_section USING GIN (content);

/*
*************************************************************************
                                RLS
*************************************************************************
 */
ALTER TABLE form  FORCE ROW LEVEL SECURITY;
ALTER TABLE form_section FORCE ROW LEVEL SECURITY;

/*
*************************************************************************
                                FUNCTION TRIGGER
*************************************************************************
 */
CREATE OR REPLACE FUNCTION update_form_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION is_valid_transition_form_status()
    RETURNS TRIGGER
AS $$
BEGIN
    IF NEW.status = OLD.status THEN
        RETURN new;
    END IF;

    IF OLD.status = 'CREATED' AND NEW.status = 'IN_PROGRESS' THEN
        RETURN NEW;
    ELSIF OLD.status = 'IN_PROGRESS' AND NEW.status = 'COMPLETED' THEN
        RETURN NEW;
    ELSIF OLD.status = 'COMPLETED' AND NEW.status = 'IN_PROGRESS' THEN
        RETURN NEW;
    ELSE
        RAISE EXCEPTION
            'Forbiden status transition: % -> %', OLD.status, NEW.status;
    END IF;
END;
$$
    LANGUAGE plpgsql;

CREATE FUNCTION update_teacher_form()
    RETURNS TRIGGER
AS $$
BEGIN
    IF OLD.stage_id != NEW.stage_id  THEN
        RAISE EXCEPTION
            'Forbiden operation, cannot modify forms users id';
    ELSIF OLD.created_at != new.created_at THEN
        RAISE EXCEPTION
            'Forbiden operation, cannot modify created_at value';
    END IF;
    RETURN NEW;
END;
$$
    LANGUAGE plpgsql;
/*
*************************************************************************
                                FUNCTION
*************************************************************************
 */

CREATE FUNCTION allow_update_form_section_user ()
    RETURNS boolean
    STABLE
AS $$
BEGIN
    RETURN EXISTS(SELECT 1 FROM form_section WHERE user_id = current_setting('app.user_id')::int);
END;
$$
    LANGUAGE plpgsql;

CREATE FUNCTION allow_update_form_section_status (form_identifier integer)
    RETURNS boolean
AS $$
BEGIN
    RETURN EXISTS(SELECT 1 FROM form_section WHERE id = form_identifier AND status != 'COMPLETED');
END;
$$
    LANGUAGE plpgsql;


CREATE FUNCTION is_current_user_teacher()
    RETURNS boolean
    STABLE
AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1
        FROM users
        WHERE users.id = current_setting('app.user_id')::int
          AND users.role = 'teacher'
    );
END;
$$
    LANGUAGE plpgsql;


/*
*************************************************************************
                                TRIGGER
*************************************************************************
 */

CREATE TRIGGER trigger_update_form_updated_at
    BEFORE UPDATE ON form
    FOR EACH ROW
EXECUTE FUNCTION update_form_updated_at();

CREATE TRIGGER transition_form_status
    BEFORE UPDATE ON form
    FOR EACH ROW
EXECUTE FUNCTION is_valid_transition_form_status();

CREATE TRIGGER valid_form_update_teacher
    BEFORE UPDATE ON form
    FOR EACH ROW
EXECUTE FUNCTION update_teacher_form();

CREATE TRIGGER transition_form_section_status
    BEFORE UPDATE ON form_section
    FOR EACH ROW
EXECUTE FUNCTION is_valid_transition_form_status();

/*
*************************************************************************
                                POLICY
*************************************************************************
 */

CREATE POLICY only_teacher_update_form
    ON form
    FOR UPDATE
    USING (
    EXISTS (
        SELECT 1 FROM stages s
        WHERE s.id = form.stage_id
          AND s.teacher_id = current_setting('app.user_id')::int
    )
        AND is_current_user_teacher()
    )
    WITH CHECK (
    stage_id = (SELECT stage_id FROM form WHERE id = form.id)
    );

CREATE POLICY only_owner_update_form_section
    ON form_section
    FOR UPDATE
    USING(
    user_id = current_setting('app.user_id')::int
    )
    WITH CHECK (
    user_id = current_setting('app.user_id')::int
        AND section_type = form_section.section_type
        AND created_at = form_section.created_at
        AND updated_at = now()
    );

/*
*************************************************************************
                                SAMPLE DATAS
*************************************************************************
 */

INSERT INTO form (
    stage_id,
    status,
    content,
    created_at,
    completed_at
) VALUES
      (1,  'CREATED',
       '{"objective":"Découverte du monde professionnel","company":"TechCorp","duration_weeks":4}',
       NOW() - INTERVAL '5 days',
       NULL),

      (2, 'IN_PROGRESS',
       '{"objective":"Stage développement web","company":"Digital Marketing Pro","duration_weeks":6}',
       NOW() - INTERVAL '15 days',
       NULL),

      (3 , 'COMPLETED',
       '{"objective":"Stage data analysis","company":"DataSolutions","duration_weeks":8}',
       NOW() - INTERVAL '40 days',
       NOW() - INTERVAL '7 days'),

      (4, 'IN_PROGRESS',
       '{"objective":"Stage marketing","company":"MarketPlus","duration_weeks":6}',
       NOW() - INTERVAL '20 days',
       NULL),

      (5, 'CREATED',
       '{"objective":"Stage réseaux","company":"NetSecure","duration_weeks":5}',
       NOW() - INTERVAL '3 days',
       NULL);

INSERT INTO form_section (
    form_id,
    section_type,
    user_id,
    status,
    content,
    created_at,
    completed_at
) VALUES
      (1, 'STUDENT', 2, 'CREATED',
       '{"motivation":"Très motivé pour découvrir le milieu professionnel"}',
       NOW() - INTERVAL '4 days',
       NULL),

      (2, 'TEACHER', 13, 'IN_PROGRESS',
       '{"evaluation":"Bon investissement","remarks":"Continue sur cette voie"}',
       NOW() - INTERVAL '13 days',
       NULL),

      (3, 'TUTOR', 19, 'COMPLETED',
       '{"integration":"Excellente","autonomy":"Très bonne","skills":"Analyse de données"}',
       NOW() - INTERVAL '35 days',
       NOW() - INTERVAL '8 days'),

      (4, 'STUDENT', 5, 'IN_PROGRESS',
       '{"feedback":"Stage intéressant","difficulties":"Gestion du temps"}',
       NOW() - INTERVAL '18 days',
       NULL),

      (5, 'TEACHER', 16, 'CREATED',
       '{"comment":"Formulaire en attente de complétion"}',
       NOW() - INTERVAL '2 days',
       NULL);
