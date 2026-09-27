package articles

import "time"

type ArticleDTO struct {
	ID                 string    `json:"id"`
	PublisherID        string    `json:"publisher_id"`
	CategoryID         string    `json:"category_id"`
	CategoryName       string    `json:"category_name"`
	Title              string    `json:"title"`
	Slug               string    `json:"slug"`
	Summary            string    `json:"summary"`
	Content            string    `json:"content"`
	CoverImageURL      string    `json:"cover_image_url"`
	CloudinaryPublicID string    `json:"cloudinary_public_id"`
	Status             string    `json:"status"`
	IsFeatured         bool      `json:"is_featured"`
	ViewsCount         int       `json:"views_count"`
	TargetCountryCode  string    `json:"target_country_code"`
	PublisherName      string    `json:"publisher_name"`
	EntityType         string    `json:"entity_type"`
	EntityID           string    `json:"entity_id"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	PublishedAt        time.Time `json:"published_at"`
	CreatedBy          string    `json:"created_by"`
	NameCreatedBy      string    `json:"name_created_by"`
}

type CreateArticleDTO struct {
	EntityType           string `json:"entity_type"`
	EntityID             string `json:"entity_id"`
	DisplayName          string `json:"display_name"`
	AvatarURL            string `json:"avatar_url"`
	IsVerified           bool   `json:"is_verified"`
	PublisherCountryCode string `json:"publisher_country_code"`
	TargetCountryCode    string `json:"target_country_code"`
	CategoryID           string `json:"category_id"`
	Title                string `json:"title"`
	Slug                 string `json:"slug"`
	Summary              string `json:"summary"`
	Content              string `json:"content"`
	CoverImageURL        string `json:"cover_image_url"`
	CloudinaryPublicID   string `json:"cloudinary_public_id"`
	Status               string `json:"status"`
	IsFeatured           bool   `json:"is_featured"`
	CreatedBy            string `json:"created_by"`
	NameCreatedBy        string `json:"name_created_by"`
}

type ListByEntityDTO struct {
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

type FeedByCountryDTO struct {
	CountryCode string `json:"country_code"`
}

type DeleteArticleDTO struct {
	ID         string `json:"id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
}

type ArticleIDDTO struct {
	ID string `json:"id"`
}
