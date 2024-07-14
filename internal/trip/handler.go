package trip

import (
	"fmt"
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

func (h *Handler) CreateTrip(c echo.Context) error {
	var trip Trip
	if err := c.Bind(&trip); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, ok := c.Get("userID").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user ID from context")
	}
	trip.CreatedBy = userID

	if err := h.repo.Create(&trip); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, trip)
}

func (h *Handler) GetTrip(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	trip, err := h.repo.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Trip not found")
	}

	return c.JSON(http.StatusOK, trip)
}

func (h *Handler) UpdateTrip(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	var updatedTrip Trip
	if err := c.Bind(&updatedTrip); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingTrip, err := h.repo.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Trip not found")
	}

	existingTrip.Name = updatedTrip.Name
	existingTrip.Description = updatedTrip.Description
	existingTrip.StartDate = updatedTrip.StartDate
	existingTrip.EndDate = updatedTrip.EndDate

	if err := h.repo.Update(existingTrip); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	participants, err := h.repo.GetUsersForTrip(existingTrip.ID)
	if err != nil {
		log.Printf("Failed to get trip participants: %v", err)
	} else {
		message := fmt.Sprintf("The trip '%s' has been updated", existingTrip.Name)
		for _, participant := range participants {
			err = h.notificationService.SendNotification(participant.Email, notification.TripUpdate, message)
			if err != nil {
				log.Printf("Failed to send trip update notification to %s: %v", participant.Email, err)
			}
		}
	}

	return c.JSON(http.StatusOK, existingTrip)
}

func (h *Handler) DeleteTrip(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	if err := h.repo.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
