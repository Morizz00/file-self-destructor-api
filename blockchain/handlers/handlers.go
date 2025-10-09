package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Morizz00/self-destruct-share-api/blockchain/services"
	"github.com/go-chi/chi/v5"
)

type UploadRequest struct {
	FileID        string `json:"file_id"`
	FileData      []byte `json:"file_data"`
	ExpiryMinutes int64  `json:"expiry_minutes"`
	MaxDownloads  int64  `json:"max_downloads"`
}

type UploadResponse struct {
	Success         bool   `json:"success"`
	TransactionHash string `json:"transaction_hash"`
	IPFSHash        string `json:"ipfs_hash"`
	FileHash        string `json:"file_hash"`
	BlockchainURL   string `json:"blockchain_url"`
}

func RegisterUpload(w http.ResponseWriter, r *http.Request) {
	var req UploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	blockchainService := services.NewBlockchainService()
	result, err := blockchainService.RegisterUpload(req.FileID, req.FileData, req.ExpiryMinutes, req.MaxDownloads)
	if err != nil {
		http.Error(w, "Failed to register upload: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := UploadResponse{
		Success:         true,
		TransactionHash: result.TransactionHash,
		IPFSHash:        result.IPFSHash,
		FileHash:        result.FileHash,
		BlockchainURL:   result.BlockchainURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RegisterDownload(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FileID     string `json:"file_id"`
		Downloader string `json:"downloader"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	blockchainService := services.NewBlockchainService()

	result, err := blockchainService.RegisterDownload(req.FileID, req.Downloader)
	if err != nil {
		http.Error(w, "Failed to register download: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func VerifyProof(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "file_id")
	providedHash := r.URL.Query().Get("hash")

	blockchainService := services.NewBlockchainService()

	isValid, err := blockchainService.VerifyFileHash(fileID, providedHash)
	if err != nil {
		http.Error(w, "Failed to verify proof: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"valid":   isValid,
		"file_id": fileID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetProof(w http.ResponseWriter, r *http.Request) {
	fileID := chi.URLParam(r, "file_id")

	blockchainService := services.NewBlockchainService()

	proof, err := blockchainService.GetTransfer(fileID)
	if err != nil {
		http.Error(w, "Failed to get proof: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(proof)
}
