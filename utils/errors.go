package utils

import "errors"

var (
	ErrFileTooLarge     = errors.New("file size exceeds 50MB limit")
	ErrInvalidFileSize  = errors.New("invalid file size")
	ErrInvalidDownloads = errors.New("downloads must be between 1 and 10")
	ErrDownloadsExceeded = errors.New("downloads cannot exceed 10")
	ErrInvalidExpiry    = errors.New("expiry must be at least 1 minute")
	ErrExpiryExceeded   = errors.New("expiry cannot exceed 7 days (10080 minutes)")
)

