package handlers

import (
	"context"
	"net/http"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/notification/channels/sms"
	"github.com/zitadel/zitadel/internal/notification/channels/twilio"
	"github.com/zitadel/zitadel/internal/notification/channels/webhook"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// GetActiveSMSConfig reads the active iam sms provider config
func (n *NotificationQueries) GetActiveSMSConfig(ctx context.Context) (*sms.Config, error) {
	config, err := n.SMSProviderConfigActive(ctx, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		return nil, err
	}

	if config.TwilioConfig != nil {
		token, err := crypto.DecryptString(config.TwilioConfig.Token, n.SMSTokenCrypto)
		if err != nil {
			return nil, err
		}
		return &sms.Config{
			TwilioConfig: &twilio.Config{
				SID:          config.TwilioConfig.SID,
				Token:        token,
				SenderNumber: config.TwilioConfig.SenderNumber,
			},
		}, nil
	}
	if config.HTTPConfig != nil {
		return &sms.Config{
			WebhookConfig: &webhook.Config{
				CallURL: config.HTTPConfig.Endpoint,
				Method:  http.MethodPost,
				Headers: nil,
			},
		}, nil
	}

	return nil, zerrors.ThrowNotFound(nil, "HANDLER-8nfow", "Errors.SMS.Twilio.NotFound")
}
