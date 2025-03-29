CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_id UUID,
    author_id UUID,
    TITLE VACHAR(200),
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);