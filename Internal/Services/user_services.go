package services

import (
	"errors"
	"fmt"
	dto "marryo/Internal/DTO"
	models "marryo/Internal/Models"
	repositories "marryo/Internal/Repositories"
	"strings"
	"time"

	"gorm.io/gorm"
	// utils "marryo/Internal/Utils"
)

type UserService struct {
	repo repositories.Repository
}

func NewUserService(repo repositories.Repository) *UserService {
	return &UserService{repo: repo}
}

// Profile
func (s *UserService) Profile(userID uint) (*models.Profile, error) {
	// var user models.User
	// err := s.repo.FindByID(&user, userID, "Profile.Images")

	var user models.Profile
	err := s.repo.FindByID(&user, userID, "user_id", "Images")
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// EditProfile
func (s *UserService) EditProfile(userID uint, input *dto.EditProfile) (*models.Profile, error) {

	var profile models.Profile

	if err := s.repo.FindOne(&profile, "user_id = ?", userID); err != nil {
		return nil, errors.New("profile not found")
	}

	if input.Name != "" {
		profile.FullName = input.Name
	}

	if input.DobDay != "" && input.DobMonth != "" && input.DobYear != "" {
		dob, err := time.Parse("2006-01-02",
			fmt.Sprintf("%s-%s-%s", input.DobYear, input.DobMonth, input.DobDay))
		if err == nil {
			profile.DOB = dob
		}
	}

	if input.MotherTongue != "" {
		profile.MotherTongue = input.MotherTongue
	}
	if input.Gender != "" {
		profile.Gender = input.Gender
	}
	if input.Height != "" {
		profile.Height = input.Height
	}
	if input.PhysicalStatus != "" {
		profile.PhysicalStatus = input.PhysicalStatus
	}
	if input.MaritalStatus != "" {
		profile.MaritalStatus = input.MaritalStatus
	}
	if input.Religion != "" {
		profile.Religion = input.Religion
	}
	if input.Country != "" {
		profile.Country = input.Country
	}
	if input.Employment != "" {
		profile.Employment = input.Employment
	}
	if input.Occupation != "" {
		profile.Occupation = input.Occupation
	}
	if input.AnnualIncome != 0 {
		profile.AnnualIncome = input.AnnualIncome
	}
	if input.Star != "" {
		profile.Star = input.Star
	}
	if input.Raasi != "" {
		profile.Raasi = input.Raasi
	}
	if input.Education != "" {
		profile.Education = input.Education
	}
	if input.College != "" {
		profile.College = input.College
	}
	if input.Organization != "" {
		profile.Organization = input.Organization
	}
	if input.EatingHabit != "" {
		profile.EatingHabit = input.EatingHabit
	}

	if err := s.repo.Save(&profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// SearchFilter
func (s *UserService) FilterProfiles(req *dto.SearchFilterRequest, userID uint) ([]models.Profile, error) {

	var profile models.Profile
	if err := s.repo.FindOne(&profile, "user_id = ?", userID); err != nil {
		return nil, errors.New("profile not found")
	}

	lookingFor := "female"
	if profile.Gender == "female" {
		lookingFor = "male"
	}

	pgRepo, ok := s.repo.(*repositories.PgSQLRepository)
	if !ok {
		return nil, errors.New("invalid repository")
	}

	db := pgRepo.DB.Model(&models.Profile{})

	db = db.
		Where("gender = ?", lookingFor).
		Where("profile_completed = true").
		Where("user_id != ?", userID)

	// if req.LookingFor != "" {
	// 	db.Where("gender = ?", req.LookingFor)
	// }
	if req.MaritalStatus != "" {
		db.Where("marital_status = ?", req.MaritalStatus)
	}
	if req.Religion != "" {
		db.Where("religion = ?", req.Religion)
	}
	if len(req.Caste) != 0 {
		db.Where("caste IN ?", req.Caste)
	}
	if req.Education != "" {
		db.Where("education = ?", req.Education)
	}
	if req.Occupation != "" {
		db.Where("occupation = ?", req.Occupation)
	}
	if req.Star != "" {
		db.Where("star = ?", req.Star)
	}
	if req.Country != "" {
		db.Where("country = ?", req.Country)
	}
	if req.State != "" {
		db.Where("state = ?", req.State)
	}
	if req.City != "" {
		db.Where("city = ?", req.City)
	}

	if req.AgeFrom > 0 && req.AgeTo > 0 {
		fromDOB := time.Now().AddDate(-req.AgeTo, 0, 0)
		toDOB := time.Now().AddDate(-req.AgeFrom, 0, 0)
		db.Where("dob BETWEEN ? AND ?", fromDOB, toDOB)
	}

	var profiles []models.Profile
	if err := db.Preload("Images").Find(&profiles).Error; err != nil {
		return nil, err
	}

	return profiles, nil
}

// FilterProfiles For User Home Page
func (s *UserService) HomeProfiles(userID uint) ([]models.Profile, error) {

	var profile models.Profile
	if err := s.repo.FindOne(&profile, "user_id = ?", userID); err != nil {
		return nil, errors.New("profile not found")
	}

	lookingFor := "female"
	if profile.Gender == "female" {
		lookingFor = "male"
	}

	pgRepo := s.repo.(*repositories.PgSQLRepository)
	db := pgRepo.DB.Model(&models.Profile{})

	db = db.
		Where("gender = ?", lookingFor).
		Where("profile_completed = true").
		Where("user_id != ?", userID)

	if !profile.DOB.IsZero() {
		now := time.Now()
		age := now.Year() - profile.DOB.Year()

		if now.YearDay() < profile.DOB.YearDay() {
			age--
		}

		minage := age - 5
		maxage := age + 5

		fromDOB := time.Now().AddDate(-maxage, 0, 0)
		toDOB := time.Now().AddDate(-minage, 0, 0)

		db = db.Where("dob IS NOT NULL").
			Where("dob BETWEEN ? AND ?", fromDOB, toDOB)
	}

	if profile.Religion != "" {
		db.Where("religion = ?", profile.Religion)
	}
	// if profile.Star != "" {
	// 	db.Where("star = ?", profile.Star)
	// }

	var matchedprofiles []models.Profile
	err := db.Preload("Images").Order("created_at DESC").
		// Limit(20).
		Find(&matchedprofiles).Error

	return matchedprofiles, err
}

// Delete
func (s *UserService) DeleteProfile(userID uint) error {

	var user models.User

	err := s.repo.DeleteByID(&user, userID)
	if err != nil {
		return errors.New("profile deletion faild")
	}
	return nil
}

// GetUserByID
func (s *UserService) GetProfile(userID uint) (*models.Profile, error) {
	var profile models.Profile
	err := s.repo.FindByID(&profile, userID, "user_id", "Images")
	if err != nil {
		return nil, errors.New("user profile not found")
	}

	return &profile, nil
}

// SearchProfiles
func (s *UserService) SearchProfiles(req *dto.SearchRequest, userID uint) ([]models.Profile, error) {

	var myProfile models.Profile
	if err := s.repo.FindByID(&myProfile, userID, "user_id = ?"); err != nil {
		return nil, errors.New("profile not found")
	}

	lookingFor := "female"
	if myProfile.Gender == "female" {
		lookingFor = "male"
	}

	pgRepo := s.repo.(*repositories.PgSQLRepository)
	db := pgRepo.DB.Model(&models.Profile{})

	db = db.Where("gender = ?", lookingFor).
		Where("profile_completed = true").
		Where("user_id != ?", userID)

	if req.Name != "" {
		db = db.Where(
			"LOWER(full_name) LIKE ?",
			"%"+strings.ToLower(req.Name)+"%",
		)
	}

	if req.Age > 0 {
		fromDOB := time.Now().AddDate(-(req.Age + 1), 0, 0)
		toDOB := time.Now().AddDate(-req.Age, 0, 0)
		db = db.Where("dob BETWEEN ? AND ?", fromDOB, toDOB)
	}

	if req.Star != "" {
		db = db.Where("star = ?", req.Star)
	}

	var profiles []models.Profile
	err := db.
		Order("created_at DESC").
		Limit(20).
		Find(&profiles).Error

	return profiles, err
}

// Send Interest
func (s *UserService) SendInterest(senderID, receiverID uint) error {

	if senderID == receiverID {
		return errors.New("cannot send interest to yourself")
	}

	var existing models.Interest
	err := s.repo.FindOne(&existing,
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID)
	if err == nil {
		return errors.New("interest alredy exists")
	}

	interest := models.Interest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     "pending",
	}

	return s.repo.Create(&interest)
}

// Accept Interest
func (s *UserService) AcceptInterest(receiverID, interestID uint) error {

	pgRepo := s.repo.(*repositories.PgSQLRepository)

	return pgRepo.DB.Transaction(func(tx *gorm.DB) error {

		var interest models.Interest
		if err := tx.First(&interest, interestID).Error; err != nil {
			return errors.New("interest not found")
		}

		if interest.ReceiverID != receiverID {
			return errors.New("unauthorized")
		}

		if err := tx.Model(&interest).
			Update("status", "accepted").Error; err != nil {
			return err
		}

		// Check if match already exists
		var existingMatch models.Match
		err := tx.Where("(user1_id = ? AND user2_id = ?) OR (user1_id = ? AND user2_id = ?)",
			interest.SenderID, interest.ReceiverID, interest.ReceiverID, interest.SenderID).First(&existingMatch).Error

		if err != nil {
			match := models.Match{
				User1ID: interest.ReceiverID,
				User2ID: interest.SenderID,
			}
			if err := tx.Create(&match).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Decline Interest
func (s *UserService) DeclineInterest(receiverID, interestID uint) error {

	var interest models.Interest
	if err := s.repo.FindByID(&interest, interestID, "id"); err != nil {
		return errors.New("interest not found")
	}

	if interest.ReceiverID != receiverID {
		return errors.New("unauthorized")
	}

	if interest.Status != "pending" {
		return errors.New("interest already pending")
	}

	interest.Status = "declined"
	return s.repo.Save(&interest)
}

// GetInterests
func (s *UserService) GetInterests(userID uint) ([]models.Interest, error) {
	pgRepo := s.repo.(*repositories.PgSQLRepository)

	var interests []models.Interest
	// err := pgRepo.DB.
	// 	Preload("Receiver").
	// 	Preload("Receiver.Profile").
	// 	Preload("Receiver.Profile.Images").
	// 	Where("sender_id = ?", userID).
	// 	Order("created_at DESC").
	// 	Find(&interests).Error
	err := pgRepo.DB.
		Preload("Sender").
		Preload("Sender.Profile").
		Preload("Sender.Profile.Images").
		Preload("Receiver").
		Preload("Receiver.Profile").
		Preload("Receiver.Profile.Images").
		Where("sender_id = ? AND status = ?", userID, "pending").
		Order("created_at DESC").
		Find(&interests).Error

	return interests, err
}

// GetReceivedInterests
func (s *UserService) GetReceivedInterests(userID uint) ([]models.Interest, error) {
	pgRepo := s.repo.(*repositories.PgSQLRepository)

	var interests []models.Interest
	// err := pgRepo.DB.
	// 	Preload("Sender.Profile.Images").
	// 	Where("receiver_id = ?", userID).
	// 	Order("created_at DESC").
	// 	Find(&interests).Error
	err := pgRepo.DB.
		Preload("Sender").
		Preload("Sender.Profile").
		Preload("Sender.Profile.Images").
		Preload("Receiver").
		Preload("Receiver.Profile").
		Preload("Receiver.Profile.Images").
		Where("receiver_id = ? AND status = ?", userID, "pending").
		Order("created_at DESC").
		Find(&interests).Error

	return interests, err
}

// GetAcceptedInterests
func (s *UserService) GetAcceptedInterests(userID uint) ([]models.Interest, error) {
	pgRepo := s.repo.(*repositories.PgSQLRepository)

	var interests []models.Interest
	// err := pgRepo.DB.
	// 	Preload("Sender.Profile.Images").
	// 	Preload("Receiver.Profile.Images").
	// 	Where(
	// 		"(sender_id = ? OR receiver_id = ?) AND status = ?",
	// 		userID, userID, "accepted",
	// 	).
	// 	Order("updated_at DESC").
	// 	Find(&interests).Error
	err := pgRepo.DB.
		Preload("Sender").
		Preload("Sender.Profile").
		Preload("Sender.Profile.Images").
		Preload("Receiver").
		Preload("Receiver.Profile").
		Preload("Receiver.Profile.Images").
		Where(
			"(sender_id = ? OR receiver_id = ?) AND status = ?",
			userID, userID, "accepted",
		).
		Order("updated_at DESC").
		Find(&interests).Error

	if err != nil {
		return nil, err
	}

	var matches []models.User
	for _, interest := range interests {
		if interest.SenderID == userID {
			matches = append(matches, interest.Receiver)
		} else {
			matches = append(matches, interest.Sender)
		}
	}

	return interests, err
}

// ReceivedInterests
func (s *UserService) ReceivedInterests(userID uint) ([]models.Interest, error) {

	pgRepo, ok := s.repo.(*repositories.PgSQLRepository)
	if !ok {
		return nil, errors.New("invalid repository")
	}

	var interests []models.Interest

	// Assuming we want pending interests. If we want all, remove Where status.
	// But usually "Requests" tab implies pending.
	// We'll return pending ones.
	err := pgRepo.DB.Model(&models.Interest{}).
		Preload("Sender").
		Preload("Sender.Profile").
		Preload("Sender.Profile.Images").
		Preload("Receiver").
		Preload("Receiver.Profile").
		Preload("Receiver.Profile.Images").
		Where("receiver_id = ? AND status = ?", userID, "pending").
		Order("created_at DESC").
		Find(&interests).Error

	if err != nil {
		return nil, err
	}

	return interests, nil
}

func (s *UserService) SaveDeviceToken(token models.DeviceToken) error {

	pg := s.repo.(*repositories.PgSQLRepository)

	return pg.DB.
		Where("token = ?", token.Token).
		Assign(models.DeviceToken{
			UserID:   token.UserID,
			Platform: token.Platform,
		}).
		FirstOrCreate(&models.DeviceToken{}).Error
}



//UploadIMG
// func (s *UserService) UploadIMG(userID uint, ImgURL string) error {

// 	var profile models.Profile
// 	if err := s.repo.FindOne(&profile, "user_id = ?", userID); err != nil {
// 		return errors.New("profile not found")
// 	}

// 	pgRepo := s.repo.(*repositories.PgSQLRepository)

// 	var count int64
// 	err := pgRepo.DB.Model(&models.Img{}).
// 		Where("profile_id = ? AND deleted_at IS NULL", userID).
// 		Count(&count).Error;
// 	if err != nil {
// 		return err
// 	}

// 	if count >= 3 {
// 		return errors.New("maximum 3 images allowed")
// 	}

// 	img := models.Img{
// 		ProfileID: userID,
// 		URL: ImgURL,
// 		IsApproved: false,
// 	}

// 	return s.repo.Create(&img)
// }
