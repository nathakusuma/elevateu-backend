CREATE TABLE courses
(
    id               UUID PRIMARY KEY,
    category_id      UUID                     NOT NULL REFERENCES categories (id),
    title            VARCHAR(50)              NOT NULL,
    description      VARCHAR(1000)            NOT NULL,
    teacher_name     VARCHAR(60)              NOT NULL,
    rating           DOUBLE PRECISION         NOT NULL DEFAULT 0,
    rating_count     BIGINT                   NOT NULL DEFAULT 0,
    total_rating     DOUBLE PRECISION         NOT NULL DEFAULT 0,
    enrollment_count BIGINT                   NOT NULL DEFAULT 0,
    content_count    INT                      NOT NULL DEFAULT 0,
    total_duration INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX courses_title_idx ON courses USING gist (title gist_trgm_ops);
CREATE INDEX courses_category_id_index ON courses (category_id);
CREATE INDEX courses_total_rating_id_index ON courses (total_rating DESC, id DESC);
