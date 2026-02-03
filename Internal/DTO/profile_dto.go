package  dto

type RegisterRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`

	Name         string `json:"name"`
	DobDay       string `json:"dobDay"`
	DobMonth     string `json:"dobMonth"`
	DobYear      string `json:"dobYear"`
	MotherTongue string `json:"motherTongue"`

	Gender         string `json:"gender"`
	Height         string `json:"height"`
	PhysicalStatus string `json:"physicalStatus"`
	MaritalStatus  string `json:"maritalStatus"`
	Religion       string `json:"religion"`

	Country      string `json:"country"`
	Employment   string `json:"employment"`
	Occupation   string `json:"occupation"`
	AnnualIncome int    `json:"annualIncome"`

	Star  string `json:"star"`
	Raasi string `json:"raasi"`

	Education    string `json:"education"`
	College      string `json:"college"`
	Organization string `json:"organization"`

	EatingHabit string `json:"eatingHabit"`
}
