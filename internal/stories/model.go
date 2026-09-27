package stories

import "time"

type StoryResponse struct {
	ID                 string    `json:"id"`
	MediaURL           string    `json:"media_url"`
	MediaType          string    `json:"media_type"`
	DurationSeconds    int       `json:"duration_seconds"`
	ExpiresAt          time.Time `json:"expires_at"`
	CreatedAt          time.Time `json:"created_at"`
	IsViewed           bool      `json:"is_viewed"`
	CreatedBy          string    `json:"created_by,omitempty"`
	NameCreatedBy      string    `json:"name_created_by,omitempty"`
	TargetCountryCode  string    `json:"target_country_code,omitempty"`
}

type PublisherFeedResponse struct {
	PublisherID string          `json:"publisher_id"`
	EntityType  string          `json:"entity_type"`
	EntityID    string          `json:"entity_id"`
	DisplayName string          `json:"display_name"`
	AvatarURL   string          `json:"avatar_url"`
	IsVerified  bool            `json:"is_verified"`
	CountryCode string          `json:"country_code"`
	Stories     []StoryResponse `json:"stories"`
}

type ListByEntityDTO struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

type DeleteStoryDTO struct {
	ID         string `json:"id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}
