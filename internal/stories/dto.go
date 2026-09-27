package stories

type CreateStoryDTO struct {
	EntityType         string `json:"entity_type"`
	EntityID           string `json:"entity_id"`
	DisplayName        string `json:"display_name"`
	AvatarURL          string `json:"avatar_url"`
	IsVerified         bool   `json:"is_verified"`
	MediaType          string `json:"media_type"`
	MediaURL           string `json:"media_url"`
	CreatedBy          string `json:"created_by"`
	NameCreatedBy      string `json:"name_created_by"`
	TargetCountryCode  string `json:"target_country_code"`
	PublisherCountryCode string `json:"publisher_country_code"`
}

type FeedByCountryDTO struct {
	CountryCode string `json:"country_code"`
}
