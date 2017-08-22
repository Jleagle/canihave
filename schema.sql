CREATE TABLE `products` (
  `id` varchar(10) NOT NULL DEFAULT '',
  `date_created` date NOT NULL,
  `date_updated` date NOT NULL,
  `times_added` int(11) NOT NULL DEFAULT '1',
  `name` varchar(255) NOT NULL,
  `desc` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
