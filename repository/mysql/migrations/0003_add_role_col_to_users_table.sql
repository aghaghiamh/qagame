-- +migrate Up
ALTER TABLE `users` ADD COLUMN `role` ENUM('user', 'admin') DEFAULT 'user' NOT NULL;

-- +migrate Down
ALTER TABLE `users` DROP COLUMN `role`;