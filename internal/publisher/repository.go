package publisher

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type Repository interface {
	Upsert(ctx context.Context, dto UpdatePublisherDTO) error
	GetIDByEntity(ctx context.Context, entityType, entityID string) (string, error)
	GetByEntity(ctx context.Context, entityType, entityID string) (*Publisher, error)
	List(ctx context.Context, entityType string) ([]Publisher, error)
	CreatePlatform(ctx context.Context, dto CreatePlatformPublisherDTO) (*Publisher, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Upsert(ctx context.Context, dto UpdatePublisherDTO) error {
	countryCode := strings.TrimSpace(dto.CountryCode)
	if countryCode == "" {
		countryCode = "GLOBAL"
	}

	query := `
		INSERT INTO publishers (id, entity_type, entity_id, display_name, avatar_url, is_verified, country_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (entity_type, entity_id) 
		DO UPDATE SET 
			display_name = EXCLUDED.display_name,
			avatar_url = EXCLUDED.avatar_url,
			is_verified = EXCLUDED.is_verified,
			country_code = EXCLUDED.country_code;
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		uuid.New().String(),
		dto.EntityType,
		dto.EntityID,
		dto.DisplayName,
		dto.AvatarURL,
		dto.IsVerified,
		countryCode,
	)
	return err
}

func (r *repository) GetIDByEntity(ctx context.Context, entityType, entityID string) (string, error) {
	var id string
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id FROM publishers WHERE entity_type = $1 AND entity_id = $2`,
		entityType,
		entityID,
	).Scan(&id)
	return id, err
}

func (r *repository) GetByEntity(ctx context.Context, entityType, entityID string) (*Publisher, error) {
	query := `
		SELECT id, entity_type, entity_id, display_name, COALESCE(avatar_url, ''), is_verified,
		       COALESCE(country_code, 'GLOBAL'), created_at
		FROM publishers
		WHERE entity_type = $1 AND entity_id = $2
	`
	var p Publisher
	err := r.db.QueryRowContext(ctx, query, entityType, entityID).Scan(
		&p.ID,
		&p.EntityType,
		&p.EntityID,
		&p.DisplayName,
		&p.AvatarURL,
		&p.IsVerified,
		&p.CountryCode,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) List(ctx context.Context, entityType string) ([]Publisher, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if strings.TrimSpace(entityType) != "" {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, entity_type, entity_id, display_name, COALESCE(avatar_url, ''), is_verified,
			       COALESCE(country_code, 'GLOBAL'), created_at
			FROM publishers
			WHERE entity_type = $1
			ORDER BY display_name ASC
		`, entityType)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, entity_type, entity_id, display_name, COALESCE(avatar_url, ''), is_verified,
			       COALESCE(country_code, 'GLOBAL'), created_at
			FROM publishers
			ORDER BY entity_type ASC, display_name ASC
		`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Publisher
	for rows.Next() {
		var p Publisher
		if err := rows.Scan(
			&p.ID,
			&p.EntityType,
			&p.EntityID,
			&p.DisplayName,
			&p.AvatarURL,
			&p.IsVerified,
			&p.CountryCode,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	if result == nil {
		result = []Publisher{}
	}
	return result, rows.Err()
}

func (r *repository) CreatePlatform(ctx context.Context, dto CreatePlatformPublisherDTO) (*Publisher, error) {
	displayName := strings.TrimSpace(dto.DisplayName)
	if displayName == "" {
		return nil, fmt.Errorf("display_name es obligatorio")
	}

	entityID := uuid.New().String()
	id := uuid.New().String()

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO publishers (id, entity_type, entity_id, display_name, avatar_url, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, id, EntityTypePlatform, entityID, displayName, dto.AvatarURL, dto.IsVerified)
	if err != nil {
		return nil, err
	}

	return r.GetByEntity(ctx, EntityTypePlatform, entityID)
}
