
-- -----------------------------------------------------
-- Schema ps4-recommendation-system
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `ps4-recommendation-system` DEFAULT CHARACTER SET utf8 ;
USE `ps4-recommendation-system` ;
 
-- -----------------------------------------------------
-- Table `ps4-recommendation-system`.`users`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `ps4-recommendation-system`.`users` (
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
-- Table `ps4-recommendation-system`.`auth`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `ps4-recommendation-system`.`auth` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `session_key` VARCHAR(64) NULL,
  `timestamp` DATETIME NULL,
  `uid` INT NULL,
  `users_uid` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_auth_users_idx` (`users_uid` ASC),
  CONSTRAINT `fk_auth_users`
    FOREIGN KEY (`users_uid`)
    REFERENCES `ps4-recommendation-system`.`users` (`uid`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;
 
 
-- -----------------------------------------------------
-- Table `ps4-recommendation-system`.`games`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `ps4-recommendation-system`.`games` (
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