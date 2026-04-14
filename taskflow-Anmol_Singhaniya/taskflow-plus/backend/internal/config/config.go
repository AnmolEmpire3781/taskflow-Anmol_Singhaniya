package config

import (
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
)

type Config struct {
    AppEnv         string
    AppPort        string
    LogLevel       string
    DatabaseURL    string
    JWTSecret      string
    JWTExpiryHours int
    BcryptCost     int
    AutoSeed       bool
    ShutdownTimeout time.Duration
}

func Load() (Config, error) {
    _ = godotenv.Load()
    cfg := Config{
        AppEnv: getEnv("APP_ENV", "development"),
        AppPort: getEnv("APP_PORT", "8080"),
        LogLevel: getEnv("LOG_LEVEL", "INFO"),
        DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@db:5432/taskflow?sslmode=disable"),
        JWTSecret: os.Getenv("JWT_SECRET"),
        JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 24),
        BcryptCost: getEnvInt("BCRYPT_COST", 12),
        AutoSeed: getEnvBool("AUTO_SEED", true),
        ShutdownTimeout: 10 * time.Second,
    }
    if cfg.JWTSecret == "" {
        return cfg, fmt.Errorf("JWT_SECRET is required")
    }
    if cfg.BcryptCost < 12 {
        return cfg, fmt.Errorf("BCRYPT_COST must be >= 12")
    }
    return cfg, nil
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" { return v }
    return fallback
}
func getEnvInt(key string, fallback int) int {
    if v := os.Getenv(key); v != "" {
        if i, err := strconv.Atoi(v); err == nil { return i }
    }
    return fallback
}
func getEnvBool(key string, fallback bool) bool {
    if v := os.Getenv(key); v != "" {
        if b, err := strconv.ParseBool(v); err == nil { return b }
    }
    return fallback
}
