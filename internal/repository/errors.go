package repository

import "errors"

var (
	ErrApplicantNotFound   = errors.New("applicant not found")
	ErrOpportunityNotFound = errors.New("opportunity not found")
	ErrEmployerNotFound    = errors.New("employer not found")
	ErrTagNotFound         = errors.New("tag not found")
	ErrTagNameExists       = errors.New("tag name already exists")
	ErrApplicationNotFound = errors.New("application not found")
)
