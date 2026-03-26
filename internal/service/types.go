package service

import (
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
)

type AuthContext struct {
	UserID uuid.UUID
	Role   model.Role
}
