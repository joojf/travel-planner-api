package invitation

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateInvitation(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	var inv Invitation
	if err := c.Bind(&inv); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	inv.TripID = tripID
	if err := h.repo.Create(&inv); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create invitation")
	}

	return c.JSON(http.StatusCreated, inv)
}

func (h *Handler) GetInvitations(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	invitations, err := h.repo.GetByTripID(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch invitations")
	}

	return c.JSON(http.StatusOK, invitations)
}

// TODO: Implement other handler methods (UpdateInvitation, DeleteInvitation, etc.)
