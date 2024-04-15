package services

import (
	"context"
	"log"

	"github.com/panzerhomer/banner/internal/domain"
	"github.com/panzerhomer/banner/internal/redis.go"
)

type BannerRepository interface {
	InsertBanner(banner domain.Banner) (int64, error)
	GetBanners(tagIDs []int64, featureID int64, limit int64, offset int64) ([]domain.Banner, error)
	GetBanner(tagIDs []int64, featureID int64, IsAdmin bool) ([]domain.Banner, error)
}

type bannerService struct {
	repo  BannerRepository
	redis *redis.Redis
}

func NewBannerService(repo BannerRepository, redis *redis.Redis) *bannerService {
	return &bannerService{repo: repo, redis: redis}
}

func (s *bannerService) CreateBanner(ctx context.Context, banner domain.Banner) (int64, error) {
	bannerID, err := s.repo.InsertBanner(banner)
	if err != nil {
		return -1, err
	}

	banner.BannerID = bannerID
	s.redis.SaveBanner(bannerID, banner)

	return bannerID, nil
}

func (s *bannerService) GetBanners(ctx context.Context, banner domain.BannerFilter) ([]domain.Banner, error) {
	if banner.Limit <= 0 && banner.Limit > 50 {
		banner.Limit = 10
		banner.Offset = 0
	}

	banners, err := s.repo.GetBanners(banner.TagIds, banner.FeatureID, banner.Limit, banner.Offset)
	if err != nil {
		return nil, err
	}

	return banners, nil
}

func (s *bannerService) GetBanner(ctx context.Context, tagIDs []int64, featureID int64, lastVersion bool, isAdmin bool) ([]domain.Banner, error) {
	if !lastVersion {
		banner, err := s.redis.GetBanner(tagIDs, featureID)
		if err != nil {
			return nil, err
		}
		if banner != nil {
			log.Println("got banner from redis cache")
			return []domain.Banner{*banner}, nil
		}
	}

	banner, err := s.repo.GetBanner(tagIDs, featureID, isAdmin)
	if err != nil {
		return nil, err
	}

	log.Println("got banner from database")

	return banner, nil
}
