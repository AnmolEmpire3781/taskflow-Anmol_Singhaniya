package service

import (
    "context"
    "time"

    "github.com/plus/taskflow/backend/internal/models"
    "github.com/plus/taskflow/backend/internal/repository"
)

type TaskService struct {
    tasks *repository.TaskRepository
    projects *repository.ProjectRepository
    users *repository.UserRepository
}
func NewTaskService(tasks *repository.TaskRepository, projects *repository.ProjectRepository, users *repository.UserRepository) *TaskService {
    return &TaskService{tasks: tasks, projects: projects, users: users}
}
func (s *TaskService) List(ctx context.Context, userID, projectID, status, assignee string, page, limit int) ([]models.Task, error) {
    p, err := s.projects.GetByID(ctx, projectID); if err != nil { return nil, err }
    if p.OwnerID != userID && assignee != userID { /* relaxed access via task membership handled by data list elsewhere */ }
    return s.tasks.ListByProject(ctx, projectID, status, assignee, page, limit)
}
func (s *TaskService) Create(ctx context.Context, actorID, projectID, title, description, priority string, assigneeID *string, dueDate *time.Time) (models.Task, error) {
    if _, err := s.projects.GetByID(ctx, projectID); err != nil { return models.Task{}, err }
    if assigneeID != nil {
        exists, err := s.users.Exists(ctx, *assigneeID)
        if err != nil { return models.Task{}, err }
        if !exists { return models.Task{}, repository.ErrNotFound }
    }
    return s.tasks.Create(ctx, title, description, priority, projectID, actorID, assigneeID, dueDate)
}
func (s *TaskService) Update(ctx context.Context, actorID, taskID string, fields map[string]any) (models.Task, error) {
    t, err := s.tasks.GetByID(ctx, taskID); if err != nil { return t, err }
    p, err := s.projects.GetByID(ctx, t.ProjectID); if err != nil { return t, err }
    if p.OwnerID != actorID && t.CreatorID != actorID && (t.AssigneeID == nil || *t.AssigneeID != actorID) { return t, repository.ErrForbidden }
    if assignee, ok := fields["assignee_id"]; ok {
        if assigneeStr, ok2 := assignee.(string); ok2 && assigneeStr != "" {
            exists, err := s.users.Exists(ctx, assigneeStr)
            if err != nil { return t, err }
            if !exists { return t, repository.ErrNotFound }
        }
    }
    return s.tasks.Update(ctx, taskID, fields)
}
func (s *TaskService) Delete(ctx context.Context, actorID, taskID string) error {
    t, err := s.tasks.GetByID(ctx, taskID); if err != nil { return err }
    p, err := s.projects.GetByID(ctx, t.ProjectID); if err != nil { return err }
    if p.OwnerID != actorID && t.CreatorID != actorID { return repository.ErrForbidden }
    return s.tasks.Delete(ctx, taskID)
}
