package invitation

import (
	"fmt"
	"github.com/joojf/travel-planner-api/internal/trip"
	"log"
	"net/http"
	"strconv"

	"github.com/joojf/travel-planner-api/internal/notification"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo                RepositoryInterface
	notificationService *notification.Service
}

func NewHandler(repo RepositoryInterface, notificationService *notification.Service) *Handler {
	return &Handler{
		repo:                repo,
		notificationService: notificationService,
	}
}

func (h *Handler) CreateInvitation(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	var invitation Invitation
	if err := c.Bind(&invitation); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	invitation.TripID = tripID

	if err := h.repo.Create(&invitation); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Fetch trip name for the notification message
	tripDetails, err := h.repo.GetTripByID(tripID)
	if err != nil {
		log.Printf("Failed to get trip details: %v", err)
		tripDetails = &trip.Trip{Name: "Unknown"} // Using trip.Trip, assuming Trip is defined in the trip package
	}

	message := fmt.Sprintf("You have been invited to join the trip '%s'", tripDetails.Name)
	err = h.notificationService.SendNotification(invitation.Email, notification.TripInvitation, message)
	if err != nil {
		log.Printf("Failed to send invitation notification: %v", err)
		// Note: We're not returning an error here, as the invitation was created successfully
		// TODO: Return an error
	}

	return c.JSON(http.StatusCreated, invitation)
}

func (h *Handler) GetInvitations(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	invitations, err := h.repo.GetByTripID(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, invitations)
}

func (h *Handler) DeleteInvitation(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("invitationId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid invitation ID")
	}

	if err := h.repo.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
