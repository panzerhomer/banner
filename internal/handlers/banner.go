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
	GetBanner(ctx context.Context, tagIDs []int64, featureID int64, lastVersion bool, isAdmin bool) ([]domain.Banner, error)
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
		utils.ResponseJSON(w, utils.Error, ErrUserNotAuthorized.Error(), http.StatusUnauthorized)
		return
	}

	if token != AdminToken {
		utils.ResponseJSON(w, utils.Error, ErrUserNotAllowed.Error(), http.StatusForbidden)
		return
	}

	var banner domain.Banner
	if err := json.NewDecoder(r.Body).Decode(&banner); err != nil {
		utils.ResponseJSON(w, utils.Error, ErrBannerIncorrectData.Error(), http.StatusInternalServerError)
		return
	}

	if err := banner.Validate(); err != nil {
		utils.ResponseJSON(w, utils.Error, ErrBannerIncorrectData.Error(), http.StatusBadRequest)
		return
	}

	bannerID, err := h.servo.CreateBanner(ctx, banner)
	if err != nil {
		utils.ResponseJSON(w, utils.Error, ErrBannerExists.Error(), http.StatusInternalServerError)
		return
	}

	utils.ResponseJSON(w, "banner_id", bannerID, http.StatusCreated)
}

func (h *bannerHandler) GetBannersWithFilter(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(string)
	if token == "" && !ok {
		utils.ResponseJSON(w, utils.Error, ErrUserNotAuthorized.Error(), http.StatusUnauthorized)
		return
	}

	if token != AdminToken {
		utils.ResponseJSON(w, utils.Error, ErrUserNotAllowed.Error(), http.StatusForbidden)
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
		utils.ResponseJSON(w, utils.Error, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, banners)
}

func (h *bannerHandler) GetUserBanner(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(string)
	if token == "" && !ok {
		utils.ResponseJSON(w, utils.Error, ErrUserNotAuthorized.Error(), http.StatusUnauthorized)
		return
	}

	var isAdmin bool
	if token == AdminToken {
		isAdmin = true
	}

	if token != UserToken {
		utils.ResponseJSON(w, utils.Error, ErrUserNotAllowed.Error(), http.StatusForbidden)
		return
	}

	queryParams := r.URL.Query()
	featureIdParam := queryParams.Get("feature_id")
	tagsParam := queryParams.Get("tag_ids")
	lastVersionParam := queryParams.Get("use_last_revision")

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

	featureId, _ := strconv.Atoi(featureIdParam)
	lastVersion := strToBool(lastVersionParam)

	banners, err := h.servo.GetBanner(ctx, tags, int64(featureId), lastVersion, isAdmin)
	if err != nil {
		utils.ResponseJSON(w, utils.Error, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, banners)
}

func strToBool(s string) bool {
	return strings.ToLower(s) == "true"
}
