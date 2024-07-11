package link

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

func (h *Handler) CreateLink(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	var link Link
	if err := c.Bind(&link); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	link.TripID = tripID

	if err := h.repo.Create(&link); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, link)
}

func (h *Handler) GetLinks(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	links, err := h.repo.GetByTripID(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, links)
}

func (h *Handler) UpdateLink(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("linkId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid link ID")
	}

	existingLink, err := h.repo.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Link not found")
	}

	var updatedLink Link
	if err := c.Bind(&updatedLink); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingLink.Title = updatedLink.Title
	existingLink.URL = updatedLink.URL
	existingLink.Description = updatedLink.Description

	if err := h.repo.Update(existingLink); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingLink)
}

func (h *Handler) DeleteLink(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("linkId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid link ID")
	}

	if err := h.repo.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}
