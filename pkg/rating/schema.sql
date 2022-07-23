CREATE TABLE `bgrating` (
  `id` int(11) NOT NULL AUTO_INCREMENT PRIMARY KEY,
  `char_name` varchar(10) NOT NULL,
  `char_guild_name` varchar(16) NOT NULL DEFAULT '',
  `char_level` int(11) NOT NULL DEFAULT 0,
  `solo_rating` int(11) NOT NULL DEFAULT 1000,
  `solo_wins` int(11) NOT NULL DEFAULT 0,
  `solo_loses` int(11) NOT NULL DEFAULT 0,
  `last_solo_match_time` datetime DEFAULT NULL,
  `party_rating` int(11) NOT NULL DEFAULT 1000,
  `party_wins` int(11) NOT NULL DEFAULT 0,
  `party_loses` int(11) NOT NULL DEFAULT 0,
  `last_party_match_time` datetime DEFAULT NULL,
  `battles_count` int(11) NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;