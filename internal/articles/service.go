package articles

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/eduardeddy/content-ms/internal/publisher"
)

type Service interface {
	CreateArticle(ctx context.Context, dto CreateArticleDTO) error
	ListByEntity(ctx context.Context, dto ListByEntityDTO) ([]ArticleDTO, error)
	ListByCountry(ctx context.Context, countryCode string) ([]ArticleDTO, error)
	ListArticles(ctx context.Context) ([]ArticleDTO, error)
	GetArticle(ctx context.Context, id string) (ArticleDTO, error)
	UpdateArticle(ctx context.Context, id string, dto ArticleDTO) error
	DeleteArticle(ctx context.Context, dto DeleteArticleDTO) error
	DeleteArticleByID(ctx context.Context, id string) error
}

type service struct {
	repo       Repository
	publishers publisher.Repository
}

func NewService(repo Repository, publishers publisher.Repository) Service {
	return &service{repo: repo, publishers: publishers}
}

func (s *service) CreateArticle(ctx context.Context, dto CreateArticleDTO) error {
	dto.Title = strings.TrimSpace(dto.Title)
	dto.Content = strings.TrimSpace(dto.Content)
	dto.CategoryID = strings.TrimSpace(dto.CategoryID)
	dto.EntityType = strings.TrimSpace(dto.EntityType)
	dto.EntityID = strings.TrimSpace(dto.EntityID)
	dto.CreatedBy = strings.TrimSpace(dto.CreatedBy)
	dto.NameCreatedBy = strings.TrimSpace(dto.NameCreatedBy)
	dto.Summary = strings.TrimSpace(dto.Summary)
	dto.Slug = strings.TrimSpace(dto.Slug)
	dto.Status = strings.TrimSpace(dto.Status)
	dto.CoverImageURL = strings.TrimSpace(dto.CoverImageURL)
	dto.CloudinaryPublicID = strings.TrimSpace(dto.CloudinaryPublicID)

	if dto.Title == "" {
		return fmt.Errorf("title es obligatorio")
	}
	if dto.Content == "" {
		return fmt.Errorf("content es obligatorio")
	}
	if dto.CategoryID == "" {
		return fmt.Errorf("category_id es obligatorio")
	}
	if dto.EntityType == "" || dto.EntityID == "" {
		return fmt.Errorf("entity_type y entity_id son obligatorios")
	}
	if dto.CreatedBy == "" {
		return fmt.Errorf("created_by es obligatorio")
	}
	if dto.Summary == "" {
		dto.Summary = summarize(dto.Content)
	}
	if dto.Slug == "" {
		dto.Slug = slugify(dto.Title)
	}
	if dto.Status == "" {
		dto.Status = "published"
	}
	if dto.NameCreatedBy == "" {
		dto.NameCreatedBy = "Usuario"
	}

	publisherCountry := strings.ToUpper(strings.TrimSpace(dto.PublisherCountryCode))
	if publisherCountry == "" && dto.EntityType == "platform" {
		publisherCountry = "GLOBAL"
	}

	targetCountry := strings.ToUpper(strings.TrimSpace(dto.TargetCountryCode))
	if targetCountry == "" {
		if dto.EntityType == "platform" {
			targetCountry = "GLOBAL"
		} else if publisherCountry != "" {
			targetCountry = publisherCountry
		} else {
			targetCountry = "GLOBAL"
		}
	}
	dto.TargetCountryCode = targetCountry

	if err := s.publishers.Upsert(ctx, publisher.UpdatePublisherDTO{
		EntityType:  dto.EntityType,
		EntityID:    dto.EntityID,
		DisplayName: dto.DisplayName,
		AvatarURL:   dto.AvatarURL,
		IsVerified:  dto.IsVerified,
		CountryCode: publisherCountry,
	}); err != nil {
		return fmt.Errorf("no se pudo sincronizar publisher: %w", err)
	}

	publisherID, err := s.publishers.GetIDByEntity(ctx, dto.EntityType, dto.EntityID)
	if err != nil {
		return fmt.Errorf("publisher no encontrado para %s/%s: %w", dto.EntityType, dto.EntityID, err)
	}

	return s.repo.Create(ctx, publisherID, dto)
}

func (s *service) ListByEntity(ctx context.Context, dto ListByEntityDTO) ([]ArticleDTO, error) {
	if strings.TrimSpace(dto.EntityType) == "" || strings.TrimSpace(dto.EntityID) == "" {
		return nil, fmt.Errorf("entity_type y entity_id son obligatorios")
	}
	return s.repo.ListByEntity(ctx, dto.EntityType, dto.EntityID)
}

func (s *service) ListByCountry(ctx context.Context, countryCode string) ([]ArticleDTO, error) {
	if strings.TrimSpace(countryCode) == "" {
		return nil, fmt.Errorf("country_code es obligatorio")
	}
	return s.repo.ListByCountry(ctx, countryCode)
}

func (s *service) ListArticles(ctx context.Context) ([]ArticleDTO, error) {
	return s.repo.ListArticles(ctx)
}

func (s *service) GetArticle(ctx context.Context, id string) (ArticleDTO, error) {
	if strings.TrimSpace(id) == "" {
		return ArticleDTO{}, fmt.Errorf("id es obligatorio")
	}
	return s.repo.GetArticle(ctx, id)
}

func (s *service) UpdateArticle(ctx context.Context, id string, dto ArticleDTO) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id es obligatorio")
	}
	dto.Title = strings.TrimSpace(dto.Title)
	dto.Content = strings.TrimSpace(dto.Content)
	if dto.Title == "" || dto.Content == "" {
		return fmt.Errorf("title y content son obligatorios")
	}
	if strings.TrimSpace(dto.TargetCountryCode) == "" {
		dto.TargetCountryCode = "GLOBAL"
	}
	return s.repo.UpdateArticle(ctx, id, dto)
}

func (s *service) DeleteArticle(ctx context.Context, dto DeleteArticleDTO) error {
	if dto.ID == "" || dto.EntityType == "" || dto.EntityID == "" {
		return fmt.Errorf("id, entity_type y entity_id son obligatorios")
	}
	err := s.repo.DeleteByEntity(ctx, dto)
	if err == sql.ErrNoRows {
		return fmt.Errorf("artículo no encontrado o no pertenece a tu publisher")
	}
	return err
}

func (s *service) DeleteArticleByID(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id es obligatorio")
	}
	return s.repo.DeleteArticle(ctx, id)
}

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	var builder strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(value)) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' {
			builder.WriteRune(r)
		}
	}

	slug := nonSlugChars.ReplaceAllString(builder.String(), "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "articulo"
	}
	return slug
}

func summarize(content string) string {
	runes := []rune(content)
	if len(runes) <= 200 {
		return content
	}
	return strings.TrimSpace(string(runes[:200])) + "..."
}
