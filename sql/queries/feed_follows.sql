-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING id, created_at, updated_at, user_id, feed_id
)
SELECT
    iff.id,
    iff.created_at,
    iff.updated_at,
    iff.user_id,
    iff.feed_id,
    users.name AS user_name,
    feeds.name AS feed_name
FROM inserted_feed_follow iff
INNER JOIN users ON iff.user_id = users.id
INNER JOIN feeds ON iff.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.id,
    feed_follows.created_at,
    feed_follows.updated_at,
    feed_follows.user_id,
    feed_follows.feed_id,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1;