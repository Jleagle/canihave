CREATE TABLE `items` (
  `id` varchar(11) NOT NULL DEFAULT '',
  `dateCreated` date NOT NULL,
  `dateUpdated` date NOT NULL,
  `name` varchar(255) NOT NULL,
  `desc` text NOT NULL,
  `link` varchar(255) NOT NULL,
  `source` int(11) unsigned NOT NULL,
  `salesRank` int(11) unsigned NOT NULL,
  `images` text NOT NULL,
  `productGroup` varchar(255) NOT NULL DEFAULT '',
  `productTypeName` varchar(255) NOT NULL DEFAULT '',
  `price` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `sources` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `domain` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;