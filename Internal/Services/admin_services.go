package services

import (
	// "errors"
	"errors"
	"fmt"
	"time"

	config "marryo/Config"
	dto "marryo/Internal/DTO"
	models "marryo/Internal/Models"
	repositories "marryo/Internal/Repositories"
	utils "marryo/Internal/Utils"
)

type AdminService struct {
	repo repositories.Repository
}

func NewAdminService(repo repositories.Repository) *AdminService {
	return &AdminService{repo: repo}
}

func (s *AdminService) GetAllProfiles(userID uint, page, limit int) ([]models.User, int64, error) {

	offset := (page - 1) * limit

	pgRepo := s.repo.(*repositories.PgSQLRepository)

	db := pgRepo.DB.
		Model(&models.User{}).
		Preload("Profile").
		Where("id != ?", userID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var users []models.User
	err := db.
		Order("users.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	return users, total, err
}

//Block Unblock
func (s *AdminService) BlockUnblock(userID uint) (bool, error) {

	var user models.User
	if err := s.repo.FindByID(&user, userID, "id = ?"); err != nil {
		return false, errors.New("user not found")
	}

	status := !user.Status

	return status, s.repo.Update(&user,  map[string]interface{}{"status": status}, "id = ?", user.ID)
}

//Count Users
func (s *AdminService) GetActiveUsersCount() (int64, error) {
	return s.repo.CountWithCondition(&models.User{}, "status = ?", false)
}

//Count Matches
func (s *AdminService) GetUserMatchesCount() (int64, error) {

	return s.repo.CountWithCondition(&models.Match{}, "")
} 

//Delete Profile By ID
func (s *AdminService) DeleteProfile(userID uint) error {

	var user models.User

	err := s.repo.DeleteByID(&user, userID)
	if err != nil {
		return errors.New("profile deletion faild")
	}
	return nil
}

//GetMatches
func (s *AdminService) GetAllMatches() ([]models.Match, error) {

	pgRepo := s.repo.(*repositories.PgSQLRepository)

	var matches []models.Match

	err := pgRepo.DB.
		Preload("User1").
		Preload("User1.Profile").
		Preload("User1.Profile.Images").
		Preload("User2").
		Preload("User2.Profile").
		Preload("User2.Profile.Images").
		Order("matches.created_at DESC").
		Find(&matches).Error

	return matches, err
}

//Create User
func (s *AdminService) CreateUser(user *models.User) error {

	hashed, err := utils.Hashing(user.Password)
	if err != nil {
		return err
	}
	user.Password=hashed
	return s.repo.Create(&user)
}

//Edit Profile
func (s *AdminService) EditProfile(userID uint, input *dto.EditProfile) (*models.Profile, error) {

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

//EditAdminEmail and Password
func (s *AdminService) EditAdminProfile(userID uint, req *dto.AdminProfile) error {

	fields := make(map[string]interface{})

	if req.Username != "" {
		fields["username"] = req.Username
	}

	if req.Email != "" {
		fields["email"] = req.Email
	}

	if len(fields) == 0 {
		return nil 
	}

	return s.repo.Update(&models.User{}, fields, "id = ?", userID)
}

//ChangeAdminPassword
func (s *AdminService) ChangeAdminPassword(userID uint, req dto.AdminEditPassword) error {

	var user models.User
	if err := s.repo.FindOne(&user, "id = ?", userID); err != nil {
		return errors.New("user not found")
	}

	if err := utils.Comparepassword(user.Password, req.Oldpassword); err != nil {
		return errors.New("old password not matching")
	}

	if req.Newpassword != req.Confirmpassword {
		return errors.New("confirm password not match")
	}

	hashedpassword, err := utils.Hashing(req.Newpassword)
	if err != nil {
		return err
	}

	field := make(map[string]interface{})

	field["password"] = hashedpassword

	return s.repo.Update(&models.User{}, field, "id = ?", userID)
}

//ForgotPassword
func (s *AdminService) ForgotPassword(userID uint) error {

	var admin models.User
	if err := s.repo.FindOne(&admin, "id = ? AND role = ?", userID, "admin"); err != nil {
		return errors.New("admin not found")
	}

	otp := utils.GenerateOTP() 

	key := "admin_reset:" + admin.Email

	// Store OTP in Redis with 5 min expiry
	err := config.Redis.Set(config.Ctx, key, otp, 5*time.Minute).Err()
	if err != nil {
		return errors.New("failed to store OTP")
	}

	// Send OTP email
	utils.SendOTPEmail(admin.Email, otp)

	return nil
}
