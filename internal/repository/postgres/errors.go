package repository

import "errors"

var (
	ErrBannerNotFound = errors.New("banner not found")
	ErrBannerExists   = errors.New("banner exists")
)
