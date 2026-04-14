package app

import (
    "encoding/json"
    "log/slog"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/go-chi/chi/v5"
    chimiddleware "github.com/go-chi/chi/v5/middleware"
    "github.com/plus/taskflow/backend/internal/config"
    "github.com/plus/taskflow/backend/internal/httpx"
    authmw "github.com/plus/taskflow/backend/internal/middleware"
    "github.com/plus/taskflow/backend/internal/repository"
    "github.com/plus/taskflow/backend/internal/service"
)

func Router(cfg config.Config, log *slog.Logger, authSvc *service.AuthService, projSvc *service.ProjectService, taskSvc *service.TaskService) http.Handler {
    r := chi.NewRouter()
    r.Use(chimiddleware.RequestID, chimiddleware.RealIP, chimiddleware.Recoverer, chimiddleware.Timeout(30*time.Second))

    r.Get("/health", func(w http.ResponseWriter, r *http.Request) { httpx.JSON(w, http.StatusOK, map[string]string{"status": "ok"}) })

    r.Route("/auth", func(r chi.Router) {
        r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
            var req struct{ Name, Email, Password string }
            if err := json.NewDecoder(r.Body).Decode(&req); err != nil { httpx.Error(w, 400, "invalid json"); return }
            fields := map[string]string{}
            if strings.TrimSpace(req.Name)=="" { fields["name"]="is required" }
            if strings.TrimSpace(req.Email)=="" { fields["email"]="is required" }
            if len(req.Password)<8 { fields["password"]="must be at least 8 characters" }
            if len(fields)>0 { httpx.ValidationError(w, fields); return }
            resp, err := authSvc.Register(r.Context(), req.Name, req.Email, req.Password)
            handleErr(w, err); if err != nil { return }
            httpx.JSON(w, http.StatusCreated, resp)
        })
        r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
            var req struct{ Email, Password string }
            if err := json.NewDecoder(r.Body).Decode(&req); err != nil { httpx.Error(w, 400, "invalid json"); return }
            fields := map[string]string{}
            if strings.TrimSpace(req.Email)=="" { fields["email"]="is required" }
            if strings.TrimSpace(req.Password)=="" { fields["password"]="is required" }
            if len(fields)>0 { httpx.ValidationError(w, fields); return }
            resp, err := authSvc.Login(r.Context(), req.Email, req.Password)
            handleErr(w, err); if err != nil { return }
            httpx.JSON(w, http.StatusOK, resp)
        })
    })

    r.Group(func(r chi.Router) {
        r.Use(authmw.Auth(cfg))
        r.Route("/projects", func(r chi.Router) {
            r.Get("/", func(w http.ResponseWriter, r *http.Request) {
                page, limit := pagination(r)
                items, err := projSvc.List(r.Context(), authmw.UserID(r.Context()), page, limit)
                handleErr(w, err)
                if err != nil{
                 return
                }
                
                httpx.JSON(w, 200, map[string]any{
                    "projects": items,
                    "page":     page,
                    "limit":    limit,
                })
            })
            
            r.Post("/", func(w http.ResponseWriter, r *http.Request) {
                var req struct{ Name, Description string }
                if err := json.NewDecoder(r.Body).Decode(&req); err != nil { httpx.Error(w, 400, "invalid json"); return }
                if strings.TrimSpace(req.Name)=="" { httpx.ValidationError(w, map[string]string{"name":"is required"}); return }
                item, err := projSvc.Create(r.Context(), authmw.UserID(r.Context()), req.Name, req.Description)
                handleErr(w, err); if err != nil { return }
                httpx.JSON(w, 201, item)
            })
            r.Route("/{id}", func(r chi.Router) {
                r.Get("/", func(w http.ResponseWriter, r *http.Request) {
                    page, limit := pagination(r)
                    item, err := projSvc.Get(r.Context(), chi.URLParam(r, "id"), page, limit)
                    handleErr(w, err); if err != nil { return }
                    httpx.JSON(w, 200, item)
                })
                r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
                    var req struct{ Name, Description string }
                    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { httpx.Error(w, 400, "invalid json"); return }
                    if strings.TrimSpace(req.Name)=="" { httpx.ValidationError(w, map[string]string{"name":"is required"}); return }
                    item, err := projSvc.Update(r.Context(), authmw.UserID(r.Context()), chi.URLParam(r, "id"), req.Name, req.Description)
                    handleErr(w, err); if err != nil { return }
                    httpx.JSON(w, 200, item)
                })
                r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
                    err := projSvc.Delete(r.Context(), authmw.UserID(r.Context()), chi.URLParam(r, "id"))
                    handleErr(w, err); if err != nil { return }
                    w.WriteHeader(http.StatusNoContent)
                })
       
                r.Get("/tasks", func(w http.ResponseWriter, r *http.Request) {
                     page, limit := pagination(r)
                     tasks, err := taskSvc.List(
                         r.Context(),
                         authmw.UserID(r.Context()),
                         chi.URLParam(r, "id"),
                         r.URL.Query().Get("status"),
                         r.URL.Query().Get("assignee"),
                         page,
                         limit,
                        )
                        handleErr(w, err)
                        if err != nil {
                        return
                        }
                        httpx.JSON(w, 200, map[string]any{
                            "tasks": tasks,
                            "limit": limit,
                            "page":  page,
                        
                        })
                    })
                r.Post("/tasks", func(w http.ResponseWriter, r *http.Request) {
                    var req struct { Title, Description, Priority string; AssigneeID *string `json:"assignee_id"`; DueDate *string `json:"due_date"` }
                    if err := json.NewDecoder(r.Body).Decode(&req); err != nil { httpx.Error(w, 400, "invalid json"); return }
                    fields := map[string]string{}
                    if strings.TrimSpace(req.Title)=="" { fields["title"]="is required" }
                    if !validPriority(req.Priority) { fields["priority"]="must be one of low, medium, high" }
                    if len(fields)>0 { httpx.ValidationError(w, fields); return }
                    var due *time.Time
                    if req.DueDate != nil && *req.DueDate != "" {
                        parsed, err := time.Parse("2006-01-02", *req.DueDate); if err != nil { httpx.ValidationError(w, map[string]string{"due_date":"must be YYYY-MM-DD"}); return }
                        due = &parsed
                    }
                    item, err := taskSvc.Create(r.Context(), authmw.UserID(r.Context()), chi.URLParam(r, "id"), req.Title, req.Description, req.Priority, req.AssigneeID, due)
                    handleErr(w, err); if err != nil { return }
                    httpx.JSON(w, 201, item)
                })
                r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
                    stats, err := projSvc.Stats(r.Context(), authmw.UserID(r.Context()), chi.URLParam(r, "id"))
                    handleErr(w, err); if err != nil { return }
                    httpx.JSON(w, 200, stats)
                })
            })
        })
        r.Route("/tasks", func(r chi.Router) {
            r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {
                var req map[string]any
                if err := json.NewDecoder(r.Body).Decode(&req); err != nil { httpx.Error(w, 400, "invalid json"); return }
                fields := map[string]string{}
                updates := map[string]any{}
                if v, ok := req["title"].(string); ok { if strings.TrimSpace(v)=="" { fields["title"]="cannot be empty" } else { updates["title"]=v } }
                if v, ok := req["description"].(string); ok { updates["description"]=v }
                if v, ok := req["status"].(string); ok { if !validStatus(v) { fields["status"]="must be one of todo, in_progress, done" } else { updates["status"]=v } }
                if v, ok := req["priority"].(string); ok { if !validPriority(v) { fields["priority"]="must be one of low, medium, high" } else { updates["priority"]=v } }
                if v, ok := req["assignee_id"].(string); ok { updates["assignee_id"]=v }
                if v, ok := req["due_date"].(string); ok {
                    parsed, err := time.Parse("2006-01-02", v); if err != nil { fields["due_date"]="must be YYYY-MM-DD" } else { updates["due_date"]=parsed }
                }
                if len(fields)>0 { httpx.ValidationError(w, fields); return }
                item, err := taskSvc.Update(r.Context(), authmw.UserID(r.Context()), chi.URLParam(r, "id"), updates)
                handleErr(w, err); if err != nil { return }
                httpx.JSON(w, 200, item)
            })
            r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
                err := taskSvc.Delete(r.Context(), authmw.UserID(r.Context()), chi.URLParam(r, "id"))
                handleErr(w, err); if err != nil { return }
                w.WriteHeader(http.StatusNoContent)
            })
        })
    })
    return r
}

func pagination(r *http.Request) (int, int) {
    page, _ := strconv.Atoi(r.URL.Query().Get("page")); if page < 1 { page = 1 }
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit")); if limit < 1 || limit > 100 { limit = 20 }
    return page, limit
}
func validStatus(v string) bool { return v=="todo" || v=="in_progress" || v=="done" }
func validPriority(v string) bool { return v=="low" || v=="medium" || v=="high" }
func handleErr(w http.ResponseWriter, err error) {
    switch err {
    case nil:
        return
    case repository.ErrNotFound:
        httpx.Error(w, 404, "not found")
    case repository.ErrForbidden:
        httpx.Error(w, 403, "forbidden")
    case repository.ErrUnauthorized:
        httpx.Error(w, 401, "unauthorized")
    case repository.ErrConflict:
        httpx.Error(w, 400, "validation failed")
    default:
        httpx.Error(w, 500, "internal server error")
    }
}
