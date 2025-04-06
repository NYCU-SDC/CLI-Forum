-- Code generated by schema merge script. DO NOT EDIT.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID REFERENCES posts(id) NOT NULL,
    author_id UUID REFERENCES users(id) NOT NULL,
    title VARCHAR(200),
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name     VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);