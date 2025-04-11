CREATE TABLE files (
    id TEXT PRIMARY KEY,
    userId TEXT NOT NULL,
    filepath TEXT NOT NULL,
    filename TEXT NOT NULL,
    uploadTime TIMESTAMPTZ DEFAULT now() NOT NULL,
    FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
);

ALTER TABLE users ADD COLUMN createdAt TIMESTAMPTZ DEFAULT now() NOT NULL;
ALTER TABLE sessions DROP CONSTRAINT sessions_userid_fkey, ADD CONSTRAINT sessions_userid_fkey FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE;
