package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string, cost int) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(password), cost)
    return string(b), err
}

func ComparePassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
