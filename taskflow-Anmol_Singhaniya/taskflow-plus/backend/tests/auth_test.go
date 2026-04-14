package tests

import (
    "testing"

    "github.com/plus/taskflow/backend/internal/auth"
    "github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
    hash, err := auth.HashPassword("password123", 12)
    require.NoError(t, err)
    require.NoError(t, auth.ComparePassword(hash, "password123"))
}

func TestJWTGenerationAndParsing(t *testing.T) {
    token, err := auth.GenerateToken("u1", "test@example.com", "secret-key", 24)
    require.NoError(t, err)
    claims, err := auth.ParseToken(token, "secret-key")
    require.NoError(t, err)
    require.Equal(t, "u1", claims.UserID)
    require.Equal(t, "test@example.com", claims.Email)
}
