package main

import (
	"github.com/joojf/travel-planner-api/internal/activity"
	"github.com/joojf/travel-planner-api/internal/auth"
	"github.com/joojf/travel-planner-api/internal/database"
	"github.com/joojf/travel-planner-api/internal/middleware"
	"github.com/joojf/travel-planner-api/internal/trip"
	"github.com/joojf/travel-planner-api/internal/validator"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	db, err := database.NewPostgresDB()
	if err != nil {
		e.Logger.Fatal(err)
	}

	if err := database.RunMigrations(db); err != nil {
		e.Logger.Fatal(err)
	}

	authRepo := auth.NewSQLRepository(db)
	authHandler := auth.NewHandler(authRepo)
	tripRepo := trip.NewRepository(db)
	tripHandler := trip.NewHandler(tripRepo)
	activityRepo := activity.NewRepository(db)
	activityHandler := activity.NewHandler(activityRepo)

	e.POST("/auth/register", authHandler.Register)
	e.POST("/auth/login", authHandler.Login)
	e.POST("/auth/reset-password", authHandler.ResetPassword)

	// Trip routes
	tripGroup := e.Group("/trips", middleware.AuthMiddleware)
	tripGroup.POST("", tripHandler.CreateTrip)
	tripGroup.GET("/:tripId", tripHandler.GetTrip)
	tripGroup.PUT("/:tripId", tripHandler.UpdateTrip)
	tripGroup.DELETE("/:tripId", tripHandler.DeleteTrip)

	// Invitation routes
	invGroup := e.Group("/trips/:tripId/invitations", middleware.AuthMiddleware)
	invGroup.POST("", invitationHandler.CreateInvitation)
	invGroup.GET("", invitationHandler.GetInvitations)

	// Activity routes
	actGroup := e.Group("/trips/:tripId/activities", middleware.AuthMiddleware)
	actGroup.POST("", activityHandler.CreateActivity)
	actGroup.GET("/:activityId", activityHandler.GetActivity)
	actGroup.PUT("/:activityId", activityHandler.UpdateActivity)
	actGroup.DELETE("/:activityId", activityHandler.DeleteActivity)

	// Expense routes
	// expGroup := e.Group("/trips/:tripId/expenses", middleware.AuthMiddleware)
	// expGroup.POST("", expenseHandler.CreateExpense)
	// expGroup.GET("", expenseHandler.GetExpenses)

	e.Logger.Fatal(e.Start(":8080"))
}
