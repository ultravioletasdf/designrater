CREATE TABLE IF NOT EXISTS users (
    id integer PRIMARY KEY,
    username text NOT NULL,
    password_hash text NOT NULL,
    
    created_at integer NOT NULL,
    updated_at integer NOT NULL
);
CREATE UNIQUE INDEX idx_users_id ON users (id);
CREATE UNIQUE INDEX idx_users_username ON users (username);