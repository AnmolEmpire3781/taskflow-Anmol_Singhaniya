package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/plus/taskflow/backend/internal/auth"
    "github.com/plus/taskflow/backend/internal/config"
    "github.com/plus/taskflow/backend/internal/httpx"
)

type ctxKey string
const UserIDKey ctxKey = "user_id"
const UserEmailKey ctxKey = "user_email"

func Auth(cfg config.Config) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            header := r.Header.Get("Authorization")
            if !strings.HasPrefix(header, "Bearer ") { httpx.Error(w, http.StatusUnauthorized, "unauthorized"); return }
            token := strings.TrimPrefix(header, "Bearer ")
            claims, err := auth.ParseToken(token, cfg.JWTSecret)
            if err != nil { httpx.Error(w, http.StatusUnauthorized, "unauthorized"); return }
            ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
            ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func UserID(ctx context.Context) string {
    v, _ := ctx.Value(UserIDKey).(string)
    return v
}
