package stories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/eduardeddy/content-ms/internal/publisher"
	"github.com/google/uuid"
)

type Service interface {
	CreateStory(ctx context.Context, dto CreateStoryDTO) error
	DeleteStory(ctx context.Context, dto DeleteStoryDTO) error
	ListByEntity(ctx context.Context, dto ListByEntityDTO) (*PublisherFeedResponse, error)
	GetFeedByPublisher(ctx context.Context, publisherId uuid.UUID) (*PublisherFeedResponse, error)
	GetFeedByCountry(ctx context.Context, countryCode string) ([]PublisherFeedResponse, error)
}

type service struct {
	repo       Repository
	publishers publisher.Repository
}

func NewService(repo Repository, publishers publisher.Repository) Service {
	return &service{repo: repo, publishers: publishers}
}

func (s *service) CreateStory(ctx context.Context, dto CreateStoryDTO) error {
	if dto.EntityType == "" || dto.EntityID == "" {
		return fmt.Errorf("entity_type y entity_id son obligatorios")
	}
	if dto.MediaURL == "" || dto.MediaType == "" {
		return fmt.Errorf("media_url y media_type son obligatorios")
	}
	if dto.CreatedBy == "" {
		return fmt.Errorf("created_by es obligatorio")
	}

	publisherCountry := dto.PublisherCountryCode
	if publisherCountry == "" && dto.EntityType == "platform" {
		publisherCountry = "GLOBAL"
	}

	err := s.publishers.Upsert(ctx, publisher.UpdatePublisherDTO{
		EntityType:  dto.EntityType,
		EntityID:    dto.EntityID,
		DisplayName: dto.DisplayName,
		AvatarURL:   dto.AvatarURL,
		IsVerified:  dto.IsVerified,
		CountryCode: publisherCountry,
	})
	if err != nil {
		return fmt.Errorf("no se pudo sincronizar publisher: %w", err)
	}

	publisherID, err := s.publishers.GetIDByEntity(ctx, dto.EntityType, dto.EntityID)
	if err != nil {
		return fmt.Errorf("publisher no encontrado para %s/%s: %w", dto.EntityType, dto.EntityID, err)
	}

	return s.repo.Create(ctx, publisherID, dto)
}

func (s *service) GetFeedByCountry(ctx context.Context, countryCode string) ([]PublisherFeedResponse, error) {
	if countryCode == "" {
		return nil, fmt.Errorf("country_code es obligatorio")
	}
	return s.repo.GetFeedByCountry(ctx, countryCode)
}

func (s *service) DeleteStory(ctx context.Context, dto DeleteStoryDTO) error {
	if dto.ID == "" || dto.EntityType == "" || dto.EntityID == "" {
		return fmt.Errorf("id, entity_type y entity_id son obligatorios")
	}
	err := s.repo.DeleteByEntity(ctx, dto)
	if err == sql.ErrNoRows {
		return fmt.Errorf("historia no encontrada o no pertenece a tu publisher")
	}
	return err
}

func (s *service) ListByEntity(ctx context.Context, dto ListByEntityDTO) (*PublisherFeedResponse, error) {
	if dto.EntityType == "" || dto.EntityID == "" {
		return nil, fmt.Errorf("entity_type y entity_id son obligatorios")
	}
	return s.repo.ListByEntity(ctx, dto.EntityType, dto.EntityID)
}

func (s *service) GetFeedByPublisher(ctx context.Context, publisherId uuid.UUID) (*PublisherFeedResponse, error) {
	return s.repo.GetByPublisher(ctx, publisherId)
}
