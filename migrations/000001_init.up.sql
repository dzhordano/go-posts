CREATE TYPE verification_type AS (
    code TEXT,
    verified BOOLEAN
);

CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    verification verification_type,
    suspended BOOLEAN DEFAULT FALSE,
    registered TIMESTAMP WITHOUT TIME ZONE,
    lastonline TIMESTAMP WITHOUT TIME ZONE
);
