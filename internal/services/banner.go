package services

import (
	"context"

	"github.com/panzerhomer/banner/internal/domain"
)

// import "github.com/panzerhomer/banner/internal/storage/postgres"

type BannerRepository interface {
	InsertBanner(banner domain.Banner) (int64, error)
	GetBanners(tagIDs []int64, featureID int64, limit int64, offset int64) ([]domain.Banner, error)
}

type bannerService struct {
	repo BannerRepository
}

func NewBannerService(repo BannerRepository) *bannerService {
	return &bannerService{repo: repo}
}

func (s *bannerService) CreateBanner(ctx context.Context, banner domain.Banner) (int64, error) {
	bannerID, err := s.repo.InsertBanner(banner)
	if err != nil {
		return -1, err
	}
	return bannerID, nil
}

func (s *bannerService) GetBanners(ctx context.Context, banner domain.BannerFilter) ([]domain.Banner, error) {
	if banner.Limit == 0 {
		banner.Limit = 30
		banner.Offset = 0
	}

	banners, err := s.repo.GetBanners(banner.TagIds, banner.FeatureID, banner.Limit, banner.Offset)
	if err != nil {
		return nil, err
	}

	return banners, nil
}
