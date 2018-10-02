CREATE TABLE `users` (
	`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
	`username` VARCHAR(255) NOT NULL DEFAULT '',
    `email_hash` VARBINARY(255) NOT NULL DEFAULT '',
    `email_crypt` VARBINARY(255) NOT NULL DEFAULT '',
    `password_hash` VARBINARY(255) NOT NULL DEFAULT '',
    `password_salt` VARBINARY(255) NOT NULL DEFAULT '',
    `bio` TEXT,
    `image` TEXT,
	PRIMARY KEY (`id`) AUTO_INCREMENT,
    INDEX `username` (`username`),
	UNIQUE INDEX `email_hash` (`email_hash`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `follows` (
    `follower_id` INT UNSIGNED NOT NULL,
    `followed_id` INT UNSIGNED NOT NULL,
	PRIMARY KEY (`follower_id`, `followed_id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `favorites` (
    `user_id` INT UNSIGNED NOT NULL,
    `article_id` INT UNSIGNED NOT NULL,
	PRIMARY KEY (`user_id`, `article_id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `article` (
    `article_id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `author_id` INT UNSIGNED NOT NULL,
    `slug` VARCHAR(255),
    `title` VARCHAR(255),
    `description` TEXT,
    `body` TEXT,
    `created_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`article_id`),
    INDEX `author_id` (`author_id`),
    INDEX `slug` (`slug`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `article_tags` (
    `article_id` INT UNSIGNED NOT NULL,
    `tag` VARCHAR(255),
	PRIMARY KEY (`article_id`, `tag`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;

CREATE TABLE `comments` (
    `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
    `user_id` INT UNSIGNED NOT NULL,
    `article_id` INT UNSIGNED NOT NULL,
    `body` TEXT,
    `created_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_date` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (`id`),
    INDEX `user_id` (`user_id`),
    INDEX `article_id` (`article_id`)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB
;
