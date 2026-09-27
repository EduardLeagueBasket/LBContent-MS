package article_categories

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Repository interface {
	CreateArticleCategory(ctx context.Context, dto ArticleCategoryDTO) (ArticleCategoryDTO, error)
	ListArticleCategories(ctx context.Context) ([]ArticleCategoryDTO, error)
	DeleteArticleCategory(ctx context.Context, id string) error
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateArticleCategory(ctx context.Context, dto ArticleCategoryDTO) (ArticleCategoryDTO, error) {
	var created ArticleCategoryDTO
	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO article_categories (name, slug, created_at)
		 VALUES ($1, $2, $3)
		 RETURNING id, name, slug, created_at`,
		dto.Name,
		dto.Slug,
		time.Now(),
	).Scan(&created.ID, &created.Name, &created.Slug, &created.CreatedAt)
	if err != nil {
		return ArticleCategoryDTO{}, fmt.Errorf("failed to create article category: %w", err)
	}
	return created, nil
}

func (r *repository) ListArticleCategories(ctx context.Context) ([]ArticleCategoryDTO, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, name, slug, created_at FROM article_categories ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list article categories: %w", err)
	}
	defer rows.Close()

	categories := make([]ArticleCategoryDTO, 0)
	for rows.Next() {
		var category ArticleCategoryDTO
		if err := rows.Scan(&category.ID, &category.Name, &category.Slug, &category.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan article category: %w", err)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate article categories: %w", err)
	}

	return categories, nil
}

func (r *repository) DeleteArticleCategory(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM article_categories WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete article category: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to delete article category: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
