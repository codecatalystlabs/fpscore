-- Seed assessment data for all thematic areas
-- This script creates sample assessments with responses for all questions across all assessment types

-- Helper function to create an assessment with responses
-- We'll use a DO block to create assessments programmatically

DO $$
DECLARE
    facility_rec RECORD;
    assessment_type_rec RECORD;
    question_rec RECORD;
    thematic_area_rec RECORD;
    assessment_id_var INTEGER;
    total_possible INTEGER;
    achieved_score_var INTEGER;
    percentage_score_var DECIMAL(5,2);
    performance_level_var VARCHAR(50);
    ta_possible INTEGER;
    ta_achieved INTEGER;
    ta_percentage DECIMAL(5,2);
    response_val VARCHAR(10);
    points_earned_var INTEGER;
    rand_val INTEGER;
    assessor_names TEXT[] := ARRAY['Dr. Sarah Nakato', 'Dr. James Ochieng', 'Nurse Mary Achieng', 'Dr. Peter Okello', 'Nurse Grace Atim'];
    client_names TEXT[] := ARRAY['Client A', 'Client B', 'Client C', 'Client D', 'Client E'];
    assessment_counter INTEGER := 0;
    health_worker_id_var INTEGER;
BEGIN
    -- Loop through each facility
    FOR facility_rec IN SELECT id FROM facilities ORDER BY id LIMIT 5 LOOP
        -- Get a random health worker for this facility (or create one if none exist)
        SELECT id INTO health_worker_id_var 
        FROM health_workers 
        WHERE facility_id = facility_rec.id 
        ORDER BY RANDOM() 
        LIMIT 1;
        
        -- If no health worker exists for this facility, skip it
        IF health_worker_id_var IS NULL THEN
            RAISE NOTICE 'No health workers found for facility %, skipping assessments', facility_rec.id;
            CONTINUE;
        END IF;
        
        -- Loop through each assessment type
        FOR assessment_type_rec IN SELECT id, name FROM assessment_types ORDER BY id LOOP
            assessment_counter := assessment_counter + 1;
            
            -- Initialize scores
            total_possible := 0;
            achieved_score_var := 0;
            
            -- Create assessment record
            INSERT INTO assessments (
                health_worker_id,
                facility_id,
                assessment_type_id,
                assessor_name,
                client_name,
                notes,
                total_possible_score,
                achieved_score,
                percentage_score,
                performance_level,
                created_at
            ) VALUES (
                health_worker_id_var,
                facility_rec.id,
                assessment_type_rec.id,
                assessor_names[1 + (assessment_counter % array_length(assessor_names, 1))],
                client_names[1 + (assessment_counter % array_length(client_names, 1))],
                'Sample assessment for ' || assessment_type_rec.name,
                0,
                0,
                0.0,
                'Not Acceptable',
                CURRENT_TIMESTAMP - (assessment_counter || ' days')::INTERVAL
            ) RETURNING id INTO assessment_id_var;
            
            -- Loop through each thematic area for this assessment type
            FOR thematic_area_rec IN 
                SELECT id, name 
                FROM thematic_areas 
                WHERE assessment_type_id = assessment_type_rec.id 
                ORDER BY display_order 
            LOOP
                ta_possible := 0;
                ta_achieved := 0;
                
                -- Loop through each question in this thematic area
                FOR question_rec IN 
                    SELECT id, score_weight, is_critical
                    FROM questions 
                    WHERE thematic_area_id = thematic_area_rec.id 
                    ORDER BY display_order 
                LOOP
                    -- Calculate response based on question type
                    -- Critical questions (10 points) - 80% Yes, 15% No, 5% NA
                    -- Important questions (5 points) - 70% Yes, 25% No, 5% NA  
                    -- Regular questions (2 points) - 60% Yes, 35% No, 5% NA
                    -- This creates varied performance levels
                    
                    IF question_rec.is_critical THEN
                        -- Critical questions: higher chance of Yes
                        response_val := CASE (assessment_counter + question_rec.id) % 10
                            WHEN 0 THEN 'NA'
                            WHEN 1 THEN 'NA'
                            WHEN 2 THEN 'No'
                            ELSE 'Yes'
                        END;
                    ELSIF question_rec.score_weight = 5 THEN
                        -- Important questions
                        response_val := CASE (assessment_counter + question_rec.id) % 10
                            WHEN 0 THEN 'NA'
                            WHEN 1 THEN 'No'
                            WHEN 2 THEN 'No'
                            WHEN 3 THEN 'No'
                            ELSE 'Yes'
                        END;
                    ELSE
                        -- Regular questions
                        response_val := CASE (assessment_counter + question_rec.id) % 10
                            WHEN 0 THEN 'NA'
                            WHEN 1 THEN 'No'
                            WHEN 2 THEN 'No'
                            WHEN 3 THEN 'No'
                            WHEN 4 THEN 'No'
                            ELSE 'Yes'
                        END;
                    END IF;
                    
                    -- Calculate points earned
                    -- Only "Yes" responses get points (score_weight)
                    -- "No" and "NA" get 0 points
                    IF response_val = 'Yes' THEN
                        points_earned_var := question_rec.score_weight;
                    ELSE
                        points_earned_var := 0;
                    END IF;
                    
                    -- Add to totals
                    -- Only "Yes" and "No" count toward possible score
                    -- "NA" questions don't count toward possible score at all
                    IF response_val = 'Yes' THEN
                        total_possible := total_possible + question_rec.score_weight;
                        achieved_score_var := achieved_score_var + points_earned_var;
                        ta_possible := ta_possible + question_rec.score_weight;
                        ta_achieved := ta_achieved + points_earned_var;
                    ELSIF response_val = 'No' THEN
                        -- "No" counts toward possible but not achieved
                        total_possible := total_possible + question_rec.score_weight;
                        ta_possible := ta_possible + question_rec.score_weight;
                    END IF;
                    -- "NA" responses don't contribute to scores at all
                    
                    -- Insert response
                    INSERT INTO assessment_responses (
                        assessment_id,
                        question_id,
                        response,
                        points_earned
                    ) VALUES (
                        assessment_id_var,
                        question_rec.id,
                        response_val,
                        points_earned_var
                    );
                END LOOP;
                
                -- Calculate thematic area percentage
                IF ta_possible > 0 THEN
                    ta_percentage := (ta_achieved::DECIMAL / ta_possible::DECIMAL) * 100.0;
                ELSE
                    ta_percentage := 0.0;
                END IF;
                
                -- Insert thematic area score
                INSERT INTO thematic_area_scores (
                    assessment_id,
                    thematic_area_id,
                    possible_score,
                    achieved_score,
                    percentage_score
                ) VALUES (
                    assessment_id_var,
                    thematic_area_rec.id,
                    ta_possible,
                    ta_achieved,
                    ta_percentage
                );
            END LOOP;
            
            -- Calculate overall percentage and performance level
            IF total_possible > 0 THEN
                percentage_score_var := (achieved_score_var::DECIMAL / total_possible::DECIMAL) * 100.0;
            ELSE
                percentage_score_var := 0.0;
            END IF;
            
            -- Determine performance level
            IF percentage_score_var > 90 THEN
                performance_level_var := 'Proficient';
            ELSIF percentage_score_var >= 70 THEN
                performance_level_var := 'Competent';
            ELSE
                performance_level_var := 'Not Acceptable';
            END IF;
            
            -- Update assessment with calculated scores
            UPDATE assessments
            SET total_possible_score = total_possible,
                achieved_score = achieved_score_var,
                percentage_score = percentage_score_var,
                performance_level = performance_level_var
            WHERE id = assessment_id_var;
            
        END LOOP;
    END LOOP;
    
    RAISE NOTICE 'Created % assessments', assessment_counter;
