package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"

	appConfig "github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/providers"
	_ "github.com/aws/smithy-go/endpoints"
)

type EventPublisher struct {
	publisher message.Publisher
	queueName string
}

func NewEventPublisher(ctx context.Context, cfg *appConfig.AWSConfig) (*EventPublisher, error) {
	logger := watermill.NewStdLogger(false, false)
	awsConfig, err := providers.CreateAwsConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create aws config: %w", err)
	}

	publisherConfig := sqs.PublisherConfig{
		AWSConfig: awsConfig,
		Marshaler: nil,
	}

	publisher, err := sqs.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return &EventPublisher{
		publisher: publisher,
		queueName: cfg.EventQueueName,
	}, nil
}

func (e *EventPublisher) Publish(eventType string, payload any, metadata map[string]string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), data)
	msg.Metadata.Set("event_type", eventType)
	for k, v := range metadata {
		msg.Metadata.Set(k, v)
	}

	return e.publisher.Publish(e.queueName, msg)
}

func (e *EventPublisher) Close() error {
	return e.publisher.Close()
}
