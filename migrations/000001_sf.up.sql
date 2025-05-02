CREATE TABLE users (
    id 	   TEXT NOT NULL PRIMARY KEY,
    email  TEXT NOT NULL UNIQUE,
    pwd    TEXT NOT NULL
);

CREATE TABLE sessions (
    id     TEXT NOT NULL PRIMARY KEY,
    user_id TEXT NOT NULL
);

ALTER TABLE sessions 
    ADD FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
