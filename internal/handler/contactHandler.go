package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nakle1ka/Tramplin/internal/dto"
	"github.com/nakle1ka/Tramplin/internal/service"
)

type ContactHandler struct {
	contactService service.ContactService
}

func NewContactHandler(contactService service.ContactService) *ContactHandler {
	return &ContactHandler{
		contactService: contactService,
	}
}

func (h *ContactHandler) Create(c *gin.Context) {
	var req dto.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("failed to bind create contact request", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	recipientID, err := uuid.Parse(req.RecipientID)
	if err != nil {
		slog.Warn("invalid recipient_id format", "id", req.RecipientID)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid recipient_id"})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	srvDTO := service.CreateContactDTO{
		Auth:        auth,
		RecipientID: recipientID,
	}

	if err := h.contactService.Create(c.Request.Context(), srvDTO); err != nil {
		slog.Error("failed to create contact", "error", err, "sender_id", auth.UserID, "recipient_id", recipientID)
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	slog.Info("contact request created", "sender_id", auth.UserID, "recipient_id", recipientID)
	c.JSON(http.StatusCreated, gin.H{"message": "contact request created"})
}

func (h *ContactHandler) ListFriends(c *gin.Context) {
	var query dto.ListContactsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		slog.Warn("failed to bind list friends query", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	srvDTO := service.ListFriendsDTO{
		Auth:   auth,
		Limit:  query.Limit,
		Offset: query.Offset,
	}

	contacts, err := h.contactService.ListFriends(c.Request.Context(), srvDTO)
	if err != nil {
		slog.Error("failed to list friends", "error", err, "user_id", auth.UserID)
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := make([]dto.ContactResponse, len(contacts))
	for i, contact := range contacts {
		response[i] = dto.ContactResponse{
			ID:          contact.ID,
			SenderID:    contact.SenderID,
			RecipientID: contact.RecipientID,
			Status:      contact.Status,

			Sender: dto.Contactinfo{
				FirstName:  contact.Sender.FirstName,
				SecondName: contact.Sender.SecondName,
				LastName:   contact.Sender.LastName,
			},
			Recipient: dto.Contactinfo{
				FirstName:  contact.Recipient.FirstName,
				SecondName: contact.Recipient.SecondName,
				LastName:   contact.Recipient.LastName,
			},
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *ContactHandler) ListSentRequests(c *gin.Context) {
	var query dto.ListContactsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		slog.Warn("failed to bind list sent requests query", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	srvDTO := service.ListSentRequestsDTO{
		Auth:   auth,
		Limit:  query.Limit,
		Offset: query.Offset,
	}

	contacts, err := h.contactService.ListSentRequests(c.Request.Context(), srvDTO)
	if err != nil {
		slog.Error("failed to list sent requests", "error", err, "user_id", auth.UserID)
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := make([]dto.ContactResponse, len(contacts))
	for i, contact := range contacts {
		response[i] = dto.ContactResponse{
			ID:          contact.ID,
			SenderID:    contact.SenderID,
			RecipientID: contact.RecipientID,
			Status:      contact.Status,

			Sender: dto.Contactinfo{
				FirstName:  contact.Sender.FirstName,
				SecondName: contact.Sender.SecondName,
				LastName:   contact.Sender.LastName,
			},
			Recipient: dto.Contactinfo{
				FirstName:  contact.Recipient.FirstName,
				SecondName: contact.Recipient.SecondName,
				LastName:   contact.Recipient.LastName,
			},
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *ContactHandler) ListReceivedRequests(c *gin.Context) {
	var query dto.ListContactsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		slog.Warn("failed to bind list received requests query", "error", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	srvDTO := service.ListReceivedRequestsDTO{
		Auth:   auth,
		Limit:  query.Limit,
		Offset: query.Offset,
	}

	contacts, err := h.contactService.ListReceivedRequests(c.Request.Context(), srvDTO)
	if err != nil {
		slog.Error("failed to list received requests", "error", err, "user_id", auth.UserID)
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	response := make([]dto.ContactResponse, len(contacts))
	for i, contact := range contacts {
		response[i] = dto.ContactResponse{
			ID:          contact.ID,
			SenderID:    contact.SenderID,
			RecipientID: contact.RecipientID,
			Status:      contact.Status,

			Sender: dto.Contactinfo{
				FirstName:  contact.Sender.FirstName,
				SecondName: contact.Sender.SecondName,
				LastName:   contact.Sender.LastName,
			},
			Recipient: dto.Contactinfo{
				FirstName:  contact.Recipient.FirstName,
				SecondName: contact.Recipient.SecondName,
				LastName:   contact.Recipient.LastName,
			},
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *ContactHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	contactID, err := uuid.Parse(idStr)
	if err != nil {
		slog.Warn("invalid contact id format", "id", idStr)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid contact id"})
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || !req.Status.IsValid() {
		slog.Warn("invalid update status request", "error", err, "status", req.Status)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid input"})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	srvDTO := service.UpdateStatusDTO{
		Auth:   auth,
		ID:     contactID,
		Status: req.Status,
	}

	if err := h.contactService.UpdateStatus(c.Request.Context(), srvDTO); err != nil {
		slog.Error("failed to update contact status", "error", err, "contact_id", contactID, "status", req.Status)
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, service.ErrInvalidInput) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	slog.Info("contact status updated", "contact_id", contactID, "status", req.Status, "by_user", auth.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "contact status updated"})
}

func (h *ContactHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	contactID, err := uuid.Parse(idStr)
	if err != nil {
		slog.Warn("invalid contact id format", "id", idStr)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid contact id"})
		return
	}

	auth, err := extractAuthContext(c)
	if err != nil {
		slog.Error("auth context extraction failed", "error", err)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return
	}

	srvDTO := service.DeleteContactDTO{
		Auth: auth,
		ID:   contactID,
	}

	if err := h.contactService.Delete(c.Request.Context(), srvDTO); err != nil {
		slog.Error("failed to delete contact", "error", err, "contact_id", contactID)
		if errors.Is(err, service.ErrForbidden) {
			c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
			return
		}
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	slog.Info("contact deleted", "contact_id", contactID, "by_user", auth.UserID)
	c.JSON(http.StatusOK, gin.H{"message": "contact deleted"})
}
