package services

import (
	models "marryo/Internal/Models"
	repositories "marryo/Internal/Repositories"

	"github.com/gofiber/fiber/v2"
)

type ChatService struct {
	repo repositories.Repository
}

func NewChatService(repo repositories.Repository) *ChatService {
	return &ChatService{repo: repo}
}

// Chat Message
func (s *ChatService) Chat(userID, otherUserID uint) bool {

	var match models.Match
	err := s.repo.FindOne(&match,
		"(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
		userID, otherUserID, otherUserID, userID)

	return err == nil
}

// SaveMessages
func (s *ChatService) SaveMessages(matchID, senderID uint, content string) error {

	msg := &models.Message{
		MatchID:  matchID,
		SenderID: senderID,
		Content:  content,
	}

	return s.repo.Create(&msg)
}

// List Chats
func (s *ChatService) GetChats(userID uint) ([]map[string]interface{}, error) {
	pg := s.repo.(*repositories.PgSQLRepository)

	var matches []models.Match

	// Self-healing: Ensure all accepted interests have a match record
	var acceptedInterests []models.Interest
	pg.DB.Where("(sender_id = ? OR receiver_id = ?) AND status = ?", userID, userID, "accepted").Find(&acceptedInterests)

	for _, inter := range acceptedInterests {
		var count int64
		pg.DB.Model(&models.Match{}).Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
			inter.SenderID, inter.ReceiverID, inter.ReceiverID, inter.SenderID).Count(&count)
		if count == 0 {
			newMatch := models.Match{
				User1ID: inter.SenderID,
				User2ID: inter.ReceiverID,
			}
			pg.DB.Create(&newMatch)
		}
	}

	if err := pg.DB.
		Where("user1_id = ? OR user2_id = ?", userID, userID).
		Find(&matches).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}

	for _, m := range matches {
		otherUserID := m.User1ID
		if otherUserID == userID {
			otherUserID = m.User2ID
		}

		var profile models.Profile
		pg.DB.
			Preload("Images", "is_primary = true").
			Where("user_id = ?", otherUserID).
			First(&profile)

		var lastMsg models.Message
		pg.DB.
			Where("match_id = ?", m.ID).
			Order("created_at DESC").
			First(&lastMsg)

		result = append(result, fiber.Map{
			"match_id": m.ID,
			"user": fiber.Map{
				"id":   otherUserID,
				"name": profile.FullName,
				"photo": func() string {
					if len(profile.Images) > 0 {
						return profile.Images[0].URL
					}
					return ""
				}(),
			},
			"last_message":    lastMsg.Content,
			"last_message_at": lastMsg.CreatedAt,
		})
	}

	return result, nil
}

// Get Chat Messages
func (s *ChatService) GetMessages(matchID uint, limit int) ([]models.Message, error) {
	pg := s.repo.(*repositories.PgSQLRepository)

	var messages []models.Message
	err := pg.DB.
		Where("match_id = ?", matchID).
		Order("created_at ASC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}
