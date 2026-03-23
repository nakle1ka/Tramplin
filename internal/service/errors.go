package service

import (
	"errors"
)

var (
	ErrInvalidToken       = errors.New("invalid token")
	ErrCreateToken        = errors.New("failed to create tokens")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrUnknownRole        = errors.New("unknown user role")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("user with this email already exists")
	ErrInvalidEmployerINN = errors.New("invalid employer INN")
	ErrApplicantNotFound  = errors.New("applicant not found")
	ErrInvalidInput       = errors.New("invalid input")
)
