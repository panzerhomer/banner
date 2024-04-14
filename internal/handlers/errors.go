package handlers

import "errors"

var (
	ErrBannerIncorrectData = errors.New("incorrect banner data")
	ErrUserNotAuthorized   = errors.New("user not authorized")
	ErrUserNotAllowed      = errors.New("user not allowed")
	ErrBannerExists        = errors.New("banner already exists")
)
