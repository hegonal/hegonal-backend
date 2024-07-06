package notification

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2/log"
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