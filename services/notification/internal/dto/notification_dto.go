package dto

type SendEmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"` // Mendukung HTML
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
