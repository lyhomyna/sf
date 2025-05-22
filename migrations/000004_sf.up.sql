CREATE TABLE directories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id TEXT REFERENCES directories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

ALTER TABLE files 
    ADD COLUMN directory_id TEXT NOT NULL
    REFERENCES directories(id)
    ON DELETE CASCADE;
