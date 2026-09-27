package publisher

import (
	"context"
	"fmt"
	"strings"
)

type Service interface {
	UpsertPublisher(ctx context.Context, dto UpdatePublisherDTO) error
	ListPublishers(ctx context.Context, dto ListPublishersDTO) ([]Publisher, error)
	CreatePlatformPublisher(ctx context.Context, dto CreatePlatformPublisherDTO) (*Publisher, error)
	GetByEntity(ctx context.Context, entityType, entityID string) (*Publisher, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) UpsertPublisher(ctx context.Context, dto UpdatePublisherDTO) error {
	if strings.TrimSpace(dto.EntityType) == "" || strings.TrimSpace(dto.EntityID) == "" {
		return fmt.Errorf("entity_type y entity_id son obligatorios")
	}
	if strings.TrimSpace(dto.DisplayName) == "" {
		return fmt.Errorf("display_name es obligatorio")
	}
	return s.repo.Upsert(ctx, dto)
}

func (s *service) ListPublishers(ctx context.Context, dto ListPublishersDTO) ([]Publisher, error) {
	return s.repo.List(ctx, dto.EntityType)
}

func (s *service) CreatePlatformPublisher(ctx context.Context, dto CreatePlatformPublisherDTO) (*Publisher, error) {
	return s.repo.CreatePlatform(ctx, dto)
}

func (s *service) GetByEntity(ctx context.Context, entityType, entityID string) (*Publisher, error) {
	return s.repo.GetByEntity(ctx, entityType, entityID)
}
