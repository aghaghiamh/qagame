-- +migrate Up
CREATE TABLE IF NOT EXISTS `source_permissions` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `title` VARCHAR(191) NOT NULL,
    `description` TEXT,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE IF EXISTS `source_permissions`;