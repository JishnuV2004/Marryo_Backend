package services

import (
	"errors"
	"fmt"
	dto "marryo/Internal/DTO"
	models "marryo/Internal/Models"
	repositories "marryo/Internal/Repositories"
	"time"
	// utils "marryo/Internal/Utils"
)

type UserService struct {
	repo repositories.Repository
}

func NewUserService(repo repositories.Repository) *UserService {
	return &UserService{repo: repo}
}

// Profile
func (s *UserService) Profile(userID uint) (*models.User, error) {
	var user models.User
	err := s.repo.FindByID(&user, userID, "Profile")
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
func (s *UserService) FilterProfiles(req *dto.SearchFilterRequest) ([]models.Profile, error) {

	pgRepo, ok := s.repo.(*repositories.PgSQLRepository)
	if !ok {
		return nil, errors.New("invalid repository")
	}

	db := pgRepo.DB.Model(&models.Profile{})

	if req.LookingFor != "" {
		db.Where("gender = ?", req.LookingFor)
	}
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
	if err := db.Find(&profiles).Error; err != nil {
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
	err := db.Order("created_at DESC").
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
