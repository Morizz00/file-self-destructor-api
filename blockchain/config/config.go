package config

import (
	"os"
	"strconv"
)

type Config struct {
	BindAddress     string
	EthRPCURL      string
	PrivateKey     string
	ContractAddress string
	ChainID        int64
	IPFSURL        string
}

func Load() *Config {
	return &Config{
		BindAddress:     getEnv("BLOCKCHAIN_PORT", ":3001"),
		EthRPCURL:       getEnv("ETH_RPC_URL", "https://polygon-mumbai.g.alchemy.com/v2/YOUR_KEY"),
		PrivateKey:      getEnv("PRIVATE_KEY", ""),
		ContractAddress: getEnv("CONTRACT_ADDRESS", ""),
		ChainID:         getEnvAsInt("CHAIN_ID", 80001), // polygon mumbai
		IPFSURL:         getEnv("IPFS_URL", "http://localhost:5001"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
