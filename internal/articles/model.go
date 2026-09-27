package articles

import "time"

type Article struct {
	ID                 string    `json:"id" db:"id"`
	PublisherID        string    `json:"publisher_id" db:"publisher_id"`
	CategoryID         string    `json:"category_id" db:"category_id"`
	Title              string    `json:"title" db:"title"`
	Slug               string    `json:"slug" db:"slug"`
	Summary            string    `json:"summary" db:"summary"`
	Content            string    `json:"content" db:"content_body"`
	CoverImageURL      string    `json:"cover_image_url" db:"cover_image_url"`
	CloudinaryPublicID string    `json:"cloudinary_public_id" db:"cloudinary_public_id"`
	Status             string    `json:"status" db:"status"`
	IsFeatured         bool      `json:"is_featured" db:"is_featured"`
	ViewsCount         int       `json:"views_count" db:"views_count"`
	TargetCountryCode  string    `json:"target_country_code" db:"target_country_code"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	PublishedAt        time.Time `json:"published_at" db:"published_at"`
	CreatedBy          string    `json:"created_by" db:"created_by"`
	NameCreatedBy      string    `json:"name_created_by" db:"name_created_by"`
}
