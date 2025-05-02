CREATE TABLE files (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    filename TEXT NOT NULL,
    filepath TEXT NOT NULL,
    size INTEGER NOT NULL,
    hash TEXT NOT NULL,
    last_accessed TIMESTAMPTZ DEFAULT now() NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE users ADD COLUMN created_at TIMESTAMPTZ DEFAULT now() NOT NULL;
ALTER TABLE sessions DROP CONSTRAINT sessions_user_id_fkey, ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
