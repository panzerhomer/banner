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

func (r *bannerRepo) GetBanners(tagIDs []int64, featureID int64, limit int64, offset int64) ([]domain.Banner, error) {
	const op = "repository.postgres.GetBanners"

	const selectBannersByFeatureAndTagsQuery = `
	SELECT 
		b.banner_id, 
		b.feature, 
		b.tags, 
		b.is_active, 
		bv.banner_info, 
		bv.created_at, 
		bv.updated_at 
	FROM 
		banners as b
	JOIN banner_version as bv
	ON 
		b.banner_id = bv.banner_id 
	WHERE 
		b.feature = $1 AND tags @> $2
	ORDER BY
		b.banner_id, bv.created_at, bv.updated_at
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

func (r *bannerRepo) GetBanner(tagIDs []int64, featureID int64, IsAdmin bool) ([]domain.Banner, error) {
	const op = "repository.postgres.GetBannerByFeatureIDandTagsId"

	selectBanner := `
	SELECT 
		b.banner_id, 
		b.feature, 
		b.tags, 
		b.is_active, 
		bv.banner_info, 
		bv.created_at, 
		bv.updated_at 
	FROM 
		banners as b
	JOIN banner_version as bv
	ON 
		b.banner_id = bv.banner_id 
	WHERE
		b.feature = $1 AND tags = $2`

	if !IsAdmin {
		selectBanner += " AND b.is_active IS NOT FALSE"
	}

	selectBanner += " ORDER BY b.banner_id, bv.created_at"

	rows, err := r.db.Query(ctx, selectBanner, featureID, tagIDs)
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

func (r *bannerRepo) UpdateBannerById(banner domain.Banner) error {
	const op = "repository.postgres.UpdateBannerById"

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	const queryBanner = `
		UPDATE 
			banners
			SET 
				feature = $1,
				tags = $2,
				is_active = $3
			WHERE 
				banner_id = $4`

	_, err = r.db.Exec(ctx, queryBanner, banner.FeatureID, banner.TagIds, banner.IsActive, banner.BannerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	const queryBannerVersion = `
		UPDATE 
			banner_version
			SET
				banner_info = $1,
				updated_at = NOW()
			WHERE
				id = (SELECT id
					FROM banner_version
					WHERE banner_id = $2
					GROUP BY id
					ORDER BY max(updated_at) DESC
					LIMIT 1)
	`
	_, err = r.db.Exec(ctx, queryBannerVersion, banner.Content, banner.BannerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *bannerRepo) DeleteBannerById(bannerID int64) error {
	const op = "repository.postgres.DeleteBannerById"

	const query = `
		DELETE FROM banners
		WHERE banners.banner_id = $1
	`
	_, err := r.db.Exec(ctx, query, bannerID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
