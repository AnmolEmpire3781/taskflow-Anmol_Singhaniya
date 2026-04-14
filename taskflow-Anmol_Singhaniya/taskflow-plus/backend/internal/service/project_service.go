package service

import (
    "context"

    "github.com/plus/taskflow/backend/internal/models"
    "github.com/plus/taskflow/backend/internal/repository"
)

type ProjectService struct {
    projects *repository.ProjectRepository
    tasks *repository.TaskRepository
}
func NewProjectService(projects *repository.ProjectRepository, tasks *repository.TaskRepository) *ProjectService { return &ProjectService{projects: projects, tasks: tasks} }

func (s *ProjectService) List(ctx context.Context, userID string, page, limit int) ([]models.Project, error) { return s.projects.ListAccessible(ctx, userID, page, limit) }
func (s *ProjectService) Create(ctx context.Context, userID, name, description string) (models.Project, error) { return s.projects.Create(ctx, userID, name, description) }
func (s *ProjectService) Get(ctx context.Context, projectID string, page, limit int) (models.Project, error) {
    p, err := s.projects.GetByID(ctx, projectID); if err != nil { return p, err }
    tasks, err := s.tasks.ListByProject(ctx, projectID, "", "", page, limit); if err != nil { return p, err }
    p.Tasks = tasks
    return p, nil
}
func (s *ProjectService) Update(ctx context.Context, actorID, projectID, name, description string) (models.Project, error) {
    p, err := s.projects.GetByID(ctx, projectID); if err != nil { return p, err }
    if p.OwnerID != actorID { return p, repository.ErrForbidden }
    return s.projects.Update(ctx, projectID, name, description)
}
func (s *ProjectService) Delete(ctx context.Context, actorID, projectID string) error {
    p, err := s.projects.GetByID(ctx, projectID); if err != nil { return err }
    if p.OwnerID != actorID { return repository.ErrForbidden }
    return s.projects.Delete(ctx, projectID)
}
func (s *ProjectService) Stats(ctx context.Context, actorID, projectID string) (models.ProjectStats, error) {
    p, err := s.projects.GetByID(ctx, projectID); if err != nil { return models.ProjectStats{}, err }
    if p.OwnerID != actorID { return models.ProjectStats{}, repository.ErrForbidden }
    return s.projects.Stats(ctx, projectID)
}
