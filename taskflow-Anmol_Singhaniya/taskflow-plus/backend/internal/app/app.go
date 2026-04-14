package app

import (
    "context"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
    "time"

    "github.com/plus/taskflow/backend/internal/config"
    "github.com/plus/taskflow/backend/internal/db"
    "github.com/plus/taskflow/backend/internal/logger"
    "github.com/plus/taskflow/backend/internal/repository"
    "github.com/plus/taskflow/backend/internal/service"
)

func Run(ctx context.Context) error {
    cfg, err := config.Load()
    if err != nil { return err }
    log := logger.New(cfg.LogLevel)
    pool, err := db.NewPool(ctx, cfg.DatabaseURL)
    if err != nil { return err }
    defer pool.Close()

    migrationDir := filepath.Join("migrations")
    if _, err := os.Stat(migrationDir); err != nil { migrationDir = filepath.Join("backend", "migrations") }
    if err := db.RunMigrations(ctx, pool, log, migrationDir); err != nil { return err }

    if cfg.AutoSeed {
        seedPath := filepath.Join("seeds", "001_seed.sql")
        if _, err := os.Stat(seedPath); err != nil { seedPath = filepath.Join("backend", "seeds", "001_seed.sql") }
        if err := db.RunSeed(ctx, pool, log, seedPath); err != nil { return err }
    }

    users := repository.NewUserRepository(pool)
    projects := repository.NewProjectRepository(pool)
    tasks := repository.NewTaskRepository(pool)

    authSvc := service.NewAuthService(cfg, users)
    projSvc := service.NewProjectService(projects, tasks)
    taskSvc := service.NewTaskService(tasks, projects, users)

    srv := &http.Server{
        Addr: ":" + cfg.AppPort,
        Handler: Router(cfg, log, authSvc, projSvc, taskSvc),
        ReadHeaderTimeout: 10 * time.Second,
    }

    go func() {
        log.Info("server starting", slog.String("addr", srv.Addr))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            panic(err)
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop
    log.Info("shutdown signal received")
    shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
    defer cancel()
    if err := srv.Shutdown(shutdownCtx); err != nil { return fmt.Errorf("graceful shutdown failed: %w", err) }
    return nil
}
