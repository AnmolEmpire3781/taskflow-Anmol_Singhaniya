package repository

import (
    "context"
    "errors"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/plus/taskflow/backend/internal/models"
)

type ProjectRepository struct { pool *pgxpool.Pool }
func NewProjectRepository(pool *pgxpool.Pool) *ProjectRepository { return &ProjectRepository{pool: pool} }

func (r *ProjectRepository) ListAccessible(ctx context.Context, userID string, page, limit int) ([]models.Project, error) {
    q := `SELECT DISTINCT p.id,p.name,p.description,p.owner_id,p.created_at
          FROM projects p
          LEFT JOIN tasks t ON t.project_id = p.id
          WHERE p.owner_id=$1 OR t.assignee_id=$1 OR t.creator_id=$1
          ORDER BY p.created_at DESC LIMIT $2 OFFSET $3`
    rows, err := r.pool.Query(ctx, q, userID, limit, (page-1)*limit)
    if err != nil { return nil, err }
    defer rows.Close()
    out := []models.Project{}
    for rows.Next() {
        var p models.Project
        if err := rows.Scan(&p.ID,&p.Name,&p.Description,&p.OwnerID,&p.CreatedAt); err != nil { return nil, err }
        out = append(out, p)
    }
    return out, rows.Err()
}

func (r *ProjectRepository) Create(ctx context.Context, ownerID, name, description string) (models.Project, error) {
    p := models.Project{}
    q := `INSERT INTO projects (id,name,description,owner_id) VALUES ($1,$2,$3,$4) RETURNING id,name,description,owner_id,created_at`
    err := r.pool.QueryRow(ctx, q, uuid.NewString(), name, description, ownerID).Scan(&p.ID,&p.Name,&p.Description,&p.OwnerID,&p.CreatedAt)
    return p, err
}

func (r *ProjectRepository) GetByID(ctx context.Context, projectID string) (models.Project, error) {
    p := models.Project{}
    err := r.pool.QueryRow(ctx, `SELECT id,name,description,owner_id,created_at FROM projects WHERE id=$1`, projectID).Scan(&p.ID,&p.Name,&p.Description,&p.OwnerID,&p.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) { return p, ErrNotFound }
    return p, err
}

func (r *ProjectRepository) Update(ctx context.Context, projectID, name, description string) (models.Project, error) {
    p := models.Project{}
    q := `UPDATE projects SET name=$2, description=$3 WHERE id=$1 RETURNING id,name,description,owner_id,created_at`
    err := r.pool.QueryRow(ctx, q, projectID, name, description).Scan(&p.ID,&p.Name,&p.Description,&p.OwnerID,&p.CreatedAt)
    if errors.Is(err, pgx.ErrNoRows) { return p, ErrNotFound }
    return p, err
}

func (r *ProjectRepository) Delete(ctx context.Context, projectID string) error {
    cmd, err := r.pool.Exec(ctx, `DELETE FROM projects WHERE id=$1`, projectID)
    if err != nil { return err }
    if cmd.RowsAffected() == 0 { return ErrNotFound }
    return nil
}

func (r *ProjectRepository) Stats(ctx context.Context, projectID string) (models.ProjectStats, error) {
    stats := models.ProjectStats{
    ProjectID:  projectID,
    ByStatus:   map[string]int{},
    ByAssignee: map[string]int{},
    }
    rows, err := r.pool.Query(ctx, `SELECT status, COUNT(*) FROM tasks WHERE project_id=$1 GROUP BY status`, projectID)
    if err != nil { return stats, err }
    defer rows.Close()
    for rows.Next() {
        var status string; var count int
        if err := rows.Scan(&status, &count); err != nil { return stats, err }
        stats.ByStatus[status] = count
    }
    rows2, err := r.pool.Query(ctx, `SELECT COALESCE(assignee_id::text, 'unassigned'), COUNT(*) FROM tasks WHERE project_id=$1 GROUP BY assignee_id`, projectID)
    if err != nil { return stats, err }
    defer rows2.Close()
    for rows2.Next() {
        var assignee string; var count int
        if err := rows2.Scan(&assignee, &count); err != nil { return stats, err }
        stats.ByAssignee[assignee] = count
    }
    return stats, nil
}
