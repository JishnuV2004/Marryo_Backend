package services

import (
	"context"
	"log"
	models "marryo/Internal/Models"
	repositories "marryo/Internal/Repositories"

	"firebase.google.com/go/v4/messaging"
)

type NotificationService struct {
	repo repositories.Repository
	fcm  *messaging.Client
}

func NewNotificationService(
	repo repositories.Repository,
	fcm *messaging.Client,
) *NotificationService {
	return &NotificationService{
		repo: repo,
		fcm:  fcm,
	}
}

func (s *NotificationService) SendPush(userID uint, title string, body string, data map[string]string) {
	log.Printf("Starting SendPush for userID: %d, title: %s\n", userID, title)
	pg := s.repo.(*repositories.PgSQLRepository)

	var tokens []models.DeviceToken
	err := pg.DB.Where("user_id = ?", userID).Find(&tokens).Error
	if err != nil {
		log.Printf("Error fetching tokens for user %d: %v\n", userID, err)
		return
	}

	if len(tokens) == 0 {
		log.Printf("No device tokens found for userID: %d. Push skipped.\n", userID)
		return
	}

	log.Printf("Found %d tokens for userID: %d. Sending push...\n", len(tokens), userID)

	var fcmTokens []string
	for _, t := range tokens {
		fcmTokens = append(fcmTokens, t.Token)
	}

	msg := &messaging.MulticastMessage{
		Tokens: fcmTokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	}

	res, err := s.fcm.SendEachForMulticast(context.Background(), msg)
	if err != nil {
		log.Println("FCM send failed (fatal):", err)
		return
	}

	for i, r := range res.Responses {
		if !r.Success {
			log.Printf("FCM failed for token %s: %v\n", fcmTokens[i], r.Error)
		} else {
			log.Printf("FCM success for token %s, message ID: %s\n", fcmTokens[i], r.MessageID)
		}
	}
}
