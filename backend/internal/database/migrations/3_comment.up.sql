CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID REFERENCES posts(id),
    author_id UUID REFERENCES users(id),
    title VARCHAR(200),
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);