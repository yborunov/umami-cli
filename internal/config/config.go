package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
}

func Load(flagEndpoint, flagToken string) (*Config, error) {
	cfg := &Config{}
	if err := cfg.read(); err != nil {
		return nil, err
	}

	if flagEndpoint != "" {
		cfg.Endpoint = flagEndpoint
	}
	if flagToken != "" {
		cfg.Token = flagToken
	}

	if cfg.Endpoint == "" {
		return nil, errors.New("missing endpoint: set --endpoint or UMAMI_URL")
	}

	cfg.Endpoint = normalizeEndpoint(cfg.Endpoint)
	return cfg, nil
}

func (c *Config) Save() error {
	path, err := path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func (c *Config) read() error {
	path, err := path()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("invalid config file: %w", err)
	}
	return nil
}

func path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "umami-cli", "config.json"), nil
}

func normalizeEndpoint(endpoint string) string {
	endpoint = strings.TrimRight(endpoint, "/")
	if !strings.HasSuffix(endpoint, "/api") {
		endpoint += "/api"
	}
	return endpoint
}
