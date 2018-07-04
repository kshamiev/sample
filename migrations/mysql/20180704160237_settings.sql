-- +goose Up
-- SQL in this section is executed when the migration is applied.
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

CREATE TABLE IF NOT EXISTS `settings` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'Уникальный идентификатор записи',
  `createAt` DATETIME NULL DEFAULT NULL COMMENT 'Дата и время создания записи',
  `updateAt` DATETIME NULL DEFAULT NULL COMMENT 'Дата и время обновления записи',
  `accessAt` DATETIME NULL DEFAULT NULL COMMENT 'Дата и время последнего доступа к записи',
  `key` VARCHAR(255) NOT NULL COMMENT 'Ключ',
  `valueString` LONGTEXT NULL DEFAULT NULL COMMENT 'Строковое значение',
  `valueDate` DATETIME NULL DEFAULT NULL COMMENT 'Значение даты и времени',
  `valueUint` BIGINT(20) UNSIGNED NULL DEFAULT NULL COMMENT 'Числовое unsigned значение',
  `valueInt` BIGINT(20) NULL DEFAULT NULL COMMENT 'Числовое значение',
  `valueDecimal` DECIMAL(16,4) NULL DEFAULT NULL COMMENT 'Значение с плавающей точкой',
  `valueFloat` DOUBLE NULL DEFAULT NULL COMMENT 'IEEE-754 64-bit floating-point number',
  `valueBit` TINYINT(1) NULL DEFAULT NULL COMMENT 'Boolean value',
  `valueBlob` LONGBLOB NULL DEFAULT NULL COMMENT 'Blob value',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `key_UNIQUE` (`key` ASC),
  INDEX `createAt` (`createAt` ASC),
  INDEX `updateAt` (`updateAt` ASC),
  INDEX `accessAt` (`accessAt` ASC))
ENGINE = InnoDB
AUTO_INCREMENT=1
DEFAULT CHARACTER SET = utf8
COMMENT = 'Хранение настроек с типизированными значениями'
ROW_FORMAT = Dynamic;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

DROP TABLE IF EXISTS `settings`;

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
