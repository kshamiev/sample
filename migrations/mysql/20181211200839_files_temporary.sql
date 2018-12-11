-- +goose Up
-- SQL in this section is executed when the migration is applied.
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

CREATE TABLE IF NOT EXISTS `filesTemporary` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Уникальный идентификатор записи',
  `createAt` DATETIME NULL DEFAULT NULL COMMENT 'Дата и время создания записи',
  `updateAt` DATETIME NULL DEFAULT NULL COMMENT 'Дата и время обновления записи',
  `deleteAt` DATETIME NULL DEFAULT NULL COMMENT 'Дата и время удаления записи (пометка на удаление)',
  `accessAt` DATETIME NULL DEFAULT NULL COMMENT 'Дата и время последнего доступа к записи',
  `filename` VARCHAR(4096) NULL DEFAULT NULL COMMENT 'Оригинальное имя файла',
  `fileExt` VARCHAR(256) NULL DEFAULT NULL COMMENT 'Расширение имени файла без точки',
  `size` BIGINT(20) UNSIGNED NOT NULL DEFAULT '0' COMMENT 'Размер файла в байтах',
  `sha512` VARCHAR(128) NULL DEFAULT NULL COMMENT 'SHA512 контрольная сумма файла в HEX формате',
  `localPath` VARCHAR(4096) NULL DEFAULT NULL COMMENT 'Относительный путь и имя файла',
  `contentType` TEXT NULL DEFAULT NULL COMMENT 'MIME Content-Type загруженного файла',
  PRIMARY KEY (`id`),
  INDEX `createAt` (`createAt` ASC),
  INDEX `deleteAt` (`deleteAt` ASC),
  INDEX `fileExt` (`fileExt` ASC))
ENGINE = InnoDB
AUTO_INCREMENT=1
DEFAULT CHARACTER SET = utf8
COMMENT = 'Закачанные на сервер временные файлы'
ROW_FORMAT = Dynamic;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

DROP TABLE IF EXISTS `filesTemporary`;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
