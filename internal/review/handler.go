package review

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

func (h *Handler) CreateReview(c echo.Context) error {
	var review Review
	if err := c.Bind(&review); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}
	review.TripID = tripID

	userID, ok := c.Get("userID").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user ID from context")
	}
	review.UserID = userID

	if err := h.repo.Create(&review); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, review)
}

func (h *Handler) GetReviews(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	reviews, err := h.repo.GetByTripID(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, reviews)
}

func (h *Handler) UpdateReview(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("reviewId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid review ID")
	}

	var updatedReview Review
	if err := c.Bind(&updatedReview); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingReview, err := h.repo.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Review not found")
	}

	// TODO: Get user ID from authenticated session
	userID := int64(1) // Placeholder

	if existingReview.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, "You can only edit your own reviews")
	}

	existingReview.Rating = updatedReview.Rating
	existingReview.Comment = updatedReview.Comment

	if err := h.repo.Update(existingReview); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingReview)
}

func (h *Handler) DeleteReview(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("reviewId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid review ID")
	}

	// TODO: Get user ID from authenticated session
	userID := int64(1) // Placeholder

	if err := h.repo.Delete(id, userID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
