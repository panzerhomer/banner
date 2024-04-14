package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/panzerhomer/banner/internal/domain"
)

var ctx = context.Background()

type bannerRepo struct {
	db *pgx.Conn
}

func NewBannerRepo(db *pgx.Conn) *bannerRepo {
	return &bannerRepo{db}
}

func (r *bannerRepo) InsertBanner(banner domain.Banner) (int64, error) {
	const op = "repository.postgres.InsertBanner"

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	const insertBannerQuery = "INSERT INTO banners(feature, tags, is_active) VALUES ($1, $2, $3) RETURNING banner_id"

	var bannerID int64
	if err := r.db.QueryRow(ctx, insertBannerQuery, banner.FeatureID, banner.TagIds, banner.IsActive).Scan(&bannerID); err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	const insertBannerVersionQuery = "INSERT INTO banner_version(banner_id, banner_info) VALUES ($1, $2)"
	if _, err := r.db.Exec(ctx, insertBannerVersionQuery, bannerID, banner.Content); err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return bannerID, nil
}

func (r *bannerRepo) GetBanners(tagIDs []int, featureID int, limit int, offset int) ([]domain.Banner, error) {
	const op = "repository.postgres.GetBanners"

	const selectBannersByFeatureAndTagsQuery = `SELECT 
		banner_id, 
		feature, 
		tags, 
		is_active, 
		banner_info, 
		created_at, 
		updated_at 
	FROM 
		banner 
	JOIN banner_version 
	ON 
		banner.banner_id = banner_version.banner_id 
	WHERE 
		banner.feature = $1 AND tags @> $2 
	LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(ctx, selectBannersByFeatureAndTagsQuery, featureID, tagIDs, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var banners []domain.Banner
	for rows.Next() {
		var b domain.Banner
		err := rows.Scan(&b.BannerID, &b.FeatureID, &b.TagIds, &b.IsActive, &b.Content, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		banners = append(banners, b)
	}

	return banners, nil
}

// func (r *bannerRepo) GetBannerByID(bannerID int64)
