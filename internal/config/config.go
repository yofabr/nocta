package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	GUI struct {
		Width  int     `yaml:"width"`
		Height int     `yaml:"height"`
		Split  float64 `yaml:"split"`
	} `yaml:"gui"`
	Refresh struct {
		Interval int  `yaml:"interval"` // seconds
		Auto     bool `yaml:"auto"`
	} `yaml:"refresh"`
	Filter struct {
		DefaultProtocol string `yaml:"default_protocol"`
		ShowClosed      bool   `yaml:"show_closed"`
	} `yaml:"filter"`
	Notifications struct {
		Enabled bool `yaml:"enabled"`
		OnNew   bool `yaml:"on_new"`
		OnClose bool `yaml:"on_close"`
	} `yaml:"notifications"`
}

func Load() (*Config, error) {
	config := &Config{}

	// Set default values
	config.GUI.Width = 900
	config.GUI.Height = 520
	config.GUI.Split = 0.35
	config.Refresh.Interval = 30
	config.Refresh.Auto = false
	config.Filter.DefaultProtocol = "ALL"
	config.Filter.ShowClosed = false
	config.Notifications.Enabled = false
	config.Notifications.OnNew = true
	config.Notifications.OnClose = true

	// Try to load from file
	configPath := getConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return config, err
		}

		if err := yaml.Unmarshal(data, config); err != nil {
			return config, err
		}
	}

	return config, nil
}

func Save(config *Config) error {
	configPath := getConfigPath()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "nocta", "config.yaml")
}
