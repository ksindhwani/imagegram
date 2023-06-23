DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS comments;

CREATE TABLE `posts` (
    `post_id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `user_id` INT NOT NULL,
    `caption` TEXT,
    `created_at`  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `comments` (
    `comment_id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `post_id` INT NOT NULL,
    `comment` TEXT,
    `created_at`  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE `images` (
    `image_id` INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `post_id` INT NOT NULL,
    `location` varchar(255),
    `converted_image_location` varchar(255),
    `uploaded_at`  DATETIME DEFAULT CURRENT_TIMESTAMP
);