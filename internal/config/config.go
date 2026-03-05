package config

import (
	"fmt"
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

		if err := Validate(config); err != nil {
			return config, err
		}
	}

	return config, nil
}

func Save(config *Config) error {
	configPath := getConfigPath()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func getConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil || homeDir == "" {
		if cwd, cwdErr := os.Getwd(); cwdErr == nil {
			return filepath.Join(cwd, "config.yaml")
		}
		return "config.yaml"
	}

	return filepath.Join(homeDir, ".config", "nocta", "config.yaml")
}

func Validate(config *Config) error {
	if config.GUI.Width < 320 || config.GUI.Height < 240 {
		return fmt.Errorf("invalid GUI size: %dx%d", config.GUI.Width, config.GUI.Height)
	}

	if config.GUI.Split <= 0 || config.GUI.Split >= 1 {
		return fmt.Errorf("invalid split value: %v", config.GUI.Split)
	}

	if config.Refresh.Interval < 1 || config.Refresh.Interval > 3600 {
		return fmt.Errorf("invalid refresh interval: %d", config.Refresh.Interval)
	}

	switch config.Filter.DefaultProtocol {
	case "ALL", "TCP", "UDP":
		return nil
	default:
		return fmt.Errorf("invalid default protocol: %s", config.Filter.DefaultProtocol)
	}
}
