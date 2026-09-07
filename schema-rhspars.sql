-- RH SPARS schema: Integrated Reproductive Health Support Supervision Tool
-- Parallel to FP proficiency assessments, with binary Yes=1/No=0/NA scoring.

CREATE TABLE IF NOT EXISTS rhspars_domains (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rhspars_thematic_areas (
    id SERIAL PRIMARY KEY,
    domain_id INTEGER NOT NULL REFERENCES rhspars_domains(id) ON DELETE CASCADE,
    name VARCHAR(500) NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rhspars_questions (
    id SERIAL PRIMARY KEY,
    thematic_area_id INTEGER NOT NULL REFERENCES rhspars_thematic_areas(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    score_weight INTEGER NOT NULL DEFAULT 1, -- Yes earns this many points (typically 1)
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- One supervision visit / assessment at a facility
CREATE TABLE IF NOT EXISTS rhspars_assessments (
    id SERIAL PRIMARY KEY,
    facility_id INTEGER NOT NULL REFERENCES facilities(id) ON DELETE RESTRICT,
    supervision_type VARCHAR(50), -- 'self' | 'external'
    supervision_date DATE,
    next_supervision_date DATE,
    rh_supplier TEXT,
    facility_type_detail TEXT,
    facility_ownership TEXT,
    primary_contact_name TEXT,
    primary_contact_phone TEXT,
    primary_contact_email TEXT,
    completed_by_name TEXT,
    completed_by_designation TEXT,
    completed_by_institution TEXT,
    assessment_completion_date DATE,
    notes TEXT,
    -- Rollups
    maternity_pct DECIMAL(6,2) NOT NULL DEFAULT 0,
    fp_pct DECIMAL(6,2) NOT NULL DEFAULT 0,
    anc_pct DECIMAL(6,2) NOT NULL DEFAULT 0,
    commodities_pct DECIMAL(6,2) NOT NULL DEFAULT 0,
    health_info_pct DECIMAL(6,2) NOT NULL DEFAULT 0,
    general_facility_pct DECIMAL(6,2) NOT NULL DEFAULT 0,
    grand_percentage DECIMAL(6,2) NOT NULL DEFAULT 0,
    total_possible_score INTEGER NOT NULL DEFAULT 0,
    achieved_score INTEGER NOT NULL DEFAULT 0,
    performance_level VARCHAR(50) NOT NULL DEFAULT 'Not Acceptable',
    created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS rhspars_attendees (
    id SERIAL PRIMARY KEY,
    assessment_id INTEGER NOT NULL REFERENCES rhspars_assessments(id) ON DELETE CASCADE,
    person_type VARCHAR(20) NOT NULL, -- 'present' | 'supervisor'
    full_name TEXT,
    cadre_position TEXT,
    affiliation TEXT,
    contact TEXT
);

CREATE TABLE IF NOT EXISTS rhspars_responses (
    id SERIAL PRIMARY KEY,
    assessment_id INTEGER NOT NULL REFERENCES rhspars_assessments(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES rhspars_questions(id) ON DELETE CASCADE,
    response VARCHAR(10) NOT NULL CHECK (response IN ('Yes', 'No', 'NA')),
    points_earned INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assessment_id, question_id)
);

CREATE TABLE IF NOT EXISTS rhspars_section_scores (
    id SERIAL PRIMARY KEY,
    assessment_id INTEGER NOT NULL REFERENCES rhspars_assessments(id) ON DELETE CASCADE,
    thematic_area_id INTEGER NOT NULL REFERENCES rhspars_thematic_areas(id) ON DELETE CASCADE,
    possible_score INTEGER NOT NULL DEFAULT 0,
    achieved_score INTEGER NOT NULL DEFAULT 0,
    percentage_score DECIMAL(6,2) NOT NULL DEFAULT 0,
    UNIQUE(assessment_id, thematic_area_id)
);

CREATE TABLE IF NOT EXISTS rhspars_domain_scores (
    id SERIAL PRIMARY KEY,
    assessment_id INTEGER NOT NULL REFERENCES rhspars_assessments(id) ON DELETE CASCADE,
    domain_id INTEGER NOT NULL REFERENCES rhspars_domains(id) ON DELETE CASCADE,
    percentage_score DECIMAL(6,2) NOT NULL DEFAULT 0,
    UNIQUE(assessment_id, domain_id)
);

CREATE TABLE IF NOT EXISTS rhspars_action_plan_items (
    id SERIAL PRIMARY KEY,
    assessment_id INTEGER NOT NULL REFERENCES rhspars_assessments(id) ON DELETE CASCADE,
    gap_identified TEXT,
    action_recommendation TEXT,
    means_of_verification TEXT,
    responsible_person TEXT,
    time_frame TEXT,
    display_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_rhspars_ta_domain ON rhspars_thematic_areas(domain_id);
CREATE INDEX IF NOT EXISTS idx_rhspars_q_ta ON rhspars_questions(thematic_area_id);
CREATE INDEX IF NOT EXISTS idx_rhspars_assess_facility ON rhspars_assessments(facility_id);
CREATE INDEX IF NOT EXISTS idx_rhspars_resp_assess ON rhspars_responses(assessment_id);
CREATE INDEX IF NOT EXISTS idx_rhspars_section_assess ON rhspars_section_scores(assessment_id);
CREATE INDEX IF NOT EXISTS idx_rhspars_domain_assess ON rhspars_domain_scores(assessment_id);
