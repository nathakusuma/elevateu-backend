CREATE TABLE students
(
    user_id  UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    instance VARCHAR(50) NOT NULL,
    major    VARCHAR(50) NOT NULL
);
