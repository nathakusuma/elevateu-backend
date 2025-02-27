CREATE TABLE categories
(
    id   UUID PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

CREATE UNIQUE INDEX categories_name_key ON categories (name);
