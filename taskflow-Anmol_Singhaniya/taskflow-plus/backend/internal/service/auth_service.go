package service

import (
    "context"
    "strings"

    "github.com/plus/taskflow/backend/internal/auth"
    "github.com/plus/taskflow/backend/internal/config"
    "github.com/plus/taskflow/backend/internal/models"
    "github.com/plus/taskflow/backend/internal/repository"
)

type AuthService struct {
    cfg config.Config
    users *repository.UserRepository
}
func NewAuthService(cfg config.Config, users *repository.UserRepository) *AuthService { return &AuthService{cfg: cfg, users: users} }

type AuthResponse struct {
    Token string `json:"token"`
    User models.User `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (AuthResponse, error) {
    hash, err := auth.HashPassword(password, s.cfg.BcryptCost)
    if err != nil { return AuthResponse{}, err }
    user, err := s.users.Create(ctx, strings.TrimSpace(name), strings.ToLower(strings.TrimSpace(email)), hash)
    if err != nil { return AuthResponse{}, err }
    token, err := auth.GenerateToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
    if err != nil { return AuthResponse{}, err }
    return AuthResponse{Token: token, User: user}, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (AuthResponse, error) {
    user, err := s.users.GetByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
    if err != nil { return AuthResponse{}, repository.ErrUnauthorized }
    if err := auth.ComparePassword(user.Password, password); err != nil { return AuthResponse{}, repository.ErrUnauthorized }
    token, err := auth.GenerateToken(user.ID, user.Email, s.cfg.JWTSecret, s.cfg.JWTExpiryHours)
    if err != nil { return AuthResponse{}, err }
    user.Password = ""
    return AuthResponse{Token: token, User: user}, nil
}
