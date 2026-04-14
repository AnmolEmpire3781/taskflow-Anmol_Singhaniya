package repository

import "strings"

func isUniqueViolation(err error) bool {
    if err == nil { return false }
    return strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "SQLSTATE 23505")
}
