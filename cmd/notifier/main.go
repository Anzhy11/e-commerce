package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/models"
	"github.com/anzhy11/go-e-commerce/internal/notifications"
	"github.com/anzhy11/go-e-commerce/internal/providers"
)

func main() {
	log.Println("Starting notifier service...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	notifierConfig := &notifications.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}

	emailNotifier := notifications.NewEmailNotifier(notifierConfig)

	awsConfig, err := providers.CreateAwsConfig(ctx, &cfg.AWS)
	if err != nil {
		log.Printf("Failed to create AWS config: %v", err)
	}

	logger := watermill.NewStdLogger(false, false)

	subscriber, err := sqs.NewSubscriber(
		sqs.SubscriberConfig{
			AWSConfig: awsConfig,
		},
		logger,
	)
	if err != nil {
		log.Printf("Failed to create SQS subscriber: %v", err)
	}
	defer func() {
		_ = subscriber.Close()
	}()

	messages, err := subscriber.Subscribe(ctx, cfg.AWS.EventQueueName)
	if err != nil {
		log.Printf("Failed to subscribe to SQS: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	log.Println("Notifier service started. Waiting for messages...")

	for {
		select {
		case msg, ok := <-messages:
			if !ok {
				log.Println("SQS message channel closed")
				return
			}

			if err := processMessage(msg, emailNotifier); err != nil {
				log.Printf("Failed to process message: %v", err)
				msg.Nack()
				continue
			}

			msg.Ack()

		case <-signalChan:
			log.Println("Shutting down notifier service...")
			cancel()
			return
		}
	}
}

func processMessage(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	eventType := msg.Metadata.Get("event_type")

	switch eventType {
	case notifications.UserLoggedInEventType:
		return handleUserLoggedIn(msg, emailNotifier)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func handleUserLoggedIn(msg *message.Message, emailNotifier *notifications.EmailNotifier) error {
	var user models.User
	if err := json.Unmarshal(msg.Payload, &user); err != nil {
		return err
	}

	userName := user.FirstName + " " + user.LastName
	if userName == "" {
		userName = "User"
	}

	log.Printf("Sending login notification to %s (%s)", userName, user.Email)

	return emailNotifier.SendLoginNotification(user.Email, userName)
}
