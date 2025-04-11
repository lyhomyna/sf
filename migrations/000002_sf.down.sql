DROP TABLE files;
ALTER TABLE users DROP COLUMN createdAt;
ALTER TABLE sessions DROP CONSTRAINT sessions_userid_fkey, ADD CONSTRAINT sessions_userid_fkey FOREIGN KEY (userid) REFERENCES users(id);
