package destination

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

func (h *Handler) GetDestination(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	destination, err := h.repo.GetByTripID(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Destination not found")
	}

	return c.JSON(http.StatusOK, destination)
}

func (h *Handler) CreateDestination(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	var destination Destination
	if err := c.Bind(&destination); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	destination.TripID = tripID

	if err := h.repo.Create(&destination); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, destination)
}

func (h *Handler) UpdateDestination(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	existingDestination, err := h.repo.GetByTripID(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Destination not found")
	}

	var updatedDestination Destination
	if err := c.Bind(&updatedDestination); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingDestination.Name = updatedDestination.Name
	existingDestination.Country = updatedDestination.Country
	existingDestination.City = updatedDestination.City
	existingDestination.Description = updatedDestination.Description

	if err := h.repo.Update(existingDestination); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingDestination)
}

func (h *Handler) DeleteDestination(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	if err := h.repo.Delete(tripID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
