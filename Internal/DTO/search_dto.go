package dto

type SearchFilterRequest struct {
	LookingFor   string   `json:"lookingFor"`
	MaritalStatus string  `json:"maritalStatus"`
	Religion     string   `json:"religion"`
	Caste        []string `json:"caste"`
	Education    string   `json:"education"`
	Occupation   string   `json:"occupation"`
	AgeFrom      int      `json:"ageFrom"`
	AgeTo        int      `json:"ageTo"`
	HeightFrom   string   `json:"heightFrom"`
	HeightTo     string   `json:"heightTo"`
	Star         string   `json:"star"`
	Country      string   `json:"country"`
	State        string   `json:"state"`
	City         string   `json:"city"`
}
