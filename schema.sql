-- Database schema for Family Planning Score Tool

-- Geographic hierarchy tables
CREATE TABLE IF NOT EXISTS regions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS districts (
    id SERIAL PRIMARY KEY,
    region_id INTEGER NOT NULL REFERENCES regions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(region_id, name)
);

CREATE TABLE IF NOT EXISTS subcounties (
    id SERIAL PRIMARY KEY,
    district_id INTEGER NOT NULL REFERENCES districts(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(district_id, name)
);

CREATE TABLE IF NOT EXISTS facilities (
    id SERIAL PRIMARY KEY,
    subcounty_id INTEGER NOT NULL REFERENCES subcounties(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(subcounty_id, name)
);

-- Health workers table
CREATE TABLE IF NOT EXISTS health_workers (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone_number VARCHAR(50),
    facility_id INTEGER NOT NULL REFERENCES facilities(id) ON DELETE RESTRICT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Ensure at least one of email or phone is provided
    CONSTRAINT check_email_or_phone CHECK (email IS NOT NULL OR phone_number IS NOT NULL),
    -- Ensure email is unique if provided
    CONSTRAINT unique_email UNIQUE (email),
    -- Ensure phone is unique if provided
    CONSTRAINT unique_phone UNIQUE (phone_number)
);

-- Assessment types
CREATE TABLE IF NOT EXISTS assessment_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    code VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Thematic areas for each assessment type
CREATE TABLE IF NOT EXISTS thematic_areas (
    id SERIAL PRIMARY KEY,
    assessment_type_id INTEGER NOT NULL REFERENCES assessment_types(id) ON DELETE CASCADE,
    name VARCHAR(500) NOT NULL,
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Questions within thematic areas
CREATE TABLE IF NOT EXISTS questions (
    id SERIAL PRIMARY KEY,
    thematic_area_id INTEGER NOT NULL REFERENCES thematic_areas(id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    score_weight INTEGER NOT NULL DEFAULT 2, -- Points for "Yes" response (2, 5, or 10)
    is_critical BOOLEAN DEFAULT FALSE, -- For bold items (10 points)
    is_important BOOLEAN DEFAULT FALSE, -- For asterisk items (5 points)
    display_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Main assessment records
CREATE TABLE IF NOT EXISTS assessments (
    id SERIAL PRIMARY KEY,
    health_worker_id INTEGER NOT NULL REFERENCES health_workers(id) ON DELETE CASCADE,
    facility_id INTEGER NOT NULL REFERENCES facilities(id) ON DELETE RESTRICT,
    assessment_type_id INTEGER NOT NULL REFERENCES assessment_types(id) ON DELETE CASCADE,
    assessor_name VARCHAR(255),
    client_name VARCHAR(255),
    notes TEXT,
    total_possible_score INTEGER NOT NULL DEFAULT 0,
    achieved_score INTEGER NOT NULL DEFAULT 0,
    percentage_score DECIMAL(5,2) NOT NULL DEFAULT 0,
    performance_level VARCHAR(50), -- 'Proficient', 'Competent', 'Not Acceptable'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Individual question responses
CREATE TABLE IF NOT EXISTS assessment_responses (
    id SERIAL PRIMARY KEY,
    assessment_id INTEGER NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    response VARCHAR(10) NOT NULL CHECK (response IN ('Yes', 'No', 'NA')),
    points_earned INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assessment_id, question_id)
);

-- Thematic area scores (for quick reference)
CREATE TABLE IF NOT EXISTS thematic_area_scores (
    id SERIAL PRIMARY KEY,
    assessment_id INTEGER NOT NULL REFERENCES assessments(id) ON DELETE CASCADE,
    thematic_area_id INTEGER NOT NULL REFERENCES thematic_areas(id) ON DELETE CASCADE,
    possible_score INTEGER NOT NULL DEFAULT 0,
    achieved_score INTEGER NOT NULL DEFAULT 0,
    percentage_score DECIMAL(5,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assessment_id, thematic_area_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_districts_region ON districts(region_id);
CREATE INDEX IF NOT EXISTS idx_subcounties_district ON subcounties(district_id);
CREATE INDEX IF NOT EXISTS idx_facilities_subcounty ON facilities(subcounty_id);
CREATE INDEX IF NOT EXISTS idx_thematic_areas_assessment ON thematic_areas(assessment_type_id);
CREATE INDEX IF NOT EXISTS idx_questions_thematic ON questions(thematic_area_id);
CREATE INDEX IF NOT EXISTS idx_health_workers_facility ON health_workers(facility_id);
CREATE INDEX IF NOT EXISTS idx_health_workers_email ON health_workers(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_health_workers_phone ON health_workers(phone_number) WHERE phone_number IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_assessments_health_worker ON assessments(health_worker_id);
CREATE INDEX IF NOT EXISTS idx_assessments_facility ON assessments(facility_id);
CREATE INDEX IF NOT EXISTS idx_assessments_type ON assessments(assessment_type_id);
CREATE INDEX IF NOT EXISTS idx_assessment_responses_assessment ON assessment_responses(assessment_id);
CREATE INDEX IF NOT EXISTS idx_thematic_area_scores_assessment ON thematic_area_scores(assessment_id);

-- Insert default assessment types
INSERT INTO assessment_types (name, code) VALUES
    ('Counselling', 'counselling'),
    ('Progesterone only pill and combined oral contraceptive pill', 'ocps'),
    ('Injectable Progesterone Only', 'injectable_po'),
    ('Implant Insertion', 'implant_insertion'),
    ('Implant Removal', 'implant_removal'),
    ('IUD/IUS Insertion', 'iud_insertion'),
    ('IUD/IUS Removal', 'iud_removal'),
    ('Mini-Laparotomy Tubal Ligation', 'mini_lap'),
    ('Vasectomy', 'vasectomy')
ON CONFLICT (code) DO NOTHING;

-- Auth and administrative areas
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    code VARCHAR(100) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INTEGER NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Administrative areas access (user can be tagged to many areas)
CREATE TABLE IF NOT EXISTS user_admin_areas (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    region_id INTEGER REFERENCES regions(id) ON DELETE CASCADE,
    district_id INTEGER REFERENCES districts(id) ON DELETE CASCADE,
    subcounty_id INTEGER REFERENCES subcounties(id) ON DELETE CASCADE,
    facility_id INTEGER REFERENCES facilities(id) ON DELETE CASCADE
);

-- Useful seed role
INSERT INTO roles (name, description) VALUES ('Admin', 'System administrator')
ON CONFLICT (name) DO NOTHING;

-- Events/Logs table
CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
    session_id VARCHAR(255),
    category VARCHAR(50) NOT NULL, -- 'page', 'interaction', 'form', 'api', 'assessment', 'navigation', 'error', 'auth', 'application', 'filter', 'performance'
    event_type VARCHAR(100) NOT NULL, -- Specific event name
    page VARCHAR(255),
    data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for events table
CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id);
CREATE INDEX IF NOT EXISTS idx_events_category ON events(category);
CREATE INDEX IF NOT EXISTS idx_events_event_type ON events(event_type);
CREATE INDEX IF NOT EXISTS idx_events_created_at ON events(created_at);
CREATE INDEX IF NOT EXISTS idx_events_session_id ON events(session_id);

