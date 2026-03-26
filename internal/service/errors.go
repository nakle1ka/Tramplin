package service

import (
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrCreateToken        = errors.New("failed to create tokens")
	ErrUnknownRole        = errors.New("unknown user role")
	ErrForbidden          = errors.New("access denied")

	ErrUserNotFound      = errors.New("user not found")
	ErrApplicantNotFound = errors.New("applicant not found")
	ErrEmailExists       = errors.New("user with this email already exists")

	ErrInvalidInput       = errors.New("invalid input")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidEmployerINN = errors.New("invalid employer INN")

	ErrOpportunityNotFound = errors.New("opportunity not found")
	ErrInvalidOpportunity  = errors.New("invalid opportunity data")
	ErrInvalidDates        = errors.New("invalid dates: event start must be before event end")
	ErrInvalidSalary       = errors.New("invalid salary: min must be less than or equal to max")
	ErrExpiredOpportunity  = errors.New("cannot modify expired opportunity")
	ErrAlreadyModerated    = errors.New("opportunity already moderated")
)
