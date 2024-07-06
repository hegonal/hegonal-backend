package notificationhandler

import (
	"fmt"

	"github.com/hegonal/hegonal-backend/app/models"
)

func ValidateNotification(notification *models.CreateNotification) error {
	switch notification.NotificationType {
	case models.DiscordWebhook:
		return ValidateDiscordWebhookConfig(notification.NotificationConfig)
	default:
		return fmt.Errorf("unsupported notification type")
	}
}