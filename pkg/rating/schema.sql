CREATE TABLE `bgrating` (
  `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `rating` int(11) NOT NULL DEFAULT 1000,
  `char_name` varchar(10) NOT NULL,
  `char_guild_name` varchar(16) NOT NULL DEFAULT '',
  `char_level` int(11) NOT NULL DEFAULT 0,
  `battles_count` int(11) NOT NULL DEFAULT 0,
  `last_match_time` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;