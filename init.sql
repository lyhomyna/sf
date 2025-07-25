CREATE TABLE users (
    id 	   TEXT NOT NULL PRIMARY KEY,
    email  TEXT NOT NULL UNIQUE,
    pwd    TEXT NOT NULL
);

CREATE TABLE sessions (
    id     TEXT NOT NULL PRIMARY KEY,
    userId TEXT NOT NULL
);

ALTER TABLE sessions 
    ADD FOREIGN KEY (userId) REFERENCES users(id); 
