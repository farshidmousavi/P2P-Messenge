package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type BootstrapConfig struct {
	Peers []string `json:"peers"`
}

func LoadBootstrapConfig() (*BootstrapConfig, error) {
	// Config file path in user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".p2p-messenger", "bootstrap.json")
	
	// If file doesn't exist, create a default one
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create directory
		os.MkdirAll(filepath.Dir(configPath), 0700)
		
		// Default config with public bootstrap peers
		defaultConfig := &BootstrapConfig{
			Peers: []string{
				// Public IPFS bootstrap peers (for initial connection)
				"/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
				"/dnsaddr/bootstrap.libp2p.io/p2p/QmQCU2EcMqAqQPR2i9bChDtGNJchTbq5TbXJJ16u19uLTa",
			},
		}
		SaveBootstrapConfig(defaultConfig)
		return defaultConfig, nil
	}

	// Read existing file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config BootstrapConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveBootstrapConfig(config *BootstrapConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".p2p-messenger", "bootstrap.json")
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

// AddBootstrapPeer adds a new bootstrap peer to the list
func AddBootstrapPeer(addr string) error {
	config, err := LoadBootstrapConfig()
	if err != nil {
		return err
	}

	// Check if address already exists
	for _, existing := range config.Peers {
		if existing == addr {
			return nil // Already exists
		}
	}

	config.Peers = append(config.Peers, addr)
	return SaveBootstrapConfig(config)
}