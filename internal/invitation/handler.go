package invitation

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo RepositoryInterface
}

func NewHandler(repo RepositoryInterface) *Handler {
	return &Handler{repo: repo}
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
	invitation.Status = "pending"

	if err := h.repo.Create(&invitation); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
