package main

import (
	"fmt"
	"log"

	"github.com/anzhy11/go-e-commerce/internal/config"
	"github.com/anzhy11/go-e-commerce/internal/notifications"
)

func main() {
	cfg, _ := config.Load()

	notifierConfig := &notifications.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}

	notifier := notifications.NewEmailNotifier(notifierConfig)

	err := notifier.SendLoginNotification(cfg.SMTP.From, "Aadel")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Email mocroservice started")
}
