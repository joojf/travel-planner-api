package main

import (
	"log"

	"github.com/joojf/travel-planner-api/config"
	"github.com/joojf/travel-planner-api/internal/activity"
	"github.com/joojf/travel-planner-api/internal/auth"
	"github.com/joojf/travel-planner-api/internal/database"
	"github.com/joojf/travel-planner-api/internal/destination"
	"github.com/joojf/travel-planner-api/internal/expense"
	"github.com/joojf/travel-planner-api/internal/invitation"
	"github.com/joojf/travel-planner-api/internal/itinerary"
	"github.com/joojf/travel-planner-api/internal/link"
	"github.com/joojf/travel-planner-api/internal/middleware"
	"github.com/joojf/travel-planner-api/internal/notification"
	"github.com/joojf/travel-planner-api/internal/review"
	"github.com/joojf/travel-planner-api/internal/trip"
	"github.com/joojf/travel-planner-api/internal/validator"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORS())

	db, err := database.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		e.Logger.Fatal(err)
	}

	if err := database.RunMigrations(db); err != nil {
		e.Logger.Fatal(err)
	}

	auth.InitJWT(cfg)

	emailService := notification.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPUsername,
		cfg.SMTPPassword,
		cfg.SMTPFromEmail,
	)
	notificationService := notification.NewService(emailService)

	authRepo := auth.NewSQLRepository(db)
	authHandler := auth.NewHandler(authRepo, notificationService)
	tripRepo := trip.NewRepository(db)
	tripHandler := trip.NewHandler(tripRepo, notificationService)
	activityRepo := activity.NewRepository(db)
	activityHandler := activity.NewHandler(activityRepo)
	invitationRepo := invitation.NewRepository(db)
	invitationHandler := invitation.NewHandler(invitationRepo, notificationService)
	destinationRepo := destination.NewRepository(db)
	destinationHandler := destination.NewHandler(destinationRepo)
	linkRepo := link.NewRepository(db)
	linkHandler := link.NewHandler(linkRepo)
	itineraryRepo := itinerary.NewRepository(db)
	itineraryHandler := itinerary.NewHandler(itineraryRepo)
	expenseRepo := expense.NewRepository(db)
	expenseHandler := expense.NewHandler(expenseRepo)
	reviewRepo := review.NewRepository(db)
	reviewHandler := review.NewHandler(reviewRepo)

	e.POST("/auth/register", authHandler.Register)
	e.POST("/auth/login", authHandler.Login)
	e.POST("/auth/reset-password", authHandler.ResetPassword)
	e.POST("/auth/set-new-password", authHandler.SetNewPassword)

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
	invGroup.DELETE("/:invitationId", invitationHandler.DeleteInvitation)

	// Activity routes
	actGroup := e.Group("/trips/:tripId/activities", middleware.AuthMiddleware)
	actGroup.POST("", activityHandler.CreateActivity)
	actGroup.GET("", activityHandler.GetActivities)
	actGroup.PUT("/:activityId", activityHandler.UpdateActivity)
	actGroup.DELETE("/:activityId", activityHandler.DeleteActivity)

	// Destination routes
	destGroup := e.Group("/trips/:tripId/destination", middleware.AuthMiddleware)
	destGroup.GET("", destinationHandler.GetDestination)
	destGroup.POST("", destinationHandler.CreateDestination)
	destGroup.PUT("", destinationHandler.UpdateDestination)
	destGroup.DELETE("", destinationHandler.DeleteDestination)

	// Link routes
	linkGroup := e.Group("/trips/:tripId/links", middleware.AuthMiddleware)
	linkGroup.POST("", linkHandler.CreateLink)
	linkGroup.GET("", linkHandler.GetLinks)
	linkGroup.PUT("/:linkId", linkHandler.UpdateLink)
	linkGroup.DELETE("/:linkId", linkHandler.DeleteLink)

	// Itinerary routes
	itineraryGroup := e.Group("/trips/:tripId/itineraries", middleware.AuthMiddleware)
	itineraryGroup.POST("", itineraryHandler.CreateItinerary)
	itineraryGroup.GET("", itineraryHandler.GetItineraries)
	itineraryGroup.PUT("/:itineraryId", itineraryHandler.UpdateItinerary)
	itineraryGroup.DELETE("/:itineraryId", itineraryHandler.DeleteItinerary)

	// Expense routes
	expenseGroup := e.Group("/trips/:tripId/expenses", middleware.AuthMiddleware)
	expenseGroup.POST("", expenseHandler.CreateExpense)
	expenseGroup.GET("", expenseHandler.GetExpenses)
	expenseGroup.PUT("/:expenseId", expenseHandler.UpdateExpense)
	expenseGroup.DELETE("/:expenseId", expenseHandler.DeleteExpense)
	expenseGroup.GET("/summary", expenseHandler.GetBudgetSummary)

	// Review routes
	reviewGroup := e.Group("/trips/:tripId/reviews", middleware.AuthMiddleware)
	reviewGroup.POST("", reviewHandler.CreateReview)
	reviewGroup.GET("", reviewHandler.GetReviews)
	reviewGroup.PUT("/:reviewId", reviewHandler.UpdateReview)
	reviewGroup.DELETE("/:reviewId", reviewHandler.DeleteReview)

	e.Logger.Fatal(e.StartTLS(":8080", "cert.pem", "key.pem"))
}
