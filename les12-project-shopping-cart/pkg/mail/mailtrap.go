package mail

import (
	"context"
	"net/http"
	"time"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/mailtrap/mailtrap-go"
	"github.com/rs/zerolog"
)

type MailtrapConfig struct {
	MailSender     string
	NameSender     string
	MailTrapApiKey string
}

type MailtrapProvider struct {
	client *http.Client
	config *MailtrapConfig
	logger *zerolog.Logger
}

func NewMailtrapProvider(config *MailConfig) (EmailProviderService, error) {
	mailtrapCfg, ok := config.ProviderConfig["mailtrap"].(map[string]any)
	if !ok {
		return nil, utils.NewError("Invalid or missing Mailtrap configuration", utils.ErrCodeInternal)
	}

	return &MailtrapProvider{
		client: &http.Client{Timeout: config.Timeout},
		config: &MailtrapConfig{
			MailSender:     mailtrapCfg["mail_sender"].(string),
			NameSender:     mailtrapCfg["name_sender"].(string),
			MailTrapApiKey: mailtrapCfg["mailtrap_api_key"].(string),
		},
		logger: config.Logger,
	}, nil
}

func (p *MailtrapProvider) SendMail(ctx context.Context, email *Email) error {
	traceID := logger.GetTraceID(ctx)
	start := time.Now()

	// Check if email is valid
	if len(email.To) == 0 {
		return utils.NewError("Email must have at least one recipient", utils.ErrCodeBadRequest)
	}

	// Set timeout
	if p.client.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.client.Timeout)
		defer cancel()
	}

	// Wait for 5 seconds
	time.Sleep(5 * time.Second)

	// Set sender
	email.From = Address{
		Email: p.config.MailSender,
		Name:  p.config.NameSender,
	}

	// Create a new Mailtrap client
	client, err := mailtrap.NewClient(p.config.MailTrapApiKey,
		mailtrap.WithSandbox(true),
		mailtrap.WithSandboxID(4082863),
	)

	if err != nil {
		p.logger.Error().Str("trace_id", traceID).
			Dur("duration", time.Since(start)).
			Str("operation", "create_mailtrap_client").
			Err(err).
			Msg("Failed to create mailtrap client")
		return utils.WrapError(err, "Failed to create mailtrap client", utils.ErrCodeInternal)
	}

	// Convert []Address -> []mailtrap.Address
	listEmail := make([]mailtrap.Address, 0, len(email.To))
	for _, addr := range email.To {
		listEmail = append(listEmail, mailtrap.Address(addr))
	}

	// Send the email using Mailtrap
	resp, _, err := client.Send(ctx, &mailtrap.SendRequest{
		From:     mailtrap.Address(email.From),
		To:       listEmail,
		Subject:  email.Subject,
		Text:     email.Text,
		Category: email.Category,
	})

	// Handle errors
	if err != nil {
		p.logger.Error().Str("trace_id", traceID).
			Dur("duration", time.Since(start)).
			Str("operation", "send_mail").
			Err(err).
			Msg("Failed to send email")
		return utils.WrapError(err, "Failed to send email", utils.ErrCodeInternal)
	}

	// Log success
	p.logger.Info().Str("trace_id", traceID).
		Dur("duration", time.Since(start)).
		Str("operation", "send_mail").
		Interface("message_ids", resp.MessageIDs).
		Msg("Email sent successfully")

	return nil
}
