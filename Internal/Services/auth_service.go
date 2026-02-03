package services

import (
	"errors"
	"fmt"

	"time"

	config "marryo/Config"
	dto "marryo/Internal/DTO"
	models "marryo/Internal/Models"
	repositories "marryo/Internal/Repositories"
	utils "marryo/Internal/Utils"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	repo  repositories.Repository
	redis *redis.Client
}

func NewAuthService(repo repositories.Repository, redis *redis.Client) *AuthService {
	return &AuthService{repo: repo, redis: redis}
}

//Signup
func (s *AuthService) Signup(u *dto.RegisterRequest) (interface{}, error) {

	var existing models.User
	if err := s.repo.FindOne(&existing, "email = ?", u.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	hashed, err := utils.Hashing(u.Password)
	if err != nil {
		return nil, err
	}

	otp := utils.GenerateOTP()

	// user := models.User{
	// 	Email:      u.Email,
	// 	Password:   hashed,
	// 	IsVerified: false,
	// }

	s.redis.HSet(config.Ctx, "signup:"+u.Email, map[string]interface{}{ "email":    u.Email, "password": hashed,},)

	s.redis.Expire(config.Ctx, "signup:"+u.Email, 10*time.Minute)

	s.redis.Set(config.Ctx, "otp:"+u.Email, otp, 5*time.Minute)

	if err := utils.SendOTPEmail(u.Email, otp); err != nil {
		return nil, errors.New("failed to send OTP email")
	}

	fmt.Println("OTP:", otp)

	return fiber.Map{
		"message": "OTP sent to email. Please verify to continue.",
	}, nil

}

//Verifiy OTP
func (s *AuthService) VerifiyOTP(email, otp string) error {

	stroedOTP, err := s.redis.Get(config.Ctx, "otp:"+email).Result()
	if err != nil || stroedOTP != otp {
		return errors.New("invalid or expired OTP")
	}

	data, err := s.redis.HGetAll(config.Ctx, "signup:"+email).Result()
	if err != nil || len(data) == 0 {
		return errors.New("signup session expired")
	}

	user := models.User{
		Email: data["email"],
		Password: data["password"],
		IsVerified: true,
	}
	if err := s.repo.Create(&user); err != nil {
		return err
	}

	s.redis.Del(config.Ctx, "otp:"+email)
	s.redis.Del(config.Ctx, "signup:"+email)

	return nil
}

// CompleteSignup
func (s *AuthService) CompleteSignup(u *dto.RegisterRequest) ( error) {

	var user models.User

	if err := s.repo.FindOne(&user, "email = ?", u.Email); err != nil {
		return errors.New("user not found")
	}

	if !user.IsVerified {
		return errors.New("email not verified")
	}

	dob, err := time.Parse(
		"2006-1-2",
		fmt.Sprintf("%s-%s-%s", u.DobYear, u.DobMonth, u.DobDay),
	)

		profile := models.Profile{
			UserID: user.ID,
			FullName: u.Name,
			DOB: dob,
			MotherTongue: u.MotherTongue,

			Gender: u.Gender,
			Height: u.Height,
			PhysicalStatus: u.PhysicalStatus,
			MaritalStatus:  u.MaritalStatus,
			Religion:       u.Religion,

			Country:      u.Country,
			Employment:   u.Employment,
			Occupation:   u.Occupation,
			AnnualIncome: u.AnnualIncome,

			Star:  u.Star,
			Raasi: u.Raasi,

			Education:    u.Education,
			College:      u.College,
			Organization: u.Organization,

			EatingHabit: u.EatingHabit,

			ProfileCompleted: true,
		}
		err = s.repo.Create(&profile)
		return err
}

// Login
func (s *AuthService) Login(data *dto.LoginRequest) (*models.User, string, string, error) {
	var user models.User
	err := s.repo.FindOne(&user, "email = ?", data.Email)
	if err != nil {
		return nil, "", "", errors.New("user not found")
	}

	if err := utils.Comparepassword(user.Password, data.Password); err != nil {
		return nil, "", "", errors.New("invalid password")
	}

	access, _ := utils.GenerateAccess(user.ID)
	refresh, _ := utils.GenerateRefresh(user.ID)

	s.redis.Set(config.Ctx, "refresh"+refresh, user.ID, 7*24*time.Hour)

	return &user, access, refresh, nil
}

// Refresh Func
func (s *AuthService) Refresh(oldRefresh string) (string, string, error) {

	id, err := s.redis.Get(config.Ctx, "refresh"+oldRefresh).Uint64()
	if err != nil {
		return "", "", errors.New("invalid or expired refresh token")
	}

	userID := uint(id)

	s.redis.Del(config.Ctx, "refresh"+oldRefresh)

	newAccess, _ := utils.GenerateAccess(userID)
	newRefresh, _ := utils.GenerateRefresh(userID)

	s.redis.Set(config.Ctx, "refresh"+newRefresh, userID, 7*24*time.Hour)

	return newAccess, newRefresh, nil
}

// Logout
func (s *AuthService) Logout(access, refresh string) {
	s.redis.Del(config.Ctx, "refresh"+refresh)
	s.redis.Set(config.Ctx, "blacklist:"+access, "1", 15*time.Minute)
}

//Profile
// func (s *AuthService) Profile(email string) (*models.User, error) {
// 	var user models.User
// 	err := s.repo.FindOne(&user, "email = ?", email)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }
