package repository

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/plus/taskflow/backend/internal/models"
)

type TaskRepository struct { pool *pgxpool.Pool }
func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository { return &TaskRepository{pool: pool} }

func (r *TaskRepository) ListByProject(ctx context.Context, projectID, status, assignee string, page, limit int) ([]models.Task, error) {
    base := `SELECT id,title,description,status,priority,project_id,assignee_id,creator_id,due_date,created_at,updated_at FROM tasks WHERE project_id=$1`
    args := []any{projectID}
    pos := 2
    if status != "" { base += fmt.Sprintf(" AND status=$%d", pos); args = append(args, status); pos++ }
    if assignee != "" { base += fmt.Sprintf(" AND assignee_id=$%d", pos); args = append(args, assignee); pos++ }
    base += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", pos, pos+1)
    args = append(args, limit, (page-1)*limit)
    rows, err := r.pool.Query(ctx, base, args...)
    if err != nil { return nil, err }
    defer rows.Close()
    tasks := []models.Task{}
    for rows.Next() {
        var t models.Task
        if err := rows.Scan(&t.ID,&t.Title,&t.Description,&t.Status,&t.Priority,&t.ProjectID,&t.AssigneeID,&t.CreatorID,&t.DueDate,&t.CreatedAt,&t.UpdatedAt); err != nil { return nil, err }
        tasks = append(tasks, t)
    }
    return tasks, rows.Err()
}

func (r *TaskRepository) Create(ctx context.Context, title, description, priority, projectID, creatorID string, assigneeID *string, dueDate *time.Time) (models.Task, error) {
    t := models.Task{}
    q := `INSERT INTO tasks (id,title,description,status,priority,project_id,assignee_id,creator_id,due_date) VALUES ($1,$2,$3,'todo',$4,$5,$6,$7,$8)
          RETURNING id,title,description,status,priority,project_id,assignee_id,creator_id,due_date,created_at,updated_at`
    err := r.pool.QueryRow(ctx, q, uuid.NewString(), title, description, priority, projectID, assigneeID, creatorID, dueDate).Scan(&t.ID,&t.Title,&t.Description,&t.Status,&t.Priority,&t.ProjectID,&t.AssigneeID,&t.CreatorID,&t.DueDate,&t.CreatedAt,&t.UpdatedAt)
    return t, err
}

func (r *TaskRepository) GetByID(ctx context.Context, taskID string) (models.Task, error) {
    t := models.Task{}
    err := r.pool.QueryRow(ctx, `SELECT id,title,description,status,priority,project_id,assignee_id,creator_id,due_date,created_at,updated_at FROM tasks WHERE id=$1`, taskID).Scan(&t.ID,&t.Title,&t.Description,&t.Status,&t.Priority,&t.ProjectID,&t.AssigneeID,&t.CreatorID,&t.DueDate,&t.CreatedAt,&t.UpdatedAt)
    if errors.Is(err, pgx.ErrNoRows) { return t, ErrNotFound }
    return t, err
}

func (r *TaskRepository) Update(ctx context.Context, taskID string, fields map[string]any) (models.Task, error) {
    if len(fields) == 0 { return r.GetByID(ctx, taskID) }
    query := `UPDATE tasks SET `
    args := []any{taskID}
    i := 2
    first := true
    for col, val := range fields {
        if !first { query += ", " }
        first = false
        query += fmt.Sprintf("%s=$%d", col, i)
        args = append(args, val)
        i++
    }
    query += fmt.Sprintf(", updated_at=NOW() WHERE id=$1 RETURNING id,title,description,status,priority,project_id,assignee_id,creator_id,due_date,created_at,updated_at")
    t := models.Task{}
    err := r.pool.QueryRow(ctx, query, args...).Scan(&t.ID,&t.Title,&t.Description,&t.Status,&t.Priority,&t.ProjectID,&t.AssigneeID,&t.CreatorID,&t.DueDate,&t.CreatedAt,&t.UpdatedAt)
    if errors.Is(err, pgx.ErrNoRows) { return t, ErrNotFound }
    return t, err
}

func (r *TaskRepository) Delete(ctx context.Context, taskID string) error {
    cmd, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id=$1`, taskID)
    if err != nil { return err }
    if cmd.RowsAffected() == 0 { return ErrNotFound }
    return nil
}
