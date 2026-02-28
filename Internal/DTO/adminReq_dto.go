package dto

type AdminProfile struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AdminEditPassword struct {
	Oldpassword     string `json:"oldpassword"`
	Newpassword     string `json:"newpassword"`
	Confirmpassword string `json:"confirmpassword"`
}
