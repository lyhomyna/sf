DROP TABLE files;
ALTER TABLE users DROP COLUMN created_at;
ALTER TABLE sessions DROP CONSTRAINT sessions_user_id_fkey, ADD CONSTRAINT sessions_user_id_fkey FOREIGN KEY (userid) REFERENCES users(id);