END $$;

-- Create additional assessments with varying performance levels for better testing
-- High performing assessments (Proficient)
DO $$
DECLARE
    facility_rec RECORD;
    assessment_type_rec RECORD;
    question_rec RECORD;
    thematic_area_rec RECORD;
    assessment_id_var INTEGER;
    total_possible INTEGER;
    achieved_score_var INTEGER;
    percentage_score_var DECIMAL(5,2);
    performance_level_var VARCHAR(50);
    ta_possible INTEGER;
    ta_achieved INTEGER;
    ta_percentage DECIMAL(5,2);
    response_val VARCHAR(10);
    points_earned_var INTEGER;
    rand_val INTEGER;
    facility_count INTEGER;
    type_count INTEGER;
    health_worker_id_var INTEGER;
BEGIN
    -- Get first facility and first 3 assessment types for high performers
    SELECT id INTO facility_rec FROM facilities ORDER BY id LIMIT 1;
    
    -- Get a health worker for this facility
    SELECT id INTO health_worker_id_var 
    FROM health_workers 
    WHERE facility_id = facility_rec.id 
    ORDER BY RANDOM() 
    LIMIT 1;
    
    -- Skip if no health worker exists
    IF health_worker_id_var IS NULL THEN
        RAISE NOTICE 'No health workers found for facility %, skipping high performer assessments', facility_rec.id;
        RETURN;
    END IF;
    
    FOR assessment_type_rec IN SELECT id, name FROM assessment_types ORDER BY id LIMIT 3 LOOP
        total_possible := 0;
            achieved_score_var := 0;
        
        INSERT INTO assessments (
            health_worker_id,
            facility_id,
            assessment_type_id,
            assessor_name,
            client_name,
            notes,
            total_possible_score,
            achieved_score,
            percentage_score,
            performance_level,
            created_at
        ) VALUES (
            health_worker_id_var,
            facility_rec.id,
            assessment_type_rec.id,
            'Dr. Excellent Provider',
            'High Performer Client',
            'High performing assessment - Proficient level',
            0,
            0,
            0.0,
            'Proficient',
            CURRENT_TIMESTAMP - (random() * 30 || ' days')::INTERVAL
        ) RETURNING id INTO assessment_id_var;
        
        FOR thematic_area_rec IN 
            SELECT id, name 
            FROM thematic_areas 
            WHERE assessment_type_id = assessment_type_rec.id 
            ORDER BY display_order 
        LOOP
            ta_possible := 0;
            ta_achieved := 0;
            
            FOR question_rec IN 
                SELECT id, score_weight, is_critical
                FROM questions 
                WHERE thematic_area_id = thematic_area_rec.id 
                ORDER BY display_order 
            LOOP
                -- High performer: 95% Yes, 3% No, 2% NA
                rand_val := floor(random() * 100);
                IF rand_val < 2 THEN
                    response_val := 'NA';
                ELSIF rand_val < 5 THEN
                    response_val := 'No';
                ELSE
                    response_val := 'Yes';
                END IF;
                
                -- Only "Yes" responses get points (score_weight)
                IF response_val = 'Yes' THEN
                    points_earned_var := question_rec.score_weight;
                ELSE
                    points_earned_var := 0;
                END IF;
                
                -- Add to totals
                -- Only "Yes" and "No" count toward possible score
                -- "NA" questions don't count toward possible score at all
                IF response_val = 'Yes' THEN
                    total_possible := total_possible + question_rec.score_weight;
                        achieved_score_var := achieved_score_var + points_earned_var;
                    ta_possible := ta_possible + question_rec.score_weight;
                    ta_achieved := ta_achieved + points_earned_var;
                ELSIF response_val = 'No' THEN
                    -- "No" counts toward possible but not achieved
                    total_possible := total_possible + question_rec.score_weight;
                    ta_possible := ta_possible + question_rec.score_weight;
                END IF;
                -- "NA" responses don't contribute to scores at all
                
                INSERT INTO assessment_responses (
                    assessment_id,
                    question_id,
                    response,
                    points_earned
                ) VALUES (
                    assessment_id_var,
                    question_rec.id,
                    response_val,
                    points_earned_var
                );
            END LOOP;
            
            IF ta_possible > 0 THEN
                ta_percentage := (ta_achieved::DECIMAL / ta_possible::DECIMAL) * 100.0;
            ELSE
                ta_percentage := 0.0;
            END IF;
            
            INSERT INTO thematic_area_scores (
                assessment_id,
                thematic_area_id,
                possible_score,
                achieved_score,
                percentage_score
            ) VALUES (
                assessment_id_var,
                thematic_area_rec.id,
                ta_possible,
                ta_achieved,
                ta_percentage
            );
        END LOOP;
        
        IF total_possible > 0 THEN
                percentage_score_var := (achieved_score_var::DECIMAL / total_possible::DECIMAL) * 100.0;
        ELSE
            percentage_score_var := 0.0;
        END IF;
        
        performance_level_var := 'Proficient';
        
        UPDATE assessments
        SET total_possible_score = total_possible,
                achieved_score = achieved_score_var,
            percentage_score = percentage_score_var,
            performance_level = performance_level_var
        WHERE id = assessment_id_var;
    END LOOP;
