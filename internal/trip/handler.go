package trip

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

func (h *Handler) CreateTrip(c echo.Context) error {
	var trip Trip
	if err := c.Bind(&trip); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO: Get user ID from authenticated session
	trip.CreatedBy = 1 // Placeholder

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
