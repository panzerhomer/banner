package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"
	"github.com/panzerhomer/banner/internal/domain"
	"github.com/panzerhomer/banner/internal/utils"
)

var ctx = context.Background()

const (
	AdminToken = "admin_token"
	UserToken  = "user_token"
)

type BannerService interface {
	CreateBanner(ctx context.Context, banner domain.Banner) (int64, error)
	GetBanners(ctx context.Context, banner domain.BannerFilter) ([]domain.Banner, error)
}

type bannerHandler struct {
	servo BannerService
}

func NewBannerHandler(servo BannerService) *bannerHandler {
	return &bannerHandler{servo: servo}
}

func (h *bannerHandler) CreateBanner(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(string)
	if token == "" && !ok {
		utils.ResponseJSON(w, utils.Error, ErrUserNotAuthorized, http.StatusUnauthorized)
		return
	}

	if token != AdminToken {
		utils.ResponseJSON(w, utils.Error, ErrUserNotAllowed, http.StatusForbidden)
		return
	}

	var banner domain.Banner
	if err := json.NewDecoder(r.Body).Decode(&banner); err != nil {
		utils.ResponseJSON(w, utils.Error, ErrBannerIncorrectData, http.StatusInternalServerError)
		return
	}

	if err := banner.Validate(); err != nil {
		utils.ResponseJSON(w, utils.Error, ErrBannerIncorrectData, http.StatusBadRequest)
		return
	}

	bannerID, err := h.servo.CreateBanner(ctx, banner)
	if err != nil {
		utils.ResponseJSON(w, utils.Error, ErrBannerExists, http.StatusInternalServerError)
		return
	}

	utils.ResponseJSON(w, "banner_id", bannerID, http.StatusCreated)
}

func (h *bannerHandler) GetBannersWithFilter(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(string)
	if token == "" && !ok {
		utils.ResponseJSON(w, utils.Error, ErrUserNotAuthorized, http.StatusUnauthorized)
		return
	}

	if token != AdminToken {
		utils.ResponseJSON(w, utils.Error, ErrUserNotAllowed, http.StatusForbidden)
		return
	}

	queryParams := r.URL.Query()

	featureIdParam := queryParams.Get("feature_id")
	tagsParam := queryParams.Get("tag_ids")
	limitParam := queryParams.Get("limit")
	offsetParam := queryParams.Get("offset")

	tagValues := strings.Split(tagsParam, ",")

	var tags []int64
	for _, tagValue := range tagValues {
		if tagValue == "" {
			continue
		}
		tag, err := strconv.Atoi(tagValue)
		if err != nil {
			utils.ResponseJSON(w, utils.Error, err.Error(), http.StatusBadRequest)
			return
		}
		tags = append(tags, int64(tag))
	}

	limit, _ := strconv.Atoi(limitParam)
	offset, _ := strconv.Atoi(offsetParam)
	featureId, _ := strconv.Atoi(featureIdParam)

	bannerWithFilter := domain.BannerFilter{
		TagIds:    tags,
		FeatureID: int64(featureId),
		Limit:     int64(limit),
		Offset:    int64(offset),
	}

	banners, err := h.servo.GetBanners(ctx, bannerWithFilter)
	if err != nil {
		utils.ResponseJSON(w, utils.Error, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, banners)
}
