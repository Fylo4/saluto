-- name: GetMessages :many
SELECT
  id,
  displayName,
  body,
  timePosted
FROM post
ORDER BY timePosted desc;

-- name: PostMessage :exec
INSERT INTO post (displayName, body, timePosted)
VALUES ($1, $2, CURRENT_TIMESTAMP(2));
