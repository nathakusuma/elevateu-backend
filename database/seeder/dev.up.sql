DO
$$
    DECLARE
        -- User IDs declarations
        student1_id   UUID;
        student2_id   UUID;
        student3_id   UUID;
        mentor1_id    UUID;
        mentor2_id    UUID;
        admin_id      UUID;

        -- Additional User IDs (UUIDv7 format)
        student4_id   UUID := '01a4bb3c-d1fa-7b23-a8ec-15cdef123a01';
        student5_id   UUID := '01a4bb3c-d2fa-7b24-a8ed-15cdef123a02';
        student6_id   UUID := '01a4bb3c-d3fa-7b25-a8ee-15cdef123a03';
        student7_id   UUID := '01a4bb3c-d4fa-7b26-a8ef-15cdef123a04';
        student8_id   UUID := '01a4bb3c-d5fa-7b27-a8f0-15cdef123a05';
        student9_id   UUID := '01a4bb3c-d6fa-7b28-a8f1-15cdef123a06';
        student10_id  UUID := '01a4bb3c-d7fa-7b29-a8f2-15cdef123a07';
        student11_id  UUID := '01a4bb3c-d8fa-7b2a-a8f3-15cdef123a08';
        student12_id  UUID := '01a4bb3c-d9fa-7b2b-a8f4-15cdef123a09';
        student13_id  UUID := '01a4bb3c-dafa-7b2c-a8f5-15cdef123a10';
        mentor3_id    UUID := '01a4bb3c-dbfa-7b2d-a8f6-15cdef123b01';
        mentor4_id    UUID := '01a4bb3c-dcfa-7b2e-a8f7-15cdef123b02';
        mentor5_id    UUID := '01a4bb3c-ddfa-7b2f-a8f8-15cdef123b03';
        mentor6_id    UUID := '01a4bb3c-defa-7b30-a8f9-15cdef123b04';
        mentor7_id    UUID := '01a4bb3c-dffa-7b31-a8fa-15cdef123b05';
        mentor8_id    UUID := '01a4bb3c-e0fa-7b32-a8fb-15cdef123b06';
        mentor9_id    UUID := '01a4bb3c-e1fa-7b33-a8fc-15cdef123b07';
        mentor10_id   UUID := '01a4bb3c-e2fa-7b34-a8fd-15cdef123b08';
        mentor11_id   UUID := '01a4bb3c-e3fa-7b35-a8fe-15cdef123b09';
        mentor12_id   UUID := '01a4bb3c-e4fa-7b36-a8ff-15cdef123b10';
        admin2_id     UUID := '01a4bb3c-e5fa-7b37-a900-15cdef123c01';

        -- Category IDs declarations (UUIDv7 format)
        webdev_id        UUID;
        datascience_id   UUID;
        mobile_id        UUID;
        cloud_id         UUID;
        ai_id            UUID;
        devops_id        UUID;
        cybersecurity_id UUID;

        -- Course IDs declarations (UUIDv7 format)
        course1_id       UUID;
        course2_id       UUID;
        course3_id       UUID;
        course4_id       UUID;
        course5_id       UUID;
        course6_id       UUID;
        course7_id       UUID;
        course8_id       UUID;
        course9_id       UUID;
        course10_id      UUID;
        course11_id      UUID;
        course12_id      UUID;
        course13_id      UUID;
        course14_id      UUID;
        course15_id      UUID;
        course16_id      UUID;
        course17_id      UUID;
        course18_id      UUID;
        course19_id      UUID;
        course20_id      UUID;

        -- Course Video IDs
        video1_id     UUID := '01a4bb3c-f1fa-7c01-b001-15cdef234a01';
        video2_id     UUID := '01a4bb3c-f2fa-7c02-b002-15cdef234a02';
        video3_id     UUID := '01a4bb3c-f3fa-7c03-b003-15cdef234a03';
        video4_id     UUID := '01a4bb3c-f4fa-7c04-b004-15cdef234a04';
        video5_id     UUID := '01a4bb3c-f5fa-7c05-b005-15cdef234a05';
        video6_id     UUID := '01a4bb3c-f6fa-7c06-b006-15cdef234a06';
        video7_id     UUID := '01a4bb3c-f7fa-7c07-b007-15cdef234a07';
        video8_id     UUID := '01a4bb3c-f8fa-7c08-b008-15cdef234a08';
        video9_id     UUID := '01a4bb3c-f9fa-7c09-b009-15cdef234a09';
        video10_id    UUID := '01a4bb3c-fafa-7c0a-b00a-15cdef234a10';

        -- More videos for JS course
        video11_id    UUID := '01a4bb3c-fbfa-7c0b-b00b-15cdef234a11';
        video12_id    UUID := '01a4bb3c-fcfa-7c0c-b00c-15cdef234a12';
        video13_id    UUID := '01a4bb3c-fdfa-7c0d-b00d-15cdef234a13';
        video14_id    UUID := '01a4bb3c-fefa-7c0e-b00e-15cdef234a14';
        video15_id    UUID := '01a4bb3c-fffa-7c0f-b00f-15cdef234a15';

        -- Videos for React course
        video16_id    UUID := '01a4bb3d-00fa-7c10-b010-15cdef234a16';
        video17_id    UUID := '01a4bb3d-01fa-7c11-b011-15cdef234a17';
        video18_id    UUID := '01a4bb3d-02fa-7c12-b012-15cdef234a18';
        video19_id    UUID := '01a4bb3d-03fa-7c13-b013-15cdef234a19';
        video20_id    UUID := '01a4bb3d-04fa-7c14-b014-15cdef234a20';

        -- Course Material IDs
        material1_id  UUID := '01a4bb3d-10fa-7d01-c001-15cdef345a01';
        material2_id  UUID := '01a4bb3d-11fa-7d02-c002-15cdef345a02';
        material3_id  UUID := '01a4bb3d-12fa-7d03-c003-15cdef345a03';
        material4_id  UUID := '01a4bb3d-13fa-7d04-c004-15cdef345a04';
        material5_id  UUID := '01a4bb3d-14fa-7d05-c005-15cdef345a05';
        material6_id  UUID := '01a4bb3d-15fa-7d06-c006-15cdef345a06';
        material7_id  UUID := '01a4bb3d-16fa-7d07-c007-15cdef345a07';
        material8_id  UUID := '01a4bb3d-17fa-7d08-c008-15cdef345a08';
        material9_id  UUID := '01a4bb3d-18fa-7d09-c009-15cdef345a09';
        material10_id UUID := '01a4bb3d-19fa-7d0a-c00a-15cdef345a10';

        -- More materials for JS course
        material11_id UUID := '01a4bb3d-1afa-7d0b-c00b-15cdef345a11';
        material12_id UUID := '01a4bb3d-1bfa-7d0c-c00c-15cdef345a12';
        material13_id UUID := '01a4bb3d-1cfa-7d0d-c00d-15cdef345a13';
        material14_id UUID := '01a4bb3d-1dfa-7d0e-c00e-15cdef345a14';
        material15_id UUID := '01a4bb3d-1efa-7d0f-c00f-15cdef345a15';

        -- Materials for React course
        material16_id UUID := '01a4bb3d-1ffa-7d10-c010-15cdef345a16';
        material17_id UUID := '01a4bb3d-20fa-7d11-c011-15cdef345a17';
        material18_id UUID := '01a4bb3d-21fa-7d12-c012-15cdef345a18';
        material19_id UUID := '01a4bb3d-22fa-7d13-c013-15cdef345a19';
        material20_id UUID := '01a4bb3d-23fa-7d14-c014-15cdef345a20';
    BEGIN
        -- Initialize User UUIDs
        student1_id := '0194bb8e-1e7c-4082-806a-1c9483a59a1b';
        student2_id := '0194bbc5-0cfc-408e-802d-a854785c39d9';
        student3_id := '0194bbfb-fb7c-40d9-8029-7ca996ab5da6';
        mentor1_id := '0194bc32-e9fc-405c-801a-08ff3b0cf28b';
        mentor2_id := '0194bc69-d87c-40d5-809a-17513f3d2b98';
        admin_id := '0194bca0-c6fc-40ec-8051-68235bef9817';

        -- Initialize Category UUIDs (UUIDv7 format)
        webdev_id := '01890e5c-b334-7c38-8c8c-b05c57968489';
        datascience_id := '01890e5c-bed5-7ebb-8dcf-ada8c40c2253';
        mobile_id := '01890e5c-c5b6-7da1-b35c-70d83a37d078';
        cloud_id := '01890e5c-cdbe-718c-9623-ec4a1bdfaa27';
        ai_id := '01890e5c-d583-752e-ab1e-24d75a5c4dfc';
        devops_id := '01890e5c-dc95-7de2-9e36-00a6a28bf8de';
        cybersecurity_id := '01890e5c-e31e-7b65-b6fb-08a3b7a00d56';

        -- Initialize Course UUIDs (UUIDv7 format)
        course1_id := '01890e5d-05a6-7b0b-a7b3-3f8ed37b5261';
        course2_id := '01890e5d-0ec9-7e29-ae5f-4a7c66f5d456';
        course3_id := '01890e5d-1657-71ae-bcba-54a2a55b1bb8';
        course4_id := '01890e5d-1e20-7e5c-b0b2-6dbc90efac62';
        course5_id := '01890e5d-25de-7d73-bf1c-8c1a13f17839';
        course6_id := '01890e5d-2dd5-7bde-994a-76e4b9e18a08';
        course7_id := '01890e5d-3530-7487-8b77-fc5cc0773b3f';
        course8_id := '01890e5d-3c4a-72fe-b78b-e1ebf0cbed87';
        course9_id := '01890e5d-437c-72ca-8ddc-4aa6b26ddc11';
        course10_id := '01890e5d-4a4c-7452-b55a-7b8fa75d84e5';
        course11_id := '01890e5d-516c-7dc0-9d9c-1d6ec6e36fc4';
        course12_id := '01890e5d-58ca-787d-9ed5-31fbe2bc9429';
        course13_id := '01890e5d-6086-7c90-a81b-9dccca9a6c97';
        course14_id := '01890e5d-6769-7b84-baba-4b3aad93a8b0';
        course15_id := '01890e5d-6eda-71bf-8a5e-a74a7f8c8d2c';
        course16_id := '01890e5d-76d1-75af-9c61-0adf9a0d8db1';
        course17_id := '01890e5d-7e5d-7bfb-9e26-2eb534b456f3';
        course18_id := '01890e5d-8540-7a65-9b7d-5d95efed8bc9';
        course19_id := '01890e5d-8c74-7b15-9ed8-15b5e14ab0c7';
        course20_id := '01890e5d-946d-7a56-a9b4-69d7c6fefb92';

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

        -- Additional users
        INSERT INTO users (id, name, email, password_hash, role, has_avatar, created_at, updated_at)
        VALUES
            -- Additional Students
            (student4_id, 'Student 4', 'student4@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '28 days',
             NOW() - INTERVAL '28 days'),

            (student5_id, 'Student 5', 'student5@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '26 days',
             NOW() - INTERVAL '26 days'),

            (student6_id, 'Student 6', 'student6@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '24 days',
             NOW() - INTERVAL '24 days'),

            (student7_id, 'Student 7', 'student7@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '22 days',
             NOW() - INTERVAL '22 days'),

            (student8_id, 'Student 8', 'student8@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '20 days',
             NOW() - INTERVAL '20 days'),

            (student9_id, 'Student 9', 'student9@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '18 days',
             NOW() - INTERVAL '18 days'),

            (student10_id, 'Student 10', 'student10@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             FALSE,
             NOW() - INTERVAL '16 days',
             NOW() - INTERVAL '16 days'),

            (student11_id, 'Student 11', 'student11@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             TRUE,
             NOW() - INTERVAL '14 days',
             NOW() - INTERVAL '14 days'),

            (student12_id, 'Student 12', 'student12@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             TRUE,
             NOW() - INTERVAL '12 days',
             NOW() - INTERVAL '12 days'),

            (student13_id, 'Student 13', 'student13@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'student',
             TRUE,
             NOW() - INTERVAL '10 days',
             NOW() - INTERVAL '10 days'),

            -- Additional Mentors
            (mentor3_id, 'Mentor 3', 'mentor3@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             FALSE,
             NOW() - INTERVAL '28 days',
             NOW() - INTERVAL '28 days'),

            (mentor4_id, 'Mentor 4', 'mentor4@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             FALSE,
             NOW() - INTERVAL '26 days',
             NOW() - INTERVAL '26 days'),

            (mentor5_id, 'Mentor 5', 'mentor5@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             FALSE,
             NOW() - INTERVAL '24 days',
             NOW() - INTERVAL '24 days'),

            (mentor6_id, 'Mentor 6', 'mentor6@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             FALSE,
             NOW() - INTERVAL '22 days',
             NOW() - INTERVAL '22 days'),

            (mentor7_id, 'Mentor 7', 'mentor7@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             FALSE,
             NOW() - INTERVAL '20 days',
             NOW() - INTERVAL '20 days'),

            (mentor8_id, 'Mentor 8', 'mentor8@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             TRUE,
             NOW() - INTERVAL '18 days',
             NOW() - INTERVAL '18 days'),

            (mentor9_id, 'Mentor 9', 'mentor9@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             TRUE,
             NOW() - INTERVAL '16 days',
             NOW() - INTERVAL '16 days'),

            (mentor10_id, 'Mentor 10', 'mentor10@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             TRUE,
             NOW() - INTERVAL '14 days',
             NOW() - INTERVAL '14 days'),

            (mentor11_id, 'Mentor 11', 'mentor11@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             TRUE,
             NOW() - INTERVAL '12 days',
             NOW() - INTERVAL '12 days'),

            (mentor12_id, 'Mentor 12', 'mentor12@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'mentor',
             TRUE,
             NOW() - INTERVAL '10 days',
             NOW() - INTERVAL '10 days'),

            -- Additional Admin
            (admin2_id, 'Admin 2', 'admin2@seeder.nathakusuma.com',
             '$2a$12$Hr/bomKRJ7QxBC2OLiZecO8clgrSxU6xJc5bkyZBcONniUDgr8E/6', 'admin',
             TRUE,
             NOW() - INTERVAL '25 days',
             NOW() - INTERVAL '25 days');

        -- Students specific data
        INSERT INTO students (user_id, instance, major, point, subscribed_boost_until, subscribed_challenge_until)
        VALUES (student1_id, 'University of Technology', 'Computer Science', 0, NULL, NULL),
               (student2_id, 'State University', 'Data Science', 0, '2026-12-31 23:59:59', '2026-12-31 23:59:59'),
               (student3_id, 'National Institute', 'Software Engineering', 0, NULL, NULL);

        -- Additional Students specific data
        INSERT INTO students (user_id, instance, major, point, subscribed_boost_until, subscribed_challenge_until)
        VALUES (student4_id, 'Technical Institute', 'Web Development', 150, '2025-06-30 23:59:59', NULL),
               (student5_id, 'University of Engineering', 'Information Systems', 75, NULL, '2025-08-31 23:59:59'),
               (student6_id, 'Tech Academy', 'Artificial Intelligence', 230, '2025-10-15 23:59:59',
                '2025-10-15 23:59:59'),
               (student7_id, 'Computer Science College', 'Machine Learning', 50, NULL, NULL),
               (student8_id, 'IT University', 'Cloud Computing', 175, '2025-07-20 23:59:59', NULL),
               (student9_id, 'State University', 'Data Analytics', 300, NULL, '2025-09-10 23:59:59'),
               (student10_id, 'Technical College', 'Mobile Applications', 125, NULL, NULL),
               (student11_id, 'National Institute', 'Cybersecurity', 275, '2025-12-01 23:59:59', '2025-12-01 23:59:59'),
               (student12_id, 'Computer Engineering School', 'DevOps', 90, NULL, NULL),
               (student13_id, 'Digital University', 'Full Stack Development', 200, '2025-11-15 23:59:59',
                '2025-11-15 23:59:59');

        -- Mentors specific data
        INSERT INTO mentors (user_id, specialization, experience, rating, rating_count, rating_total, price, balance)
        VALUES (mentor1_id, 'Web Development',
                'Over 8 years of experience in full-stack development with expertise in React, Node.js, and database design. Led development teams for enterprise applications and mentored junior developers.',
                4.8, 125, 4.8 * 125, 750000, 3500000),

               (mentor2_id, 'Machine Learning',
                'Data scientist with 6 years of experience in implementing machine learning models for production. Background in computer vision and natural language processing. Published research in top AI conferences.',
                4.7, 98, 4.7 * 98, 850000, 2750000);

        -- Additional Mentors specific data
        INSERT INTO mentors (user_id, specialization, experience, rating, rating_count, rating_total, price, balance)
        VALUES (mentor3_id, 'Mobile App Development',
                'iOS and Android developer with 7 years of experience. Created over 20 published apps with millions of downloads. Specialized in React Native and Swift.',
                4.6, 85, 4.6 * 85, 700000, 2900000),

               (mentor4_id, 'Data Engineering',
                'Big data specialist with 5 years of experience designing data pipelines and ETL processes. Expert in Hadoop, Spark, and cloud-based data solutions.',
                4.8, 62, 4.8 * 62, 900000, 3200000),

               (mentor5_id, 'Cloud Architecture',
                'AWS certified solutions architect with 9 years of experience. Helped dozens of companies migrate to cloud infrastructure and optimize costs.',
                4.9, 120, 4.9 * 120, 950000, 5700000),

               (mentor6_id, 'Frontend Development',
                'UI/UX focused developer with 6 years of experience. Expert in React, Vue, and modern CSS frameworks. Previously worked at top tech companies.',
                4.7, 95, 4.7 * 95, 650000, 3100000),

               (mentor7_id, 'Backend Development',
                'Backend engineer with 8 years of experience building scalable APIs and microservices. Proficient in Node.js, Go, and distributed systems.',
                4.8, 105, 4.8 * 105, 800000, 4200000),

               (mentor8_id, 'DevOps Engineering',
                'DevOps specialist with 7 years of experience implementing CI/CD pipelines and infrastructure as code. Expert in Docker, Kubernetes, and GitOps.',
                4.7, 88, 4.7 * 88, 850000, 3750000),

               (mentor9_id, 'Cybersecurity',
                'Information security expert with 10 years of experience. Certified ethical hacker with background in penetration testing and security architecture.',
                4.9, 75, 4.9 * 75, 1000000, 3800000),

               (mentor10_id, 'Database Design',
                'Database architect with 8 years of experience. Specializes in SQL and NoSQL database optimization, data modeling, and performance tuning.',
                4.6, 90, 4.6 * 90, 750000, 3400000),

               (mentor11_id, 'Blockchain Development',
                'Blockchain developer with 5 years of experience building dApps and smart contracts. Expert in Ethereum, Solidity, and web3 technologies.',
                4.7, 60, 4.7 * 60, 950000, 2850000),

               (mentor12_id, 'Game Development',
                'Game developer with 9 years of experience. Worked on multiple published titles using Unity and Unreal Engine. Expert in C# and C++ for games.',
                4.8, 110, 4.8 * 110, 800000, 4400000);

        -- Categories seeder
        INSERT INTO categories (id, name)
        VALUES (webdev_id, 'Web Development'),
               (datascience_id, 'Data Science'),
               (mobile_id, 'Mobile Development'),
               (cloud_id, 'Cloud Computing'),
               (ai_id, 'Artificial Intelligence'),
               (devops_id, 'DevOps'),
               (cybersecurity_id, 'Cybersecurity');

        -- Courses seeder
        INSERT INTO courses (id, category_id, title, description, teacher_name, rating, rating_count, total_rating,
                             enrollment_count, content_count, created_at, updated_at)
        VALUES
            -- Web Development courses
            (course1_id,
             webdev_id,
             'Modern JavaScript Fundamentals',
             'Master the core concepts of JavaScript including ES6+ features, asynchronous programming, and functional programming techniques. This comprehensive course covers everything from basic syntax to advanced topics like closures, promises, and modern tooling.',
             'Dr. Sarah Johnson',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '60 days',
             NOW() - INTERVAL '15 days'),

            (course2_id,
             webdev_id,
             'Full-Stack React & Node.js',
             'Build production-ready applications with React and Node.js. Learn state management with Redux, server-side rendering, RESTful API design, authentication, testing, and deployment strategies for modern web applications.',
             'Michael Chen',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '45 days',
             NOW() - INTERVAL '10 days'),

            (course3_id,
             webdev_id,
             'CSS Mastery: Advanced Styling',
             'Take your CSS skills to the next level with advanced selectors, animations, grid layouts, and CSS variables. Learn responsive design techniques, CSS architecture patterns, and optimization strategies for modern websites.',
             'Sophia Martinez',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '85 days',
             NOW() - INTERVAL '12 days'),

            (course4_id,
             webdev_id,
             'Vue.js for Frontend Development',
             'Learn Vue.js from scratch and build dynamic user interfaces. This course covers Vue components, directives, Vuex for state management, Vue Router for navigation, and integration with backend APIs.',
             'David Wilson',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '70 days',
             NOW() - INTERVAL '8 days'),

            -- Data Science courses
            (course5_id,
             datascience_id,
             'Python for Data Analysis',
             'Learn how to use Python libraries like Pandas, NumPy, and Matplotlib to analyze, visualize, and interpret complex datasets. This course covers data cleaning, transformation, exploratory analysis, and creating insightful visualizations.',
             'Dr. Emily Rodriguez',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '90 days',
             NOW() - INTERVAL '25 days'),

            (course6_id,
             datascience_id,
             'Statistical Methods for Data Science',
             'Develop a strong foundation in statistics essential for data science. Topics include probability distributions, hypothesis testing, regression analysis, experimental design, and Bayesian methods with practical applications using Python and R.',
             'Prof. James Wilson',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '75 days',
             NOW() - INTERVAL '20 days'),

            (course7_id,
             datascience_id,
             'Big Data Processing with Spark',
             'Master distributed data processing using Apache Spark. Learn to handle large-scale datasets, perform batch and stream processing, and implement machine learning pipelines in a distributed environment.',
             'Dr. Aisha Patel',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '65 days',
             NOW() - INTERVAL '18 days'),

            (course8_id,
             datascience_id,
             'Data Visualization with D3.js',
             'Create interactive and dynamic data visualizations for the web using D3.js. Learn how to transform data into compelling visual stories through charts, graphs, and interactive dashboards.',
             'Marcus Johnson',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '55 days',
             NOW() - INTERVAL '14 days'),

            -- Mobile Development courses
            (course9_id,
             mobile_id,
             'iOS App Development with Swift',
             'Create beautiful and functional iOS applications using Swift and UIKit. Learn app architecture, interface design principles, Core Data, networking, and publishing to the App Store. Build several complete applications throughout the course.',
             'Alex Turner',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '120 days',
             NOW() - INTERVAL '30 days'),

            (course10_id,
             mobile_id,
             'Android Development with Kotlin',
             'Build modern Android applications using Kotlin and Jetpack components. Master Android architecture, UI design with Material Design, data persistence, background processing, and Google Play deployment.',
             'Priya Sharma',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '110 days',
             NOW() - INTERVAL '28 days'),

            (course11_id,
             mobile_id,
             'Cross-Platform Apps with Flutter',
             'Develop cross-platform mobile applications with a single codebase using Flutter and Dart. Learn to create beautiful UI components, manage state, integrate APIs, and deploy to both iOS and Android platforms.',
             'Ryan Park',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '100 days',
             NOW() - INTERVAL '26 days'),

            -- Cloud Computing courses
            (course12_id,
             cloud_id,
             'AWS Solutions Architect',
             'Prepare for the AWS Solutions Architect certification while learning to design distributed systems on AWS. Cover EC2, S3, RDS, Lambda, API Gateway, CloudFormation, and best practices for security, cost optimization, and high availability.',
             'Jennifer Parker',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '100 days',
             NOW() - INTERVAL '5 days'),

            (course13_id,
             cloud_id,
             'Google Cloud Platform Essentials',
             'Learn the fundamentals of Google Cloud Platform services and architecture. Explore Compute Engine, App Engine, Cloud Storage, BigQuery, Cloud Functions, and deployment automation with hands-on projects.',
             'Thomas Lee',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '95 days',
             NOW() - INTERVAL '22 days'),

            (course14_id,
             cloud_id,
             'Microsoft Azure Administration',
             'Master Azure infrastructure management including virtual machines, storage, networking, security, and monitoring. Prepare for the Azure Administrator certification with practical exercises and real-world scenarios.',
             'Samantha Wright',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '85 days',
             NOW() - INTERVAL '16 days'),

            (course15_id,
             cloud_id,
             'Serverless Architecture Patterns',
             'Design and implement scalable serverless applications across major cloud providers. Learn event-driven architecture, function-as-a-service, managed services integration, and deployment automation for serverless systems.',
             'Daniel Garcia',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '75 days',
             NOW() - INTERVAL '10 days'),

            -- AI courses
            (course16_id,
             ai_id,
             'Machine Learning Fundamentals',
             'Build a strong foundation in machine learning algorithms and techniques. This course covers supervised and unsupervised learning, model evaluation, feature engineering, and practical implementation using Python''s scikit-learn library.',
             'Dr. Robert Kim',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '85 days',
             NOW() - INTERVAL '12 days'),

            (course17_id,
             ai_id,
             'Deep Learning with TensorFlow',
             'Dive into neural networks and deep learning using TensorFlow and Keras. Learn to build, train and deploy convolutional neural networks, recurrent neural networks, and transformers for computer vision, NLP, and other applications.',
             'Dr. Lisa Zhang',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '70 days',
             NOW() - INTERVAL '8 days'),

            (course18_id,
             ai_id,
             'Natural Language Processing',
             'Master techniques for processing and analyzing text data using Python. Learn text preprocessing, feature extraction, sentiment analysis, topic modeling, and building language models with transformers like BERT and GPT.',
             'Prof. Alan Foster',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '60 days',
             NOW() - INTERVAL '15 days'),

            (course19_id,
             ai_id,
             'Computer Vision Applications',
             'Learn to build applications that can interpret and process visual information from images and videos. Cover image classification, object detection, facial recognition, and generative models using PyTorch and OpenCV.',
             'Dr. Nina Zhou',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '50 days',
             NOW() - INTERVAL '7 days'),

            -- DevOps courses
            (course20_id,
             devops_id,
             'Docker and Kubernetes for DevOps',
             'Master containerization with Docker and orchestration with Kubernetes. Learn to build, deploy, and scale microservices architectures using container technologies with CI/CD pipeline integration.',
             'Richard Stevens',
             0,
             0,
             0,
             0,
             0,
             NOW() - INTERVAL '90 days',
             NOW() - INTERVAL '14 days');

        -- Course Videos seeder
        -- Get actual course IDs from the database
        -- JavaScript course
        SELECT id INTO course1_id FROM courses WHERE title = 'Modern JavaScript Fundamentals';
        -- React course
        SELECT id INTO course2_id FROM courses WHERE title = 'Full-Stack React & Node.js';
        -- Python course
        SELECT id INTO course5_id FROM courses WHERE title = 'Python for Data Analysis';
        -- ML course
        SELECT id INTO course16_id FROM courses WHERE title = 'Machine Learning Fundamentals';
        -- iOS course
        SELECT id INTO course9_id FROM courses WHERE title = 'iOS App Development with Swift';
        -- Android course
        SELECT id INTO course10_id FROM courses WHERE title = 'Android Development with Kotlin';
        -- AWS course
        SELECT id INTO course12_id FROM courses WHERE title = 'AWS Solutions Architect';
        -- Another cloud course
        SELECT id INTO course13_id FROM courses WHERE title = 'Google Cloud Platform Essentials';
        -- Docker course
        SELECT id INTO course20_id FROM courses WHERE title = 'Docker and Kubernetes for DevOps';

        -- Check if we successfully got the course IDs
        IF course1_id IS NOT NULL THEN
            -- Course Videos seeder for JavaScript course
            INSERT INTO course_videos (id, course_id, title, description, duration, is_free, "order", created_at,
                                       updated_at)
            VALUES (video1_id, course1_id, 'Introduction to Modern JavaScript',
                    'Overview of JavaScript''s evolution and the modern features we''ll cover in this course. Learn about the course structure and what to expect.',
                    1200, TRUE, 1,
                    NOW() - INTERVAL '60 days',
                    NOW() - INTERVAL '60 days'),

                   (video2_id, course1_id, 'ES6+ Variables and Scoping',
                    'Deep dive into let, const, and var declarations. Understanding block scope vs function scope and temporal dead zone concepts.',
                    1800, TRUE, 2,
                    NOW() - INTERVAL '59 days',
                    NOW() - INTERVAL '59 days'),

                   (video3_id, course1_id, 'Arrow Functions and Lexical This',
                    'Master arrow function syntax and understand how they handle the "this" keyword differently from traditional functions.',
                    2100, FALSE, 3,
                    NOW() - INTERVAL '58 days',
                    NOW() - INTERVAL '58 days'),

                   (video4_id, course1_id, 'Destructuring Objects and Arrays',
                    'Learn efficient syntax for extracting values from complex data structures using ES6 destructuring patterns.',
                    1950, FALSE, 4,
                    NOW() - INTERVAL '57 days',
                    NOW() - INTERVAL '57 days'),

                   (video5_id, course1_id, 'Spread and Rest Operators',
                    'Understand the powerful spread and rest operators for working with arrays, objects, and function parameters.',
                    1650, FALSE, 5,
                    NOW() - INTERVAL '56 days',
                    NOW() - INTERVAL '56 days'),

                   (video11_id, course1_id, 'Template Literals and String Methods',
                    'Explore template literals for expressive string formatting and the new string methods introduced in ES6+.',
                    1750, FALSE, 6,
                    NOW() - INTERVAL '55 days',
                    NOW() - INTERVAL '55 days'),

                   (video12_id, course1_id, 'Asynchronous JavaScript: Promises',
                    'Learn how to use Promises to write cleaner asynchronous code and handle errors more effectively.',
                    2400, FALSE, 7,
                    NOW() - INTERVAL '54 days',
                    NOW() - INTERVAL '54 days'),

                   (video13_id, course1_id, 'Async/Await Syntax',
                    'Master the async/await pattern for writing asynchronous code that looks and behaves like synchronous code.',
                    2200, FALSE, 8,
                    NOW() - INTERVAL '53 days',
                    NOW() - INTERVAL '53 days'),

                   (video14_id, course1_id, 'JavaScript Modules and Import/Export',
                    'Understand the module system in JavaScript and how to organize code using imports and exports.',
                    1900, FALSE, 9,
                    NOW() - INTERVAL '52 days',
                    NOW() - INTERVAL '52 days'),

                   (video15_id, course1_id, 'Modern JavaScript Tooling',
                    'Overview of essential development tools like Babel, Webpack, ESLint, and npm for modern JavaScript projects.',
                    2300, FALSE, 10,
                    NOW() - INTERVAL '51 days',
                    NOW() - INTERVAL '51 days');

            -- Course Materials seeder for JavaScript course
            INSERT INTO course_materials (id, course_id, title, subtitle, is_free, "order", created_at, updated_at)
            VALUES (material1_id, course1_id, 'JS Fundamentals Cheatsheet',
                    'Core JavaScript concepts reference',
                    TRUE, 1,
                    NOW() - INTERVAL '60 days',
                    NOW() - INTERVAL '60 days'),

                   (material2_id, course1_id, 'ES6+ Features Reference',
                    'Guide to modern JavaScript features',
                    TRUE, 2,
                    NOW() - INTERVAL '59 days',
                    NOW() - INTERVAL '59 days'),

                   (material3_id, course1_id, 'Async JavaScript Patterns',
                    'Explanation of async patterns and practices',
                    FALSE, 3,
                    NOW() - INTERVAL '58 days',
                    NOW() - INTERVAL '58 days'),

                   (material4_id, course1_id, 'Functional Programming in JS',
                    'Guide to functional programming concepts',
                    FALSE, 4,
                    NOW() - INTERVAL '57 days',
                    NOW() - INTERVAL '57 days'),

                   (material5_id, course1_id, 'JS Performance Optimization',
                    'Techniques for writing efficient code',
                    FALSE, 5,
                    NOW() - INTERVAL '56 days',
                    NOW() - INTERVAL '56 days'),

                   (material11_id, course1_id, 'Modern JS Coding Standards',
                    'Best practices for clean code',
                    FALSE, 6,
                    NOW() - INTERVAL '55 days',
                    NOW() - INTERVAL '55 days'),

                   (material12_id, course1_id, 'Debugging JS Applications',
                    'Advanced debugging techniques and tools',
                    FALSE, 7,
                    NOW() - INTERVAL '54 days',
                    NOW() - INTERVAL '54 days'),

                   (material13_id, course1_id, 'JavaScript Design Patterns',
                    'Common design patterns in JavaScript',
                    FALSE, 8,
                    NOW() - INTERVAL '53 days',
                    NOW() - INTERVAL '53 days'),

                   (material14_id, course1_id, 'JavaScript Testing Strategies',
                    'Guide to unit, integration, and E2E testing',
                    FALSE, 9,
                    NOW() - INTERVAL '52 days',
                    NOW() - INTERVAL '52 days'),

                   (material15_id, course1_id, 'JavaScript Project Structure',
                    'Best practices for organizing JS applications',
                    FALSE, 10,
                    NOW() - INTERVAL '51 days',
                    NOW() - INTERVAL '51 days');

            -- Update content_count for JavaScript course
            UPDATE courses
            SET content_count = 20
            WHERE id = course1_id;
        END IF;

        -- Check if we successfully got the React course ID
        IF course2_id IS NOT NULL THEN
            -- Course Videos seeder for React course
            INSERT INTO course_videos (id, course_id, title, description, duration, is_free, "order", created_at,
                                       updated_at)
            VALUES (video16_id, course2_id, 'Introduction to Full-Stack Development',
                    'Overview of full-stack development with React and Node.js. Learn about the course structure and modern web application architecture.',
                    1500, TRUE, 1,
                    NOW() - INTERVAL '45 days',
                    NOW() - INTERVAL '45 days'),

                   (video17_id, course2_id, 'Setting Up a React Application',
                    'Learn how to create a React application from scratch using Create React App and understand the generated project structure.',
                    1800, TRUE, 2,
                    NOW() - INTERVAL '44 days',
                    NOW() - INTERVAL '44 days'),

                   (video18_id, course2_id, 'React Components and Props',
                    'Master the concept of components in React and understand how to use props to pass data between components.',
                    2000, FALSE, 3,
                    NOW() - INTERVAL '43 days',
                    NOW() - INTERVAL '43 days'),

                   (video19_id, course2_id, 'State Management in React',
                    'Learn about state in React, how to use useState hook, and understand when and how to lift state up in your component hierarchy.',
                    2200, FALSE, 4,
                    NOW() - INTERVAL '42 days',
                    NOW() - INTERVAL '42 days'),

                   (video20_id, course2_id, 'Introduction to Node.js and Express',
                    'Get started with Node.js backend development using Express framework. Learn about routing, middleware, and API development.',
                    1900, FALSE, 5,
                    NOW() - INTERVAL '41 days',
                    NOW() - INTERVAL '41 days');

            -- Course Materials seeder for React course
            INSERT INTO course_materials (id, course_id, title, subtitle, is_free, "order", created_at, updated_at)
            VALUES (material16_id, course2_id, 'React Component Lifecycle',
                    'Understanding lifecycle methods and hooks',
                    TRUE, 1,
                    NOW() - INTERVAL '45 days',
                    NOW() - INTERVAL '45 days'),

                   (material17_id, course2_id, 'State Management with Redux',
                    'Guide to Redux architecture and patterns',
                    FALSE, 2,
                    NOW() - INTERVAL '44 days',
                    NOW() - INTERVAL '44 days'),

                   (material18_id, course2_id, 'RESTful API Design with Node.js',
                    'Best practices for designing robust APIs',
                    FALSE, 3,
                    NOW() - INTERVAL '43 days',
                    NOW() - INTERVAL '43 days'),

                   (material19_id, course2_id, 'Authentication and Authorization',
                    'Implementing secure auth flows in applications',
                    FALSE, 4,
                    NOW() - INTERVAL '42 days',
                    NOW() - INTERVAL '42 days'),

                   (material20_id, course2_id, 'Deployment for Web Apps',
                    'Guide to deploying full-stack applications',
                    FALSE, 5,
                    NOW() - INTERVAL '41 days',
                    NOW() - INTERVAL '41 days');

            -- Update content_count for React course
            UPDATE courses
            SET content_count = 10
            WHERE id = course2_id;
        END IF;

        -- Course Enrollments seeder
        -- Insert enrollments only for courses that exist
        IF course1_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES
                -- Student 1 enrollments
                (course1_id, student1_id, 10, TRUE,
                 NOW() - INTERVAL '28 days',
                 NOW() - INTERVAL '2 days');
        END IF;

        IF course2_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course2_id, student1_id, 8, FALSE,
                    NOW() - INTERVAL '25 days',
                    NOW() - INTERVAL '1 day');
        END IF;

        IF course5_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course5_id, student1_id, 4, FALSE,
                    NOW() - INTERVAL '20 days',
                    NOW() - INTERVAL '3 days');
        END IF;

        -- Student 2 enrollments
        IF course1_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course1_id, student2_id, 7, FALSE,
                    NOW() - INTERVAL '26 days',
                    NOW() - INTERVAL '5 days');
        END IF;

        IF course16_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course16_id, student2_id, 12, TRUE,
                    NOW() - INTERVAL '24 days',
                    NOW() - INTERVAL '10 days');
        END IF;

        -- Student 3 enrollments
        IF course9_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course9_id, student3_id, 5, FALSE,
                    NOW() - INTERVAL '22 days',
                    NOW() - INTERVAL '2 days');
        END IF;

        IF course10_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course10_id, student3_id, 8, FALSE,
                    NOW() - INTERVAL '20 days',
                    NOW() - INTERVAL '1 day');
        END IF;

        -- Student 4 enrollments
        IF course1_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course1_id, student4_id, 15, FALSE,
                    NOW() - INTERVAL '27 days',
                    NOW() - INTERVAL '1 day');
        END IF;

        IF course2_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course2_id, student4_id, 10, TRUE,
                    NOW() - INTERVAL '25 days',
                    NOW() - INTERVAL '5 days');
        END IF;

        -- Student 5 enrollments
        IF course5_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course5_id, student5_id, 8, FALSE,
                    NOW() - INTERVAL '24 days',
                    NOW() - INTERVAL '2 days');
        END IF;

        -- Student 6 enrollments
        IF course16_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course16_id, student6_id, 7, FALSE,
                    NOW() - INTERVAL '23 days',
                    NOW() - INTERVAL '1 day');
        END IF;

        -- Student 7 enrollments
        IF course12_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course12_id, student7_id, 5, FALSE,
                    NOW() - INTERVAL '21 days',
                    NOW() - INTERVAL '3 days');
        END IF;

        IF course13_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course13_id, student7_id, 8, FALSE,
                    NOW() - INTERVAL '19 days',
                    NOW() - INTERVAL '2 days');
        END IF;

        -- Student 8 enrollments
        IF course12_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course12_id, student8_id, 10, TRUE,
                    NOW() - INTERVAL '19 days',
                    NOW() - INTERVAL '5 days');
        END IF;

        -- Student 9 enrollments
        IF course5_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course5_id, student9_id, 9, FALSE,
                    NOW() - INTERVAL '17 days',
                    NOW() - INTERVAL '3 days');
        END IF;

        -- Student 10 enrollments
        IF course9_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course9_id, student10_id, 8, FALSE,
                    NOW() - INTERVAL '15 days',
                    NOW() - INTERVAL '3 days');
        END IF;

        IF course10_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course10_id, student10_id, 12, TRUE,
                    NOW() - INTERVAL '13 days',
                    NOW() - INTERVAL '6 days');
        END IF;

        -- Student 11 enrollments
        IF course1_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course1_id, student11_id, 10, TRUE,
                    NOW() - INTERVAL '13 days',
                    NOW() - INTERVAL '4 days');
        END IF;

        IF course20_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course20_id, student11_id, 6, FALSE,
                    NOW() - INTERVAL '11 days',
                    NOW() - INTERVAL '1 day');
        END IF;

        -- Student 12 enrollments
        IF course20_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course20_id, student12_id, 9, FALSE,
                    NOW() - INTERVAL '11 days',
                    NOW() - INTERVAL '2 days');
        END IF;

        -- Student 13 enrollments
        IF course2_id IS NOT NULL THEN
            INSERT INTO course_enrollments (course_id, student_id, content_completed, is_completed, created_at,
                                            last_accessed_at)
            VALUES (course2_id, student13_id, 10, TRUE,
                    NOW() - INTERVAL '9 days',
                    NOW() - INTERVAL '3 days');
        END IF;

        -- Update enrollment_count for all courses
        UPDATE courses c
        SET enrollment_count = (SELECT COUNT(*)
                                FROM course_enrollments ce
                                WHERE ce.course_id = c.id);

        -- Update ratings for selected courses
        -- Only update if the course exists
        IF course1_id IS NOT NULL THEN
            UPDATE courses
            SET rating       = 4.7,
                rating_count = 125,
                total_rating = 4.7 * 125
            WHERE id = course1_id;
        END IF;

        IF course2_id IS NOT NULL THEN
            UPDATE courses
            SET rating       = 4.8,
                rating_count = 95,
                total_rating = 4.8 * 95
            WHERE id = course2_id;
        END IF;

        IF course5_id IS NOT NULL THEN
            UPDATE courses
            SET rating       = 4.6,
                rating_count = 80,
                total_rating = 4.6 * 80
            WHERE id = course5_id;
        END IF;

        IF course16_id IS NOT NULL THEN
            UPDATE courses
            SET rating       = 4.9,
                rating_count = 110,
                total_rating = 4.9 * 110
            WHERE id = course16_id;
        END IF;

        IF course9_id IS NOT NULL THEN
            UPDATE courses
            SET rating       = 4.5,
                rating_count = 75,
                total_rating = 4.5 * 75
            WHERE id = course9_id;
        END IF;

        IF course12_id IS NOT NULL THEN
            UPDATE courses
            SET rating       = 4.7,
                rating_count = 90,
                total_rating = 4.7 * 90
            WHERE id = course12_id;
        END IF;

        IF course20_id IS NOT NULL THEN
            UPDATE courses
            SET rating       = 4.8,
                rating_count = 85,
                total_rating = 4.8 * 85
            WHERE id = course20_id;
        END IF;

    END
$$;
