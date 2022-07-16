-- name: GetTop30ByRating :many
SELECT * FROM bgrating WHERE battles_count > 9 ORDER BY rating DESC, battles_count LIMIT 30