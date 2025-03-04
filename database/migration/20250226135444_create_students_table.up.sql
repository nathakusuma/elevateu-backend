CREATE TABLE students
(
    user_id                    UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    instance                   VARCHAR(50) NOT NULL,
    major                      VARCHAR(50) NOT NULL,
    point                      INT         NOT NULL DEFAULT 0,
    subscribed_boost_until     TIMESTAMP WITH TIME ZONE,
    subscribed_challenge_until TIMESTAMP WITH TIME ZONE
);
