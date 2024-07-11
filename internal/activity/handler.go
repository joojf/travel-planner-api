package activity

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

func (h *Handler) CreateActivity(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	var activity Activity
	if err := c.Bind(&activity); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	activity.TripID = tripID

	if err := h.repo.Create(&activity); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, activity)
}

func (h *Handler) GetActivities(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	activities, err := h.repo.GetByTripID(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, activities)
}

func (h *Handler) UpdateActivity(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("activityId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid activity ID")
	}

	existingActivity, err := h.repo.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Activity not found")
	}

	var updatedActivity Activity
	if err := c.Bind(&updatedActivity); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingActivity.Name = updatedActivity.Name
	existingActivity.Description = updatedActivity.Description
	existingActivity.Location = updatedActivity.Location
	existingActivity.StartTime = updatedActivity.StartTime
	existingActivity.EndTime = updatedActivity.EndTime

	if err := h.repo.Update(existingActivity); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingActivity)
}

func (h *Handler) DeleteActivity(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("activityId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid activity ID")
	}

	if err := h.repo.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
