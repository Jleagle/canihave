CREATE TABLE `items` (
  `id`           VARCHAR(11)      NOT NULL DEFAULT '',
  `dateCreated`  DATETIME         NOT NULL,
  `dateUpdated`  DATETIME         NOT NULL,
  `name`         VARCHAR(255)     NOT NULL,
  `link`         VARCHAR(255)     NOT NULL,
  `source`       VARCHAR(255)     NOT NULL DEFAULT '',
  `salesRank`    INT(11) UNSIGNED NOT NULL,
  `photo`        VARCHAR(255)     NOT NULL DEFAULT '',
  `productGroup` VARCHAR(255)     NOT NULL DEFAULT '',
  `price`        INT(11) UNSIGNED NOT NULL,
  `region`       VARCHAR(4)       NOT NULL DEFAULT ''
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8;

CREATE TABLE `categories` (
  `id`     INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `amazon` VARCHAR(255)     NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `amazon` (`amazon`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;
