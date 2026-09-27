package articles

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, publisherID string, dto CreateArticleDTO) error
	ListByEntity(ctx context.Context, entityType, entityID string) ([]ArticleDTO, error)
	ListByCountry(ctx context.Context, countryCode string) ([]ArticleDTO, error)
	ListArticles(ctx context.Context) ([]ArticleDTO, error)
	GetArticle(ctx context.Context, id string) (ArticleDTO, error)
	UpdateArticle(ctx context.Context, id string, dto ArticleDTO) error
	DeleteByEntity(ctx context.Context, dto DeleteArticleDTO) error
	DeleteArticle(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

const articleSelect = `
	a.id, a.publisher_id, a.category_id, a.title, a.slug, a.summary, a.content_body,
	a.cover_image_url, a.cloudinary_public_id, a.status, a.is_featured, a.views_count,
	a.created_at, a.updated_at, a.published_at, a.created_by::text, a.name_created_by,
	COALESCE(a.target_country_code, 'GLOBAL'),
	COALESCE(c.name, ''),
	COALESCE(p.display_name, ''),
	COALESCE(p.entity_type, ''),
	COALESCE(p.entity_id::text, '')
`

const articleFrom = `
	FROM articles a
	INNER JOIN publishers p ON p.id = a.publisher_id
	LEFT JOIN article_categories c ON c.id = a.category_id
`

func (r *repository) Create(ctx context.Context, publisherID string, dto CreateArticleDTO) error {
	now := time.Now()
	id := uuid.New().String()
	slug := dto.Slug
	if slug == "" {
		slug = slugify(dto.Title)
	}
	slug = fmt.Sprintf("%s-%s", slug, id[:8])

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO articles (
			id, publisher_id, category_id, title, slug, summary, content_body,
			cover_image_url, cloudinary_public_id, status, is_featured,
			target_country_code, created_at, updated_at, published_at, created_by, name_created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)`,
		id,
		publisherID,
		nullableString(dto.CategoryID),
		dto.Title,
		slug,
		dto.Summary,
		dto.Content,
		nullableString(dto.CoverImageURL),
		nullableString(dto.CloudinaryPublicID),
		dto.Status,
		dto.IsFeatured,
		dto.TargetCountryCode,
		now,
		now,
		now,
		dto.CreatedBy,
		dto.NameCreatedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}

func (r *repository) ListByEntity(ctx context.Context, entityType, entityID string) ([]ArticleDTO, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT `+articleSelect+articleFrom+`
		 WHERE p.entity_type = $1 AND p.entity_id = $2
		 ORDER BY a.published_at DESC, a.created_at DESC`,
		entityType,
		entityID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list articles by entity: %w", err)
	}
	return scanArticles(rows)
}

func (r *repository) ListByCountry(ctx context.Context, countryCode string) ([]ArticleDTO, error) {
	code := strings.ToUpper(strings.TrimSpace(countryCode))
	if code == "" {
		return []ArticleDTO{}, nil
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT `+articleSelect+articleFrom+`
		 WHERE a.status = 'published'
		   AND (
		     (p.entity_type <> 'platform' AND UPPER(COALESCE(p.country_code, '')) = $1)
		     OR
		     (
		       p.entity_type = 'platform'
		       AND (
		         UPPER(COALESCE(a.target_country_code, 'GLOBAL')) = $1
		         OR UPPER(COALESCE(a.target_country_code, 'GLOBAL')) = 'GLOBAL'
		       )
		     )
		   )
		 ORDER BY a.is_featured DESC, a.published_at DESC, a.created_at DESC`,
		code,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list articles by country: %w", err)
	}
	return scanArticles(rows)
}

func (r *repository) ListArticles(ctx context.Context) ([]ArticleDTO, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT `+articleSelect+articleFrom+`
		 ORDER BY a.published_at DESC, a.created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list articles: %w", err)
	}
	return scanArticles(rows)
}

func (r *repository) GetArticle(ctx context.Context, id string) (ArticleDTO, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT `+articleSelect+articleFrom+` WHERE a.id = $1`,
		id,
	)
	article, err := scanArticle(row)
	if err == sql.ErrNoRows {
		return ArticleDTO{}, err
	}
	if err != nil {
		return ArticleDTO{}, fmt.Errorf("failed to get article: %w", err)
	}
	return article, nil
}

func (r *repository) UpdateArticle(ctx context.Context, id string, dto ArticleDTO) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE articles SET
			category_id = $1,
			title = $2,
			slug = $3,
			summary = $4,
			content_body = $5,
			cover_image_url = $6,
			cloudinary_public_id = $7,
			status = $8,
			is_featured = $9,
			target_country_code = $10,
			updated_at = $11
		WHERE id = $12`,
		nullableString(dto.CategoryID),
		dto.Title,
		dto.Slug,
		dto.Summary,
		dto.Content,
		nullableString(dto.CoverImageURL),
		nullableString(dto.CloudinaryPublicID),
		dto.Status,
		dto.IsFeatured,
		dto.TargetCountryCode,
		time.Now(),
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) DeleteByEntity(ctx context.Context, dto DeleteArticleDTO) error {
	result, err := r.db.ExecContext(
		ctx,
		`DELETE FROM articles a
		 USING publishers p
		 WHERE a.id = $1
		   AND a.publisher_id = p.id
		   AND p.entity_type = $2
		   AND p.entity_id = $3`,
		dto.ID,
		dto.EntityType,
		dto.EntityID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *repository) DeleteArticle(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM articles WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func scanArticles(rows *sql.Rows) ([]ArticleDTO, error) {
	defer rows.Close()
	articles := make([]ArticleDTO, 0)
	for rows.Next() {
		article, err := scanArticle(rows)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate articles: %w", err)
	}
	return articles, nil
}

func scanArticle(scanner interface{ Scan(dest ...any) error }) (ArticleDTO, error) {
	var dto ArticleDTO
	var categoryID, coverURL, cloudinaryID sql.NullString
	var publishedAt sql.NullTime

	err := scanner.Scan(
		&dto.ID,
		&dto.PublisherID,
		&categoryID,
		&dto.Title,
		&dto.Slug,
		&dto.Summary,
		&dto.Content,
		&coverURL,
		&cloudinaryID,
		&dto.Status,
		&dto.IsFeatured,
		&dto.ViewsCount,
		&dto.CreatedAt,
		&dto.UpdatedAt,
		&publishedAt,
		&dto.CreatedBy,
		&dto.NameCreatedBy,
		&dto.TargetCountryCode,
		&dto.CategoryName,
		&dto.PublisherName,
		&dto.EntityType,
		&dto.EntityID,
	)
	if err != nil {
		return ArticleDTO{}, err
	}

	dto.CategoryID = categoryID.String
	dto.CoverImageURL = coverURL.String
	dto.CloudinaryPublicID = cloudinaryID.String
	dto.PublishedAt = publishedAt.Time
	return dto, nil
}

func nullableString(value string) any {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return value
}
