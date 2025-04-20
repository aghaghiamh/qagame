-- +migrate Up
CREATE TABLE IF NOT EXISTS `access_controls` (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `priviledged_type` ENUM('role', 'user') NOT NULL,
    `priviledged_id` INTEGER NOT NULL,
    `source_permission_id` INTEGER NOT NULL,
    `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (`source_permission_id`) REFERENCES `source_permissions`(`id`)
);

-- +migrate Down
DROP TABLE IF EXISTS `access_controls`;