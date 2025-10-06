package models

import (
	"time"
)

// FileTransfer represents a blockchain record of file transfer
type FileTransfer struct {
	FileHash         string    `json:"file_hash"`
	IPFSHash         string    `json:"ipfs_hash"`
	Uploader         string    `json:"uploader"`
	UploadTimestamp  time.Time `json:"upload_timestamp"`
	ExpiryTimestamp  time.Time `json:"expiry_timestamp"`
	IsDownloaded     bool      `json:"is_downloaded"`
	Downloader       string    `json:"downloader"`
	DownloadTimestamp time.Time `json:"download_timestamp"`
	IsSelfDestructed bool      `json:"is_self_destructed"`
	MaxDownloads     int64     `json:"max_downloads"`
	DownloadCount    int64     `json:"download_count"`
}

// UploadResult represents the result of a blockchain upload registration
type UploadResult struct {
	Success         bool   `json:"success"`
	TransactionHash string `json:"transaction_hash"`
	IPFSHash        string `json:"ipfs_hash"`
	FileHash        string `json:"file_hash"`
	BlockchainURL   string `json:"blockchain_url"`
}

// DownloadResult represents the result of a blockchain download registration
type DownloadResult struct {
	Success         bool   `json:"success"`
	TransactionHash string `json:"transaction_hash"`
	DownloadCount   int64  `json:"download_count"`
	BlockchainURL   string `json:"blockchain_url"`
}

// VerificationResult represents the result of file hash verification
type VerificationResult struct {
	Valid   bool   `json:"valid"`
	FileID  string `json:"file_id"`
	Message string `json:"message"`
}
