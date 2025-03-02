DROP EXTENSION IF EXISTS pgcrypto;
CREATE EXTENSION pgcrypto SCHEMA public;

CREATE OR REPLACE FUNCTION generate_ulid_at_time(target_time TIMESTAMP WITH TIME ZONE)
    RETURNS UUID AS
$$
DECLARE
    timestamp_ms BIGINT;
    rand_bytes   BYTEA;
    result       UUID;
BEGIN
    timestamp_ms := FLOOR(EXTRACT(EPOCH FROM target_time) * 1000);
    rand_bytes := gen_random_bytes(10);
    result := CONCAT_WS('-',
                        LPAD(TO_HEX(timestamp_ms >> 16), 8, '0'),
                        LPAD(TO_HEX((timestamp_ms & x'FFFF'::int)), 4, '0'),
                        LPAD(TO_HEX(x'4000'::int | (get_byte(rand_bytes, 0) & x'0FFF'::int)), 4, '0'),
                        LPAD(TO_HEX(x'8000'::int | (get_byte(rand_bytes, 1) & x'3FFF'::int)), 4, '0'),
                        CONCAT(
                                LPAD(TO_HEX(get_byte(rand_bytes, 2)), 2, '0'),
                                LPAD(TO_HEX(get_byte(rand_bytes, 3)), 2, '0'),
                                LPAD(TO_HEX(get_byte(rand_bytes, 4)), 2, '0'),
                                LPAD(TO_HEX(get_byte(rand_bytes, 5)), 2, '0'),
                                LPAD(TO_HEX(get_byte(rand_bytes, 6)), 2, '0'),
                                LPAD(TO_HEX(get_byte(rand_bytes, 7)), 2, '0')
                        )
              )::UUID;
    RETURN result;
END;
$$ LANGUAGE plpgsql;

DO
$$
    DECLARE
        -- User IDs declarations
        student1_id UUID;
        student2_id UUID;
        student3_id UUID;
        mentor1_id  UUID;
        mentor2_id  UUID;
        admin_id    UUID;
    BEGIN
        -- Initialize UUIDs
        student1_id := '0194bb8e-1e7c-4082-806a-1c9483a59a1b';
        student2_id := '0194bbc5-0cfc-408e-802d-a854785c39d9';
        student3_id := '0194bbfb-fb7c-40d9-8029-7ca996ab5da6';
        mentor1_id := '0194bc32-e9fc-405c-801a-08ff3b0cf28b';
        mentor2_id := '0194bc69-d87c-40d5-809a-17513f3d2b98';
        admin_id := '0194bca0-c6fc-40ec-8051-68235bef9817';

        -- Users seeder
        INSERT INTO users (id, name, email, password_hash, role, has_avatar, created_at, updated_at)
        VALUES
            -- Students
            (student1_id, 'Student 1', 'student1@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '30 days',
             NOW() - INTERVAL '30 days'),

            (student2_id, 'Student 2', 'student2@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '1 hour',
             NOW() - INTERVAL '30 days' + INTERVAL '1 hour'),

            (student3_id, 'Student 3', 'student3@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '2 hours',
             NOW() - INTERVAL '30 days' + INTERVAL '2 hours'),

            -- Mentors
            (mentor1_id, 'Mentor 1', 'mentor1@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '3 hours',
             NOW() - INTERVAL '30 days' + INTERVAL '3 hours'),

            (mentor2_id, 'Mentor 2', 'mentor2@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '4 hours',
             NOW() - INTERVAL '30 days' + INTERVAL '4 hours'),

            -- Admin
            (admin_id, 'Admin 1', 'admin1@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'admin',
             FALSE,
             NOW() - INTERVAL '30 days' + INTERVAL '5 hours',
             NOW() - INTERVAL '30 days' + INTERVAL '5 hours');

        -- Students specific data
        INSERT INTO students (user_id, instance, major)
        VALUES (student1_id, 'University of Technology', 'Computer Science'),
               (student2_id, 'State University', 'Data Science'),
               (student3_id, 'National Institute', 'Software Engineering');

        -- Mentors specific data
        INSERT INTO mentors (user_id, specialization, experience, rating, rating_count, rating_total, price, balance)
        VALUES (mentor1_id, 'Web Development',
                'Over 8 years of experience in full-stack development with expertise in React, Node.js, and database design. Led development teams for enterprise applications and mentored junior developers.',
                4.8, 125, 4.8 * 125, 750000, 3500000),

               (mentor2_id, 'Machine Learning',
                'Data scientist with 6 years of experience in implementing machine learning models for production. Background in computer vision and natural language processing. Published research in top AI conferences.',
                4.7, 98, 4.7 * 98, 850000, 2750000);
    END
$$;
