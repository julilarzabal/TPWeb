-- name: CreateUser :one
INSERT INTO "user" (name, email, age, industry, role)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, name, email, age, industry, role;

-- name: GetUser :one
SELECT id, name, email, age, industry, role
FROM "user"
WHERE id = $1;

-- name: ListUsers :many
SELECT id, name, email, age, industry, role
FROM "user"
ORDER BY name;

-- name: UpdateUser :exec
UPDATE "user"
SET name = $2,
    email = $3,
    age = $4,
    industry = $5,
    role = $6
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM "user"
WHERE id = $1;

-- name: CreateCategory :one
INSERT INTO category (name)
VALUES ($1)
RETURNING id, name;

-- name: GetCategory :one
SELECT id, name
FROM category
WHERE id = $1;

-- name: ListCategories :many
SELECT id, name
FROM category
ORDER BY name;

-- name: DeleteCategory :exec
DELETE FROM category
WHERE id = $1;

-- name: AddUserPreference :exec
INSERT INTO user_preference (user_id, category_id)
VALUES ($1, $2)
ON CONFLICT (user_id, category_id) DO NOTHING;

-- name: GetPreferencesByUser :many
SELECT c.id, c.name
FROM category c
JOIN user_preference up ON c.id = up.category_id
WHERE up.user_id = $1;

-- name: RemoveUserPreference :exec
DELETE FROM user_preference
WHERE user_id = $1 AND category_id = $2;

-- name: CreateEvent :one
INSERT INTO event (title, description, location, age_range, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, title, description, location, age_range, user_id;

-- name: GetEvent :one
SELECT id, title, description, location, age_range, user_id
FROM event
WHERE id = $1;

-- name: ListEvents :many
SELECT id, title, description, location, age_range, user_id
FROM event
ORDER BY id DESC;

-- name: UpdateEvent :exec
UPDATE event
SET title = $2,
    description = $3,
    location = $4,
    age_range = $5,
    user_id = $6
WHERE id = $1;

-- name: DeleteEvent :exec
DELETE FROM event
WHERE id = $1;

-- name: AddEventCategory :exec
INSERT INTO event_category (event_id, category_id)
VALUES ($1, $2)
ON CONFLICT (event_id, category_id) DO NOTHING;

-- name: GetCategoriesByEvent :many
SELECT c.id, c.name
FROM category c
JOIN event_category ec ON c.id = ec.category_id
WHERE ec.event_id = $1;

-- name: RemoveEventCategory :exec
DELETE FROM event_category
WHERE event_id = $1 AND category_id = $2;