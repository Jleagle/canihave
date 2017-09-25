CREATE TABLE `categories` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `amazonName` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `amazon` (`amazonName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `items` (
  `id` varchar(10) NOT NULL DEFAULT '',
  `dateCreated` int(10) NOT NULL DEFAULT '0',
  `dateUpdated` int(10) NOT NULL DEFAULT '0',
  `dateScanned` int(10) NOT NULL DEFAULT '0',
  `name` varchar(511) NOT NULL DEFAULT '',
  `link` varchar(255) NOT NULL,
  `source` varchar(255) NOT NULL DEFAULT '',
  `salesRank` int(10) unsigned NOT NULL,
  `photo` varchar(255) NOT NULL DEFAULT '',
  `node` varchar(255) NOT NULL,
  `nodeName` varchar(255) NOT NULL DEFAULT '',
  `price` int(10) unsigned NOT NULL,
  `region` varchar(4) NOT NULL DEFAULT '',
  `hits` int(10) unsigned NOT NULL DEFAULT '0',
  `type` varchar(255) NOT NULL DEFAULT '',
  `companyName` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  KEY `hits` (`hits`),
  KEY `dateCreated` (`dateCreated`),
  KEY `salesRank` (`salesRank`),
  KEY `price` (`price`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `relations` (
  `id` varchar(10) NOT NULL DEFAULT '',
  `relatedId` varchar(10) NOT NULL DEFAULT '',
  `dateCreated` int(10) NOT NULL,
  `type` varchar(10) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`,`relatedId`,`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
