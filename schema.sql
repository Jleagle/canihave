CREATE TABLE `items` (
  `id` varchar(10) NOT NULL DEFAULT '',
  `dateCreated` datetime NOT NULL,
  `dateUpdated` datetime NOT NULL,
  `name` varchar(255) NOT NULL,
  `link` varchar(255) NOT NULL,
  `source` varchar(255) NOT NULL DEFAULT '',
  `salesRank` int(10) unsigned NOT NULL,
  `photo` varchar(255) NOT NULL DEFAULT '',
  `productGroup` varchar(255) NOT NULL DEFAULT '',
  `price` int(10) unsigned NOT NULL,
  `region` varchar(4) NOT NULL DEFAULT '',
  `hits` int(10) unsigned NOT NULL DEFAULT '0',
  `status` varchar(255) NOT NULL DEFAULT '',
  `type` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `categories` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `amazon_name` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `amazon` (`amazon_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `relation` (
  `id` varchar(10) NOT NULL DEFAULT '',
  `related_id` varchar(10) NOT NULL DEFAULT '',
  `date_created` datetime NOT NULL,
  `type` varchar(10) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;