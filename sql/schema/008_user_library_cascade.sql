-- +goose up
ALTER TABLE library
ADD COLUMN user_id UUID NOT NULL,
ADD FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE users
DROP COLUMN library_id,
DROP CONSTRAINT users_library_id_fkey;

-- +goose down
ALTER TABLE library
DROP COLUMN user_id;

ALTER TABLE users
ADD COLUMN library_id UUID NOT NULL,
ADD FOREIGN KEY (library_id) REFERENCES library(id) ON DELETE CASCADE;
