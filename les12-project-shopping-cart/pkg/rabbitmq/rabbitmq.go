package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

type rabbitMQService struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	logger  *zerolog.Logger
}

func NewRabbitMQService(amqpURL string, logger *zerolog.Logger) (RabbitMQService, error) {
	// Connect to RabbitMQ server
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to connect to RabbitMQ")
		return nil, err
	}

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		logger.Error().Err(err).Msg("Failed to open channel")
		return nil, err
	}

	return &rabbitMQService{
		conn:    conn,
		channel: ch,
		logger:  logger,
	}, nil
}

func (r *rabbitMQService) Publish(ctx context.Context, queue string, message any) error {
	// Declare a queue
	_, err := r.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to declare queue")
		return err
	}

	// Convert the message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to parse message")
		return err
	}

	// Publish the message
	err = r.channel.PublishWithContext(ctx, "", queue, false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})

	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to publish message")
		return err
	}

	return nil
}

func (r *rabbitMQService) Consume(ctx context.Context, queue string, handler func([]byte) error) error {
	// Declare a queue
	_, err := r.channel.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to declare queue")
		return err
	}

	// Start consuming messages
	msgs, err := r.channel.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to declare consume")
		return err
	}


	// Start a goroutine to handle messages
	go func() {
		for {
			select {
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				if err := handler(msg.Body); err != nil {
					msg.Nack(false, false)
				} else {
					msg.Ack(false)
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *rabbitMQService) Close() error {
	// Close the channel and connection
	if r.channel != nil {
		if err := r.channel.Close(); err != nil {
			r.logger.Error().Err(err).Msg("Failed to close channel")
			return err
		}
	}

	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			r.logger.Error().Err(err).Msg("Failed to close connection")
			return err
		}
	}

	return nil
}
