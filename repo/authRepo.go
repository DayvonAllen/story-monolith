package repo

import "story-app-monolith/domain"

type AuthRepo interface {
	Login(username string, password string) (*domain.UserDto, string, error)
	ResetPassword(token, password string) error
	ResetPasswordQuery(email string) error
	VerifyCode(code string) error
}

