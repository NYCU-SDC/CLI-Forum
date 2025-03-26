CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    author_id UUID,
    title VARCHAR(200),
    content TEXT,
    create_at TIMESTAMPTZ DEFAULT now()
);