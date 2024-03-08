CREATE TYPE verification_type AS (
    code TEXT,
    verified BOOLEAN
);

CREATE TYPE session_type AS (
    rtoken TEXT,
    expiresat TIMESTAMP
);

CREATE TABLE IF NOT EXISTS admins (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    session session_type,
    registered TIMESTAMP WITHOUT TIME ZONE,
    lastonline TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    verification verification_type,
    session session_type,
    suspended BOOLEAN DEFAULT FALSE,
    registered TIMESTAMP WITHOUT TIME ZONE,
    lastonline TIMESTAMP WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    suspended BOOLEAN DEFAULT FALSE,
    created TIMESTAMP WITHOUT TIME ZONE,
    updated TIMESTAMP WITHOUT TIME ZONE,
    likes INTEGER DEFAULT 0,
    watched INTEGER DEFAULT 0
);

CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    comment TEXT NOT NULL,
    created TIMESTAMP WITHOUT TIME ZONE,
    updated TIMESTAMP WITHOUT TIME ZONE,
    censored BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS users_posts (
    id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts (id) ON DELETE CASCADE NOT NULL,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS users_comments (
    id SERIAL PRIMARY KEY,
    comment_id INTEGER REFERENCES comments (id) ON DELETE CASCADE NOT NULL,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE IF NOT EXISTS post_comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts (id) ON DELETE CASCADE NOT NULL,
    comment_id INTEGER REFERENCES comments (id) ON DELETE CASCADE NOT NULL
);

-- add admin
INSERT INTO admins (name, email, password, registered, lastonline) VALUES ('dzhalal', 'dzhalal.2020@mail.ru', '686a6476626975616e34626e343538392d7976343039753731337639676e75696f676e647338399b53521c219540495d17f03d80294a3ddf3a6ddfe5cd84ea9faf3c6d01535e23', now(), now());