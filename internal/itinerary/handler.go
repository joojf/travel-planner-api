package itinerary

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo RepositoryInterface
}

func NewHandler(repo RepositoryInterface) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) CreateItinerary(c echo.Context) error {
	var itinerary Itinerary
	if err := c.Bind(&itinerary); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}
	itinerary.TripID = tripID

	userID, ok := c.Get("userID").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user ID from context")
	}
	itinerary.CreatedBy = userID

	if err := h.repo.Create(&itinerary); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, itinerary)
}

func (h *Handler) GetItineraries(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	itineraries, err := h.repo.GetByTripID(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, itineraries)
}

func (h *Handler) UpdateItinerary(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("itineraryId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid itinerary ID")
	}

	var updatedItinerary Itinerary
	if err := c.Bind(&updatedItinerary); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingItinerary, err := h.repo.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Itinerary not found")
	}

	existingItinerary.Title = updatedItinerary.Title
	existingItinerary.Description = updatedItinerary.Description
	existingItinerary.Date = updatedItinerary.Date

	if err := h.repo.Update(existingItinerary); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingItinerary)
}

func (h *Handler) DeleteItinerary(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("itineraryId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid itinerary ID")
	}

	if err := h.repo.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