END $$;

-- Create medium performing assessments (Competent)
DO $$
DECLARE
    facility_rec RECORD;
    assessment_type_rec RECORD;
    question_rec RECORD;
    thematic_area_rec RECORD;
    assessment_id_var INTEGER;
    total_possible INTEGER;
    achieved_score_var INTEGER;
    percentage_score_var DECIMAL(5,2);
    performance_level_var VARCHAR(50);
    ta_possible INTEGER;
    ta_achieved INTEGER;
    ta_percentage DECIMAL(5,2);
    response_val VARCHAR(10);
    points_earned_var INTEGER;
    rand_val INTEGER;
    health_worker_id_var INTEGER;
BEGIN
    -- Get second facility and middle assessment types
    SELECT id INTO facility_rec FROM facilities ORDER BY id OFFSET 1 LIMIT 1;
    
    -- Get a health worker for this facility
    SELECT id INTO health_worker_id_var 
    FROM health_workers 
    WHERE facility_id = facility_rec.id 
    ORDER BY RANDOM() 
    LIMIT 1;
    
    -- Skip if no health worker exists
    IF health_worker_id_var IS NULL THEN
        RAISE NOTICE 'No health workers found for facility %, skipping competent assessments', facility_rec.id;
        RETURN;
    END IF;
    
    FOR assessment_type_rec IN SELECT id, name FROM assessment_types ORDER BY id OFFSET 3 LIMIT 3 LOOP
        total_possible := 0;
            achieved_score_var := 0;
        
        INSERT INTO assessments (
            health_worker_id,
            facility_id,
            assessment_type_id,
            assessor_name,
            client_name,
            notes,
            total_possible_score,
            achieved_score,
            percentage_score,
            performance_level,
            created_at
        ) VALUES (
            health_worker_id_var,
            facility_rec.id,
            assessment_type_rec.id,
            'Nurse Average Provider',
            'Standard Client',
            'Competent level assessment',
            0,
            0,
            0.0,
            'Competent',
            CURRENT_TIMESTAMP - (random() * 20 || ' days')::INTERVAL
        ) RETURNING id INTO assessment_id_var;
        
        FOR thematic_area_rec IN 
            SELECT id, name 
            FROM thematic_areas 
            WHERE assessment_type_id = assessment_type_rec.id 
            ORDER BY display_order 
        LOOP
            ta_possible := 0;
            ta_achieved := 0;
            
            FOR question_rec IN 
                SELECT id, score_weight, is_critical
                FROM questions 
                WHERE thematic_area_id = thematic_area_rec.id 
                ORDER BY display_order 
            LOOP
                -- Competent: 75% Yes, 20% No, 5% NA
                rand_val := floor(random() * 100);
                IF rand_val < 5 THEN
                    response_val := 'NA';
                ELSIF rand_val < 25 THEN
                    response_val := 'No';
                ELSE
                    response_val := 'Yes';
                END IF;
                
                -- Only "Yes" responses get points (score_weight)
                IF response_val = 'Yes' THEN
                    points_earned_var := question_rec.score_weight;
                ELSE
                    points_earned_var := 0;
                END IF;
                
                -- Add to totals
                -- Only "Yes" and "No" count toward possible score
                -- "NA" questions don't count toward possible score at all
                IF response_val = 'Yes' THEN
                    total_possible := total_possible + question_rec.score_weight;
                        achieved_score_var := achieved_score_var + points_earned_var;
                    ta_possible := ta_possible + question_rec.score_weight;
                    ta_achieved := ta_achieved + points_earned_var;
                ELSIF response_val = 'No' THEN
                    -- "No" counts toward possible but not achieved
                    total_possible := total_possible + question_rec.score_weight;
                    ta_possible := ta_possible + question_rec.score_weight;
                END IF;
                -- "NA" responses don't contribute to scores at all
                
                INSERT INTO assessment_responses (
                    assessment_id,
                    question_id,
                    response,
                    points_earned
                ) VALUES (
                    assessment_id_var,
                    question_rec.id,
                    response_val,
                    points_earned_var
                );
            END LOOP;
            
            IF ta_possible > 0 THEN
                ta_percentage := (ta_achieved::DECIMAL / ta_possible::DECIMAL) * 100.0;
            ELSE
                ta_percentage := 0.0;
            END IF;
            
            INSERT INTO thematic_area_scores (
                assessment_id,
                thematic_area_id,
                possible_score,
                achieved_score,
                percentage_score
            ) VALUES (
                assessment_id_var,
                thematic_area_rec.id,
                ta_possible,
                ta_achieved,
                ta_percentage
            );
        END LOOP;
        
        IF total_possible > 0 THEN
            percentage_score_var := (achieved_score_var::DECIMAL / total_possible::DECIMAL) * 100.0;
        ELSE
            percentage_score_var := 0.0;
        END IF;
        
        IF percentage_score_var > 90 THEN
            performance_level_var := 'Proficient';
        ELSIF percentage_score_var >= 70 THEN
            performance_level_var := 'Competent';
        ELSE
            performance_level_var := 'Not Acceptable';
        END IF;
        
        UPDATE assessments
        SET total_possible_score = total_possible,
            achieved_score = achieved_score_var,
            percentage_score = percentage_score_var,
            performance_level = performance_level_var
        WHERE id = assessment_id_var;
    END LOOP;
