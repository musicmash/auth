-- name: EnsureUserExists :exec
INSERT INTO users (name, photo)
VALUES (@name, @photo)
ON CONFLICT DO NOTHING;
