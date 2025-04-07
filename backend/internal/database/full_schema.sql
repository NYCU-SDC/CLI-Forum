-- Code generated by schema merge script. DO NOT EDIT.

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users
(
    id       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name     VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS posts (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     author_id UUID REFERENCES users(id) NOT NULL,
     title VARCHAR(200),
     content TEXT,
     create_at TIMESTAMPTZ DEFAULT now()
);