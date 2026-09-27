package article_categories

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

type Service interface {
	CreateArticleCategory(ctx context.Context, dto ArticleCategoryDTO) (ArticleCategoryDTO, error)
	ListArticleCategories(ctx context.Context) ([]ArticleCategoryDTO, error)
	DeleteArticleCategory(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateArticleCategory(ctx context.Context, dto ArticleCategoryDTO) (ArticleCategoryDTO, error) {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Slug = strings.TrimSpace(dto.Slug)
	if dto.Name == "" {
		return ArticleCategoryDTO{}, fmt.Errorf("name es obligatorio")
	}
	if dto.Slug == "" {
		dto.Slug = slugify(dto.Name)
	}
	return s.repo.CreateArticleCategory(ctx, dto)
}

func (s *service) ListArticleCategories(ctx context.Context) ([]ArticleCategoryDTO, error) {
	return s.repo.ListArticleCategories(ctx)
}

func (s *service) DeleteArticleCategory(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id es obligatorio")
	}
	return s.repo.DeleteArticleCategory(ctx, id)
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
		return "categoria"
	}
	return slug
}
