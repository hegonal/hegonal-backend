package notificationhandler

import (
	"encoding/json"
	"time"

	"github.com/bensch777/discord-webhook-golang"
	"github.com/gofiber/fiber/v2/log"
	"github.com/hegonal/hegonal-backend/app/models"
	"github.com/hegonal/hegonal-backend/pkg/utils"
)

type DiscordWebhookConfig struct {
	WebhookURL string `json:"webhook_url" validate:"required,url"`
	Username   string `json:"username"`
	AvatarURL  string `json:"avatar_url"`
}

func ValidateDiscordWebhookConfig(config json.RawMessage) error {
	var discordConfig DiscordWebhookConfig
	if err := json.Unmarshal(config, &discordConfig); err != nil {
		log.Error(err)
		return err
	}
	validate := utils.NewValidator()
	return validate.Struct(discordConfig)
}

func sendDiscordWebhookNotification(httpMonitor models.HttpMonitor, incident models.Incident, notification models.Notification) error {
	var discordConfig DiscordWebhookConfig
	if err := json.Unmarshal(notification.NotificationConfig, &discordConfig); err != nil {
		log.Error(err)
		return err
	}
	var color int
	if incident.IncidentEnd == nil {
		color = 16711680
	} else {
		color = 65408
	}
	
    embeds := discordwebhook.Embed{
        Title:     httpMonitor.URL + "Is down!",
        Color:     color,
        Url:       "https://avatars.githubusercontent.com/u/6016509?s=48&v=4",
        Timestamp: time.Now(),
        Thumbnail: discordwebhook.Thumbnail{
            Url: "https://avatars.githubusercontent.com/u/6016509?s=48&v=4",
        },
        Author: discordwebhook.Author{
            Name:     "Author Name",
            Icon_URL: "https://avatars.githubusercontent.com/u/6016509?s=48&v=4",
        },
        Footer: discordwebhook.Footer{
            Text:     "Footer Text",
            Icon_url: "https://avatars.githubusercontent.com/u/6016509?s=48&v=4",
        },
    }

	hook := discordwebhook.Hook{
        Username:   "Captain Hook",
        Avatar_url: "https://avatars.githubusercontent.com/u/6016509?s=48&v=4",
        Content:    "@subscriber",
        Embeds:     []discordwebhook.Embed{embeds},
    }

    payload, err := json.Marshal(hook)

	if err != nil {
        log.Fatal(err)
    }

    err = discordwebhook.ExecuteWebhook(discordConfig.WebhookURL, payload)
    return err
}