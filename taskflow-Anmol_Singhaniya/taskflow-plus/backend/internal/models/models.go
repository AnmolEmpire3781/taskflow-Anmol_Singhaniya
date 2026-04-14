package models

import "time"

type User struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    Password  string    `json:"-"`
    CreatedAt time.Time `json:"created_at"`
}

type Project struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    OwnerID     string    `json:"owner_id"`
    CreatedAt   time.Time `json:"created_at"`
    Tasks       []Task    `json:"tasks,omitempty"`
}

type Task struct {
    ID          string     `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description,omitempty"`
    Status      string     `json:"status"`
    Priority    string     `json:"priority"`
    ProjectID   string     `json:"project_id"`
    AssigneeID  *string    `json:"assignee_id,omitempty"`
    CreatorID   string     `json:"creator_id"`
    DueDate     *time.Time `json:"due_date,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type ProjectStats struct {
    ProjectID  string         `json:"project_id"`
    ByStatus   map[string]int `json:"by_status"`
    ByAssignee map[string]int `json:"by_assignee"`
}