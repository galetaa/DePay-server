package models

type UpdateProfileRequest struct {
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Address     string `json:"address"`
}

type SubmitKYCRequest struct {
	DocumentType string `json:"document_type" binding:"required"`
	DocumentURL  string `json:"document_url" binding:"required"`
}
