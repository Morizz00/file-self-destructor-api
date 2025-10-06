package services

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Morizz00/self-destruct-share-api/blockchain/config"
	"github.com/Morizz00/self-destruct-share-api/blockchain/models"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-ipfs-api"
)

type BlockchainService struct {
	config     *config.Config
	ethClient  *ethclient.Client
	ipfsClient *shell.Shell
}

func NewBlockchainService() *BlockchainService {
	cfg := config.Load()
	
	// Initialize Ethereum client
	client, err := ethclient.Dial(cfg.EthRPCURL)
	if err != nil {
		log.Printf("Failed to connect to Ethereum client: %v", err)
		// Continue without blockchain for now
		client = nil
	}

	// Initialize IPFS client
	ipfsClient := shell.NewShell(cfg.IPFSURL)

	return &BlockchainService{
		config:     cfg,
		ethClient:  client,
		ipfsClient: ipfsClient,
	}
}

// RegisterUpload registers a file upload on the blockchain
func (bs *BlockchainService) RegisterUpload(fileID string, fileData []byte, expiryMinutes int64, maxDownloads int64) (*models.UploadResult, error) {
	// 1. Calculate file hash
	fileHash := bs.calculateFileHash(fileData)
	
	// 2. Create metadata
	metadata := map[string]interface{}{
		"file_id":         fileID,
		"file_hash":       fileHash,
		"timestamp":       time.Now().UTC().Format(time.RFC3339),
		"expiry_minutes":  expiryMinutes,
		"max_downloads":   maxDownloads,
		"file_size":       len(fileData),
	}

	// 3. Upload metadata to IPFS
	ipfsHash, err := bs.uploadToIPFS(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to IPFS: %v", err)
	}

	// 4. For now, simulate blockchain transaction
	// In a real implementation, you would interact with a smart contract
	txHash := bs.simulateBlockchainTransaction(fileID, fileHash, ipfsHash)

	result := &models.UploadResult{
		Success:         true,
		TransactionHash: txHash,
		IPFSHash:        ipfsHash,
		FileHash:        fileHash,
		BlockchainURL:   fmt.Sprintf("https://mumbai.polygonscan.com/tx/%s", txHash),
	}

	return result, nil
}

// RegisterDownload registers a file download on the blockchain
func (bs *BlockchainService) RegisterDownload(fileID string, downloader string) (*models.DownloadResult, error) {
	// Simulate download registration
	txHash := bs.simulateBlockchainTransaction(fileID, "download", downloader)

	result := &models.DownloadResult{
		Success:         true,
		TransactionHash: txHash,
		DownloadCount:   1, // This would be fetched from blockchain
		BlockchainURL:   fmt.Sprintf("https://mumbai.polygonscan.com/tx/%s", txHash),
	}

	return result, nil
}

// VerifyFileHash verifies if a provided hash matches the stored hash
func (bs *BlockchainService) VerifyFileHash(fileID string, providedHash string) (bool, error) {
	// In a real implementation, you would query the blockchain
	// For now, we'll simulate verification
	log.Printf("Verifying file %s with hash %s", fileID, providedHash)
	
	// Simulate verification logic
	// In reality, you'd query the smart contract
	return true, nil
}

// GetTransfer retrieves transfer information from blockchain
func (bs *BlockchainService) GetTransfer(fileID string) (*models.FileTransfer, error) {
	// In a real implementation, you would query the blockchain
	// For now, return a mock response
	transfer := &models.FileTransfer{
		FileHash:         "mock_hash_" + fileID,
		IPFSHash:         "mock_ipfs_hash",
		Uploader:         "0x1234567890123456789012345678901234567890",
		UploadTimestamp:  time.Now().Add(-1 * time.Hour),
		ExpiryTimestamp:  time.Now().Add(24 * time.Hour),
		IsDownloaded:     false,
		Downloader:       "",
		DownloadTimestamp: time.Time{},
		IsSelfDestructed: false,
		MaxDownloads:     5,
		DownloadCount:    0,
	}

	return transfer, nil
}

// Helper methods

func (bs *BlockchainService) calculateFileHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (bs *BlockchainService) uploadToIPFS(data interface{}) (string, error) {
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// Upload to IPFS
	hash, err := bs.ipfsClient.Add(jsonData)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (bs *BlockchainService) simulateBlockchainTransaction(fileID string, data ...string) string {
	// Generate a mock transaction hash
	combined := fileID
	for _, d := range data {
		combined += d
	}
	
	hash := sha256.Sum256([]byte(combined))
	return hex.EncodeToString(hash[:])
}

// GetIPFSData retrieves data from IPFS
func (bs *BlockchainService) GetIPFSData(hash string) ([]byte, error) {
	return bs.ipfsClient.Cat(hash)
}
