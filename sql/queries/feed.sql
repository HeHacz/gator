-- name: CreateFeed :one
INSERT INTO feeds (name, id, created_at, updated_at, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedbyName :one
SELECT * FROM feeds WHERE name = $1;

-- name: GetFeedbyURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: DropFeeds :exec
DELETE FROM feeds;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsWithUsername :many
SELECT feeds.ID, feeds.created_at, feeds.updated_at,feeds.name, feeds.url, users.name as username
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id
ORDER BY username; 