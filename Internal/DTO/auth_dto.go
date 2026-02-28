package dto

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	// Role string `json:"role"`
}

// type EditProfileRequest struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }
