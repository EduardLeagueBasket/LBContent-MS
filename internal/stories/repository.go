package stories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, publisherID string, dto CreateStoryDTO) error
	DeleteByEntity(ctx context.Context, dto DeleteStoryDTO) error
	ListByEntity(ctx context.Context, entityType, entityID string) (*PublisherFeedResponse, error)
	GetByPublisher(ctx context.Context, publisherID uuid.UUID) (*PublisherFeedResponse, error)
	GetFeedByCountry(ctx context.Context, countryCode string) ([]PublisherFeedResponse, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, publisherID string, dto CreateStoryDTO) error {
	targetCountry := strings.TrimSpace(dto.TargetCountryCode)
	if targetCountry == "" {
		targetCountry = "GLOBAL"
	}

	query := `
		INSERT INTO stories (
			id, publisher_id, media_url, media_type, cloudinary_public_id,
			duration_seconds, expires_at, created_by, name_created_by, target_country_code
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		uuid.New().String(),
		publisherID,
		dto.MediaURL,
		dto.MediaType,
		nil,
		nil,
		time.Now().Add(24*time.Hour),
		dto.CreatedBy,
		dto.NameCreatedBy,
		targetCountry,
	)
	return err
}

func (r *repository) DeleteByEntity(ctx context.Context, dto DeleteStoryDTO) error {
	query := `
		DELETE FROM stories s
		USING publishers p
		WHERE s.id = $1
		  AND s.publisher_id = p.id
		  AND p.entity_type = $2
		  AND p.entity_id = $3
	`
	result, err := r.db.ExecContext(ctx, query, dto.ID, dto.EntityType, dto.EntityID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) ListByEntity(ctx context.Context, entityType, entityID string) (*PublisherFeedResponse, error) {
	query := `
		SELECT 
			p.id AS publisher_id,
			p.entity_type,
			p.entity_id,
			p.display_name,
			COALESCE(p.avatar_url, '') AS avatar_url,
			p.is_verified,
			COALESCE(p.country_code, '') AS country_code,
			s.id AS story_id,
			s.media_url,
			s.media_type,
			COALESCE(s.duration_seconds, 5) AS duration_seconds,
			s.expires_at,
			s.created_at,
			s.created_by::text,
			s.name_created_by
		FROM publishers p
		LEFT JOIN stories s ON s.publisher_id = p.id AND s.expires_at > NOW()
		WHERE p.entity_type = $1 AND p.entity_id = $2
		ORDER BY s.created_at DESC NULLS LAST;
	`

	rows, err := r.db.QueryContext(ctx, query, entityType, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result *PublisherFeedResponse

	for rows.Next() {
		var pubID, entityTypeVal, entityIDVal, displayName, avatarURL, countryCode string
		var isVerified bool
		var storyID, mediaURL, mediaType, createdBy, nameCreatedBy sql.NullString
		var durationSeconds sql.NullInt64
		var expiresAt, createdAt sql.NullTime

		err := rows.Scan(
			&pubID,
			&entityTypeVal,
			&entityIDVal,
			&displayName,
			&avatarURL,
			&isVerified,
			&countryCode,
			&storyID,
			&mediaURL,
			&mediaType,
			&durationSeconds,
			&expiresAt,
			&createdAt,
			&createdBy,
			&nameCreatedBy,
		)
		if err != nil {
			return nil, err
		}

		if result == nil {
			result = &PublisherFeedResponse{
				PublisherID: pubID,
				EntityType:  entityTypeVal,
				EntityID:    entityIDVal,
				DisplayName: displayName,
				AvatarURL:   avatarURL,
				IsVerified:  isVerified,
				CountryCode: countryCode,
				Stories:     []StoryResponse{},
			}
		}

		if !storyID.Valid {
			continue
		}

		story := StoryResponse{
			ID:              storyID.String,
			MediaURL:        mediaURL.String,
			MediaType:       mediaType.String,
			DurationSeconds: int(durationSeconds.Int64),
			ExpiresAt:       expiresAt.Time,
			CreatedAt:       createdAt.Time,
			CreatedBy:       createdBy.String,
			NameCreatedBy:   nameCreatedBy.String,
		}
		result.Stories = append(result.Stories, story)
	}

	if result == nil {
		return &PublisherFeedResponse{
			EntityType: entityType,
			EntityID:   entityID,
			Stories:    []StoryResponse{},
		}, nil
	}

	return result, nil
}

func (r *repository) GetByPublisher(ctx context.Context, publisherId uuid.UUID) (*PublisherFeedResponse, error) {
	query := `
		SELECT 
			p.id AS publisher_id,
			p.entity_type,
			p.entity_id,
			p.display_name,
			COALESCE(p.avatar_url, '') AS avatar_url,
			p.is_verified,
			p.country_code,
			s.id AS story_id,
			s.media_url,
			s.media_type,
			COALESCE(s.duration_seconds, 5) AS duration_seconds,
			s.expires_at,
			s.created_at
		FROM publishers p
		INNER JOIN stories s ON s.publisher_id = p.id
		WHERE p.id = $1 AND s.expires_at > NOW()
		ORDER BY s.created_at DESC;
	`

	rows, err := r.db.QueryContext(ctx, query, publisherId.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result *PublisherFeedResponse

	for rows.Next() {
		var pubID, entityType, entityID, displayName, avatarURL, countryCode string
		var isVerified bool
		var story StoryResponse

		err := rows.Scan(
			&pubID,
			&entityType,
			&entityID,
			&displayName,
			&avatarURL,
			&isVerified,
			&countryCode,
			&story.ID,
			&story.MediaURL,
			&story.MediaType,
			&story.DurationSeconds,
			&story.ExpiresAt,
			&story.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if result == nil {
			result = &PublisherFeedResponse{
				PublisherID: pubID,
				EntityType:  entityType,
				EntityID:    entityID,
				DisplayName: displayName,
				AvatarURL:   avatarURL,
				IsVerified:  isVerified,
				CountryCode: countryCode,
				Stories:     []StoryResponse{},
			}
		}

		result.Stories = append(result.Stories, story)
	}

	if result == nil {
		return nil, sql.ErrNoRows
	}

	return result, nil
}

func (r *repository) GetFeedByCountry(ctx context.Context, countryCode string) ([]PublisherFeedResponse, error) {
	code := strings.ToUpper(strings.TrimSpace(countryCode))
	if code == "" {
		return []PublisherFeedResponse{}, nil
	}

	// Feed del país:
	// - publishers locales (team / national team / league) con country_code = país
	// - stories de platform (League Basket) con target_country_code = país o GLOBAL
	query := `
		SELECT 
			p.id AS publisher_id,
			p.entity_type,
			p.entity_id,
			p.display_name,
			COALESCE(p.avatar_url, '') AS avatar_url,
			p.is_verified,
			COALESCE(p.country_code, 'GLOBAL') AS country_code,
			s.id AS story_id,
			s.media_url,
			s.media_type,
			COALESCE(s.duration_seconds, 5) AS duration_seconds,
			s.expires_at,
			s.created_at,
			COALESCE(s.target_country_code, 'GLOBAL') AS target_country_code
		FROM publishers p
		INNER JOIN stories s ON s.publisher_id = p.id
		WHERE s.expires_at > NOW()
		  AND (
		    (p.entity_type <> 'platform' AND UPPER(p.country_code) = $1)
		    OR
		    (
		      p.entity_type = 'platform'
		      AND (
		        UPPER(s.target_country_code) = $1
		        OR UPPER(s.target_country_code) = 'GLOBAL'
		      )
		    )
		  )
		ORDER BY s.created_at DESC;
	`
	print("rows #########################################################")
	rows, err := r.db.QueryContext(ctx, query, code)
	fmt.Printf("query:\n%s\n", query)
	fmt.Printf("err:\n%s\n", err)
	fmt.Printf("RIW:\n%s\n", rows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	publisherMap := make(map[string]*PublisherFeedResponse)
	var orderedPublisherIDs []string

	for rows.Next() {
		var pubID, entityType, entityID, displayName, avatarURL, country string
		var isVerified bool
		var story StoryResponse

		err := rows.Scan(
			&pubID,
			&entityType,
			&entityID,
			&displayName,
			&avatarURL,
			&isVerified,
			&country,
			&story.ID,
			&story.MediaURL,
			&story.MediaType,
			&story.DurationSeconds,
			&story.ExpiresAt,
			&story.CreatedAt,
			&story.TargetCountryCode,
		)
		if err != nil {
			return nil, err
		}

		// Solo publishers con story real (INNER JOIN, pero por si acaso).
		if story.ID == "" {
			continue
		}

		if _, exists := publisherMap[pubID]; !exists {
			publisherMap[pubID] = &PublisherFeedResponse{
				PublisherID: pubID,
				EntityType:  entityType,
				EntityID:    entityID,
				DisplayName: displayName,
				AvatarURL:   avatarURL,
				IsVerified:  isVerified,
				CountryCode: country,
				Stories:     []StoryResponse{},
			}
			orderedPublisherIDs = append(orderedPublisherIDs, pubID)
		}

		publisherMap[pubID].Stories = append(publisherMap[pubID].Stories, story)
	}

	var result []PublisherFeedResponse
	for _, id := range orderedPublisherIDs {
		pub := *publisherMap[id]
		if len(pub.Stories) == 0 {
			continue
		}
		result = append(result, pub)
	}
	if result == nil {
		result = []PublisherFeedResponse{}
	}

	return result, nil
}
