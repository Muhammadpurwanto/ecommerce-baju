package dto

type UpdateUserRequest struct {
	Name      string  `json:"name" validate:"omitempty,min=2,max=255"`
	Phone     *string `json:"phone" validate:"omitempty,min=8,max=20"`
	AvatarURL *string `json:"avatar_url" validate:"omitempty,url"`
}