END $$;

-- Create low performing assessments (Not Acceptable)
DO $$
DECLARE
    facility_rec RECORD;
    assessment_type_rec RECORD;
    question_rec RECORD;
    thematic_area_rec RECORD;
    assessment_id_var INTEGER;
    total_possible INTEGER;
    achieved_score_var INTEGER;
    percentage_score_var DECIMAL(5,2);
    performance_level_var VARCHAR(50);
    ta_possible INTEGER;
    ta_achieved INTEGER;
    ta_percentage DECIMAL(5,2);
    response_val VARCHAR(10);
    points_earned_var INTEGER;
    rand_val INTEGER;
    health_worker_id_var INTEGER;
BEGIN
    -- Get third facility and last assessment types
    SELECT id INTO facility_rec FROM facilities ORDER BY id OFFSET 2 LIMIT 1;
    
    -- Get a health worker for this facility
    SELECT id INTO health_worker_id_var 
    FROM health_workers 
    WHERE facility_id = facility_rec.id 
    ORDER BY RANDOM() 
    LIMIT 1;
    
    -- Skip if no health worker exists
    IF health_worker_id_var IS NULL THEN
        RAISE NOTICE 'No health workers found for facility %, skipping low performer assessments', facility_rec.id;
        RETURN;
    END IF;
    
    FOR assessment_type_rec IN SELECT id, name FROM assessment_types ORDER BY id OFFSET 6 LOOP
        total_possible := 0;
            achieved_score_var := 0;
        
        INSERT INTO assessments (
            health_worker_id,
            facility_id,
            assessment_type_id,
            assessor_name,
            client_name,
            notes,
            total_possible_score,
            achieved_score,
            percentage_score,
            performance_level,
            created_at
        ) VALUES (
            health_worker_id_var,
            facility_rec.id,
            assessment_type_rec.id,
            'Provider Needs Training',
            'Client Requiring Support',
            'Low performing assessment - needs improvement',
            0,
            0,
            0.0,
            'Not Acceptable',
            CURRENT_TIMESTAMP - (random() * 10 || ' days')::INTERVAL
        ) RETURNING id INTO assessment_id_var;
        
        FOR thematic_area_rec IN 
            SELECT id, name 
            FROM thematic_areas 
            WHERE assessment_type_id = assessment_type_rec.id 
            ORDER BY display_order 
        LOOP
            ta_possible := 0;
            ta_achieved := 0;
            
            FOR question_rec IN 
                SELECT id, score_weight, is_critical
                FROM questions 
                WHERE thematic_area_id = thematic_area_rec.id 
                ORDER BY display_order 
            LOOP
                -- Low performer: 40% Yes, 55% No, 5% NA
                rand_val := floor(random() * 100);
                IF rand_val < 5 THEN
                    response_val := 'NA';
                ELSIF rand_val < 60 THEN
                    response_val := 'No';
                ELSE
                    response_val := 'Yes';
                END IF;
                
                -- Only "Yes" responses get points (score_weight)
                IF response_val = 'Yes' THEN
                    points_earned_var := question_rec.score_weight;
                ELSE
                    points_earned_var := 0;
                END IF;
                
                -- Add to totals
                -- Only "Yes" and "No" count toward possible score
                -- "NA" questions don't count toward possible score at all
                IF response_val = 'Yes' THEN
                    total_possible := total_possible + question_rec.score_weight;
                        achieved_score_var := achieved_score_var + points_earned_var;
                    ta_possible := ta_possible + question_rec.score_weight;
                    ta_achieved := ta_achieved + points_earned_var;
                ELSIF response_val = 'No' THEN
                    -- "No" counts toward possible but not achieved
                    total_possible := total_possible + question_rec.score_weight;
                    ta_possible := ta_possible + question_rec.score_weight;
                END IF;
                -- "NA" responses don't contribute to scores at all
                
                INSERT INTO assessment_responses (
                    assessment_id,
                    question_id,
                    response,
                    points_earned
                ) VALUES (
                    assessment_id_var,
                    question_rec.id,
                    response_val,
                    points_earned_var
                );
            END LOOP;
            
            IF ta_possible > 0 THEN
                ta_percentage := (ta_achieved::DECIMAL / ta_possible::DECIMAL) * 100.0;
            ELSE
                ta_percentage := 0.0;
            END IF;
            
            INSERT INTO thematic_area_scores (
                assessment_id,
                thematic_area_id,
                possible_score,
                achieved_score,
                percentage_score
            ) VALUES (
                assessment_id_var,
                thematic_area_rec.id,
                ta_possible,
                ta_achieved,
                ta_percentage
            );
        END LOOP;
        
        IF total_possible > 0 THEN
                percentage_score_var := (achieved_score_var::DECIMAL / total_possible::DECIMAL) * 100.0;
        ELSE
            percentage_score_var := 0.0;
        END IF;
        
        performance_level_var := 'Not Acceptable';
        
        UPDATE assessments
        SET total_possible_score = total_possible,
                achieved_score = achieved_score_var,
            percentage_score = percentage_score_var,
            performance_level = performance_level_var
        WHERE id = assessment_id_var;
    END LOOP;
END $$;

-- Summary
SELECT 
    'Assessment seeding completed!' as status,
    COUNT(*) as total_assessments,
    COUNT(DISTINCT facility_id) as facilities_with_assessments,
    COUNT(DISTINCT assessment_type_id) as assessment_types_covered,
    COUNT(CASE WHEN performance_level = 'Proficient' THEN 1 END) as proficient_count,
    COUNT(CASE WHEN performance_level = 'Competent' THEN 1 END) as competent_count,
    COUNT(CASE WHEN performance_level = 'Not Acceptable' THEN 1 END) as not_acceptable_count
FROM assessments;

