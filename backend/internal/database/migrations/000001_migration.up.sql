CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    author_id UUID,
    title VARCHAR(200),
    content TEXT,
    create_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    name VARCHAR(30) NOT NULL,
    password VARCHAR(200) NOT NULL
);