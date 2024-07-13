package notification

import (
	"log"
)

type NotificationType string

const (
	TripUpdate     NotificationType = "trip_update"
	TripInvitation NotificationType = "trip_invitation"
	TripReminder   NotificationType = "trip_reminder"
)

type Service struct {
	emailService *EmailService
}

func NewService(emailService *EmailService) *Service {
	return &Service{
		emailService: emailService,
	}
}

func (s *Service) SendNotification(to string, notificationType NotificationType, message string) error {
	subject := s.getSubjectForNotificationType(notificationType)
	err := s.emailService.SendEmail(to, subject, message)
	if err != nil {
		log.Printf("Failed to send email notification: %v", err)
		return err
	}
	return nil
}

func (s *Service) getSubjectForNotificationType(notificationType NotificationType) string {
	switch notificationType {
	case TripUpdate:
		return "Trip Update"
	case TripInvitation:
		return "New Trip Invitation"
	case TripReminder:
		return "Trip Reminder"
	default:
		return "Travel Planner Notification"
	}
}
