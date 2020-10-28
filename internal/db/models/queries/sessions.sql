-- name: GetSession :one
SELECT * FROM sessions
where value = @value;

-- name: CreateSession :exec
INSERT INTO sessions (user_name, value)
VALUES (@user_name, @value);
