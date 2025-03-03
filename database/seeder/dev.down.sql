BEGIN;

-- Disable foreign key constraints temporarily to avoid dependency issues
SET CONSTRAINTS ALL DEFERRED;

-- Delete all seeded courses
DELETE
FROM courses
WHERE id IN (
             '01890e5d-05a6-7b0b-a7b3-3f8ed37b5261',
             '01890e5d-0ec9-7e29-ae5f-4a7c66f5d456',
             '01890e5d-1657-71ae-bcba-54a2a55b1bb8',
             '01890e5d-1e20-7e5c-b0b2-6dbc90efac62',
             '01890e5d-25de-7d73-bf1c-8c1a13f17839',
             '01890e5d-2dd5-7bde-994a-76e4b9e18a08',
             '01890e5d-3530-7487-8b77-fc5cc0773b3f',
             '01890e5d-3c4a-72fe-b78b-e1ebf0cbed87',
             '01890e5d-437c-72ca-8ddc-4aa6b26ddc11',
             '01890e5d-4a4c-7452-b55a-7b8fa75d84e5',
             '01890e5d-516c-7dc0-9d9c-1d6ec6e36fc4',
             '01890e5d-58ca-787d-9ed5-31fbe2bc9429',
             '01890e5d-6086-7c90-a81b-9dccca9a6c97',
             '01890e5d-6769-7b84-baba-4b3aad93a8b0',
             '01890e5d-6eda-71bf-8a5e-a74a7f8c8d2c',
             '01890e5d-76d1-75af-9c61-0adf9a0d8db1',
             '01890e5d-7e5d-7bfb-9e26-2eb534b456f3',
             '01890e5d-8540-7a65-9b7d-5d95efed8bc9',
             '01890e5d-8c74-7b15-9ed8-15b5e14ab0c7',
             '01890e5d-946d-7a56-a9b4-69d7c6fefb92',
             '01890e5d-9b95-72e2-a064-a37c64ac6d3a',
             '01890e5d-a316-7e25-bad9-0ec1a4eebf79',
             '01890e5d-aa41-7ec4-9278-a1eaf2fa3b40',
             '01890e5d-b181-7f2d-aeda-c4e69b0b3d75',
             '01890e5d-b8cf-7d50-b3a1-1acf37eff9af'
    );

-- Delete all seeded categories
DELETE
FROM categories
WHERE id IN (
             '01890e5c-b334-7c38-8c8c-b05c57968489',
             '01890e5c-bed5-7ebb-8dcf-ada8c40c2253',
             '01890e5c-c5b6-7da1-b35c-70d83a37d078',
             '01890e5c-cdbe-718c-9623-ec4a1bdfaa27',
             '01890e5c-d583-752e-ab1e-24d75a5c4dfc',
             '01890e5c-dc95-7de2-9e36-00a6a28bf8de',
             '01890e5c-e31e-7b65-b6fb-08a3b7a00d56'
    );

-- Delete seeded mentor-specific data
DELETE
FROM mentors
WHERE user_id IN (
                  '0194bc32-e9fc-405c-801a-08ff3b0cf28b',
                  '0194bc69-d87c-40d5-809a-17513f3d2b98'
    );

-- Delete seeded student-specific data
DELETE
FROM students
WHERE user_id IN (
                  '0194bb8e-1e7c-4082-806a-1c9483a59a1b',
                  '0194bbc5-0cfc-408e-802d-a854785c39d9',
                  '0194bbfb-fb7c-40d9-8029-7ca996ab5da6'
    );

-- Delete all seeded users
DELETE
FROM users
WHERE id IN (
             '0194bb8e-1e7c-4082-806a-1c9483a59a1b', -- Student 1
             '0194bbc5-0cfc-408e-802d-a854785c39d9', -- Student 2
             '0194bbfb-fb7c-40d9-8029-7ca996ab5da6', -- Student 3
             '0194bc32-e9fc-405c-801a-08ff3b0cf28b', -- Mentor 1
             '0194bc69-d87c-40d5-809a-17513f3d2b98', -- Mentor 2
             '0194bca0-c6fc-40ec-8051-68235bef9817' -- Admin
    );

-- Re-enable foreign key constraints
SET CONSTRAINTS ALL IMMEDIATE;

COMMIT;
