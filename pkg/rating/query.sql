-- name: GetTop30SoloRating :many
SELECT char_name, char_guild_name, char_level, solo_rating, solo_wins, solo_loses, (solo_wins + solo_loses) as match_count
FROM bgrating WHERE solo_wins + solo_loses > 9 ORDER BY solo_rating DESC, solo_wins + solo_loses LIMIT 30;

-- name: GetTop30PartyRating :many
SELECT char_name, char_guild_name, char_level, party_rating, party_wins, party_loses, (party_wins + party_loses) as match_count
FROM bgrating WHERE party_wins + party_loses > 9 ORDER BY party_rating DESC, party_wins + party_loses LIMIT 30;