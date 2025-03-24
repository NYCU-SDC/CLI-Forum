CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    author_id UUID,
    title VARCHAR(200),
    content TEXT,
    create_at TIMESTAMPTZ DEFAULT now()
);