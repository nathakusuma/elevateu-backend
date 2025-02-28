CREATE TABLE courses
(
    id                 UUID PRIMARY KEY,
    category_id        UUID                     NOT NULL REFERENCES categories (id),
    title              VARCHAR(50)              NOT NULL,
    description        VARCHAR(1000)            NOT NULL,
    teacher_name       VARCHAR(60)              NOT NULL,
    teacher_avatar_url TEXT                     NOT NULL,
    thumbnail_url      TEXT                     NOT NULL,
    preview_video_url  TEXT                     NOT NULL,
    rating             DOUBLE PRECISION         NOT NULL DEFAULT 0,
    rating_count       BIGINT                   NOT NULL DEFAULT 0,
    total_rating       DOUBLE PRECISION         NOT NULL DEFAULT 0,
    enrollment_count   BIGINT                   NOT NULL DEFAULT 0,
    content_count      INT                      NOT NULL DEFAULT 0,
    created_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
