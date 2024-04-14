package services

import (
	"context"

	"github.com/panzerhomer/banner/internal/domain"
)

// import "github.com/panzerhomer/banner/internal/storage/postgres"

type BannerRepository interface {
	InsertBanner(banner domain.Banner) (int64, error)
	GetBanners(tagIDs []int, featureID int, limit int, offset int) ([]domain.Banner, error)
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
