-- -----------------------------------------------------
-- Schema playstation
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `playstation` DEFAULT CHARACTER SET utf8 ;
USE `playstation` ;

-- -----------------------------------------------------
-- Table `playstation`.`auth`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `playstation`.`auth` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `session_key` VARCHAR(64) NULL,
  `timestamp` DATETIME NULL,
  `uid` INT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `playstation`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `playstation`.`users` (
  `timestamp_creation` DATETIME NULL,
  `uid` INT NOT NULL AUTO_INCREMENT,
  `login` VARCHAR(16) NULL,
  `password` VARCHAR(64) NULL,
  `mail` VARCHAR(32) NULL,
  `wishlist` MEDIUMTEXT NULL,
  `favorite_genres` LONGTEXT NULL,
  PRIMARY KEY (`uid`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `playstation`.`games`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `playstation`.`games` (
  `game_id` INT NOT NULL AUTO_INCREMENT,
  `genres` TEXT(250) NULL,
  `rating` VARCHAR(1) NULL,
  `developer` VARCHAR(45) NULL,
  `ofplayers` INT NULL,
  `name` VARCHAR(45) NULL,
  `img_link` VARCHAR(200) NULL,
  `summary` LONGTEXT NULL,
  `metascore` INT NULL,
  `users_score` FLOAT NULL,
  PRIMARY KEY (`game_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`genres`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `playstation`.`genres` (
  `genre_id` INT NOT NULL AUTO_INCREMENT,
  `original_genre_name` VARCHAR(45) NULL,
  `rus_genre_name` VARCHAR(45) NULL,
  PRIMARY KEY (`genre_id`))
ENGINE = InnoDB;

