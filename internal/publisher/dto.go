package publisher

type UpdatePublisherDTO struct {
	EntityType  string `json:"entity_type"`
	EntityID    string `json:"entity_id"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	IsVerified  bool   `json:"is_verified"`
	CountryCode string `json:"country_code"`
}

// CreatePlatformPublisherDTO crea un publisher de marca/plataforma (SUPER_ADMIN).
type CreatePlatformPublisherDTO struct {
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	IsVerified  bool   `json:"is_verified"`
}

type ListPublishersDTO struct {
	EntityType string `json:"entity_type"`
}
