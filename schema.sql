CREATE TABLE `items` (
  `id` varchar(11) NOT NULL DEFAULT '',
  `dateCreated` date NOT NULL,
  `dateUpdated` date NOT NULL,
  `name` varchar(255) NOT NULL,
  `link` varchar(255) NOT NULL,
  `source` varchar(255) NOT NULL DEFAULT '',
  `salesRank` int(11) unsigned NOT NULL,
  `photo` varchar(255) NOT NULL DEFAULT '',
  `productGroup` varchar(255) NOT NULL DEFAULT '',
  `productTypeName` varchar(255) NOT NULL DEFAULT '',
  `price` decimal(11,2) unsigned NOT NULL,
  `currency` varchar(255) NOT NULL DEFAULT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;