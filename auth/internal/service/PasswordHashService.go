package service

import "golang.org/x/crypto/bcrypt"

type PasswordHashService struct {
}

func NewPasswordHashService() *PasswordHashService {
	return &PasswordHashService{}
}

func (s *PasswordHashService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *PasswordHashService) CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
