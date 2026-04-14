package repository

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/plus/taskflow/backend/internal/models"
)

type UserRepository struct { pool *pgxpool.Pool }
func NewUserRepository(pool *pgxpool.Pool) *UserRepository { return &UserRepository{pool: pool} }

func (r *UserRepository) Create(ctx context.Context, name, email, password string) (models.User, error) {
    u := models.User{}
    q := `INSERT INTO users (id, name, email, password) VALUES ($1,$2,$3,$4) RETURNING id,name,email,created_at`
    id := uuid.NewString()
    err := r.pool.QueryRow(ctx, q, id, name, email, password).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
    if err != nil {
        if isUniqueViolation(err) { return u, ErrConflict }
        return u, err
    }
    return u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (models.User, error) {
    u := models.User{}
    q := `SELECT id,name,email,password,created_at FROM users WHERE email=$1`
    err := r.pool.QueryRow(ctx, q, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) { return u, ErrNotFound }
    return u, err
}

func (r *UserRepository) Exists(ctx context.Context, id string) (bool, error) {
    var exists bool
    err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`, id).Scan(&exists)
    return exists, err
}
