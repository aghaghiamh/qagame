-- +migrate Up
ALTER TABLE users ADD COLUMN hashed_password VARCHAR(256) AFTER phone_number;

-- +migrate Down
ALTER TABLE users DROP COLUMN hashed_password;