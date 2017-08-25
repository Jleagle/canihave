CREATE TABLE `items` (
  `id` varchar(11) NOT NULL DEFAULT '',
  `date_created` date NOT NULL,
  `date_updated` date NOT NULL,
  `name` varchar(255) NOT NULL,
  `desc` text NOT NULL,
  `source` int(11) NOT NULL,
  `node` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `sources` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `domain` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;