-- Migration script to add health workers support to existing database
-- Run this after running the updated schema.sql

-- Step 1: Create health_workers table if it doesn't exist
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

-- Step 2: Add health_worker_id column to assessments table
-- First, check if column already exists
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'assessments' AND column_name = 'health_worker_id'
    ) THEN
        -- Add the column as nullable first (we'll populate it, then make it NOT NULL)
        ALTER TABLE assessments 
        ADD COLUMN health_worker_id INTEGER;
        
        -- Add foreign key constraint
        ALTER TABLE assessments 
        ADD CONSTRAINT fk_assessments_health_worker 
        FOREIGN KEY (health_worker_id) REFERENCES health_workers(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Step 3: Create a temporary default health worker for each facility
-- This allows us to assign existing assessments to health workers
DO $$
DECLARE
    facility_rec RECORD;
    default_hw_id INTEGER;
    phone_num TEXT;
    email_addr TEXT;
BEGIN
    FOR facility_rec IN SELECT DISTINCT id, name FROM facilities LOOP
        -- Generate unique identifiers for this facility
        email_addr := 'default_' || facility_rec.id || '@example.com';
        -- Use a larger padding to avoid conflicts, and include facility_id in a way that's unique
        phone_num := '+256700' || LPAD(facility_rec.id::TEXT, 6, '0');
        
        -- Check if a default health worker already exists for this facility
        SELECT id INTO default_hw_id 
        FROM health_workers 
        WHERE facility_id = facility_rec.id 
        AND (email = email_addr OR email LIKE 'default_%')
        LIMIT 1;
        
        -- If no default health worker exists, create one
        IF default_hw_id IS NULL THEN
            -- Try to insert, handling conflicts on both email and phone
            BEGIN
                INSERT INTO health_workers (full_name, email, phone_number, facility_id)
                VALUES (
                    'Default Health Worker - ' || facility_rec.name,
                    email_addr,
                    phone_num,
                    facility_rec.id
                )
                RETURNING id INTO default_hw_id;
            EXCEPTION WHEN unique_violation THEN
                -- If there's a conflict, try to find an existing health worker for this facility
                SELECT id INTO default_hw_id 
                FROM health_workers 
                WHERE facility_id = facility_rec.id 
                LIMIT 1;
                
                -- If still no health worker, create one with a different phone number
                IF default_hw_id IS NULL THEN
                    phone_num := '+256701' || LPAD(facility_rec.id::TEXT, 6, '0');
                    INSERT INTO health_workers (full_name, email, phone_number, facility_id)
                    VALUES (
                        'Default Health Worker - ' || facility_rec.name,
                        email_addr,
                        phone_num,
                        facility_rec.id
                    )
                    RETURNING id INTO default_hw_id;
                END IF;
            END;
        END IF;
        
        -- Assign all assessments for this facility to the default health worker
        IF default_hw_id IS NOT NULL THEN
            UPDATE assessments 
            SET health_worker_id = default_hw_id
            WHERE facility_id = facility_rec.id 
            AND health_worker_id IS NULL;
        END IF;
    END LOOP;
END $$;

-- Step 4: Make health_worker_id NOT NULL now that all assessments have been assigned
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'assessments' 
        AND column_name = 'health_worker_id' 
        AND is_nullable = 'YES'
    ) THEN
        -- Check if all assessments have health_worker_id
        IF NOT EXISTS (SELECT 1 FROM assessments WHERE health_worker_id IS NULL) THEN
            ALTER TABLE assessments 
            ALTER COLUMN health_worker_id SET NOT NULL;
        ELSE
            RAISE NOTICE 'Warning: Some assessments still have NULL health_worker_id. Please assign them before making the column NOT NULL.';
        END IF;
    END IF;
END $$;

-- Step 5: Add indexes for performance
CREATE INDEX IF NOT EXISTS idx_health_workers_facility ON health_workers(facility_id);
CREATE INDEX IF NOT EXISTS idx_health_workers_email ON health_workers(email) WHERE email IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_health_workers_phone ON health_workers(phone_number) WHERE phone_number IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_assessments_health_worker ON assessments(health_worker_id);

-- Step 6: Add health workers permissions
INSERT INTO permissions (code, description) VALUES
    ('health_workers.view', 'View health workers'),
    ('health_workers.create', 'Create health workers'),
    ('health_workers.edit', 'Edit health workers'),
    ('health_workers.delete', 'Delete health workers')
ON CONFLICT (code) DO NOTHING;

