package handlers

import (
	"context"
	"encoding/json"
	"net/http"

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
