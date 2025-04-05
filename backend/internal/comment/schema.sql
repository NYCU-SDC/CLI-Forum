CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID NOT NULL,
    author_id UUID REFERENCES users(id) NOT NULL,
    title VARCHAR(200),
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);