package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// type EditProfileRequest struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }
