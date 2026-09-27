package publisher

import "time"

const EntityTypePlatform = "platform"

type Publisher struct {
	ID          string    `json:"id" db:"id"`
	EntityType  string    `json:"entity_type" db:"entity_type"`
	EntityID    string    `json:"entity_id" db:"entity_id"`
	DisplayName string    `json:"display_name" db:"display_name"`
	AvatarURL   string    `json:"avatar_url" db:"avatar_url"`
	IsVerified  bool      `json:"is_verified" db:"is_verified"`
	CountryCode string    `json:"country_code" db:"country_code"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
