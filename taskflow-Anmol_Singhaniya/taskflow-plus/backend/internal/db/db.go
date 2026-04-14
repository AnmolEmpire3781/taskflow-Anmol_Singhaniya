package db

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "path/filepath"
    "strings"

    "github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
    cfg, err := pgxpool.ParseConfig(databaseURL)
    if err != nil { return nil, err }
    cfg.MaxConns = 10
    return pgxpool.NewWithConfig(ctx, cfg)
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, log *slog.Logger, dir string) error {
    entries, err := os.ReadDir(dir)
    if err != nil { return err }
    for _, entry := range entries {
        if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".up.sql") { continue }
        sqlBytes, err := os.ReadFile(filepath.Join(dir, entry.Name()))
        if err != nil { return err }
        if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
            return fmt.Errorf("migration %s failed: %w", entry.Name(), err)
        }
        log.Info("migration applied", "file", entry.Name())
    }
    return nil
}

func RunSeed(ctx context.Context, pool *pgxpool.Pool, log *slog.Logger, path string) error {
    sqlBytes, err := os.ReadFile(path)
    if err != nil { return err }
    if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil { return err }
    log.Info("seed applied", "path", path)
    return nil
}
