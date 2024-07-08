package monitorengine

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	notificationhandler "github.com/hegonal/hegonal-backend/pkg/notification_handler"
	"github.com/hegonal/hegonal-backend/platform/database"
)

func sendNotifications(incident models.Incident,httpMonitor models.HttpMonitor, db *database.Queries) {
	httpMonitorNotifications,err := db.GetHttpMonitorNotifications(httpMonitor.HttpMonitorID) 
	if err != nil {
		log.Error(err)
		return
	}

	notificationIDArray := []string{}
	for _, notifications := range httpMonitorNotifications {
		notificationIDArray = append(notificationIDArray, notifications.NotificationID)
	}

	if len(notificationIDArray) == 0 {
		return
	}

	notifications, err:= db.GetAllMatchIDNotifications(notificationIDArray)
	if err != nil {
		log.Error(err)
		return
	}

	for _, notification := range notifications {
		notificationhandler.SendNotification(httpMonitor, incident,notification)
	}
}