package notificationhandler

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
)

func SendNotification(httpMonitor models.HttpMonitor, incident models.Incident, notification models.Notification) {
	switch notification.NotificationType {
	case models.DiscordWebhook:
		sendDiscordWebhookNotification(httpMonitor, incident, notification)
	default:
		log.Error(fmt.Errorf("error to find this notification type"))
	}
}