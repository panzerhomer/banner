package domain

import (
	"errors"
	"time"
)

type Banner struct {
	BannerID  int64     `json:"banner_id,omitempty"`
	TagIds    []int64   `json:"tag_ids,omitempty"`
	FeatureID int64     `json:"feature_id,omitempty"`
	Content   any       `json:"content,omitempty"`
	IsActive  bool      `json:"is_active,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type BannerRequest struct {
	TagIds    []int64 `json:"tag_ids,omitempty"`
	FeatureID int64   `json:"feature_id,omitempty"`
	Content   any     `json:"content,omitempty"`
	IsActive  bool    `json:"is_active,omitempty"`
}

func (b *Banner) Validate() error {
	if len(b.TagIds) == 0 {
		return errors.New("tag_ids cannot be empty")
	}

	if b.FeatureID < 0 {
		return errors.New("feature_id must be a positive integer")
	}

	return nil
}

type BannerFilter struct {
	TagIds        []int64
	FeatureID     int64
	Limit, Offset int64
}
