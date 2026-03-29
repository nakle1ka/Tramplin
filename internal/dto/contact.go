package dto

import (
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/model"
)

type CreateContactRequest struct {
	RecipientID string `json:"recipient_id" binding:"required,uuid"`
}

type ListContactsQuery struct {
	Limit  int `form:"limit" binding:"min=1,max=100"`
	Offset int `form:"offset" binding:"min=0"`
}

type UpdateStatusRequest struct {
	Status model.ContactStatus `json:"status"`
}

type ContactResponse struct {
	ID     uuid.UUID           `json:"id"`
	Status model.ContactStatus `json:"status"`

	SenderID uuid.UUID   `json:"sender_id"`
	Sender   Contactinfo `json:"sender"`

	RecipientID uuid.UUID   `json:"recipient_id"`
	Recipient   Contactinfo `json:"recipient"`
}

type Contactinfo struct {
	FirstName  string `json:"first_name"`
	SecondName string `json:"second_name"`
	LastName   string `json:"last_name"`
}
