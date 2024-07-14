package expense

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

func (h *Handler) CreateExpense(c echo.Context) error {
	var expense Expense
	if err := c.Bind(&expense); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}
	expense.TripID = tripID

	userID, ok := c.Get("userID").(int64)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get user ID from context")
	}
	expense.CreatedBy = userID

	if err := h.repo.Create(&expense); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, expense)
}

func (h *Handler) GetExpenses(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	expenses, err := h.repo.GetByTripID(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, expenses)
}

func (h *Handler) UpdateExpense(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("expenseId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid expense ID")
	}

	var updatedExpense Expense
	if err := c.Bind(&updatedExpense); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	existingExpense, err := h.repo.GetByID(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Expense not found")
	}

	existingExpense.Category = updatedExpense.Category
	existingExpense.Amount = updatedExpense.Amount
	existingExpense.Description = updatedExpense.Description
	existingExpense.Date = updatedExpense.Date

	if err := h.repo.Update(existingExpense); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, existingExpense)
}

func (h *Handler) DeleteExpense(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("expenseId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid expense ID")
	}

	if err := h.repo.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) GetBudgetSummary(c echo.Context) error {
	tripID, err := strconv.ParseInt(c.Param("tripId"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid trip ID")
	}

	summary, err := h.repo.GetBudgetSummary(tripID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, summary)
}
