package config

import (
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AllowedUsers []string `yaml:"allowed_users"`
	KeysFile     string   `yaml:"-"`
}

func Load(path string) Config {
	cfg := Config{
		KeysFile: "/home/git/.ssh/authorized_keys",
	}

	// YAML file
	if data, err := os.ReadFile(path); err == nil {
		var fileCfg Config
		if err := yaml.Unmarshal(data, &fileCfg); err != nil {
			slog.Warn("failed to parse config file", "path", path, "error", err)
		} else {
			cfg.AllowedUsers = append(cfg.AllowedUsers, fileCfg.AllowedUsers...)
			slog.Info("loaded config file", "path", path, "users", len(fileCfg.AllowedUsers))
		}
	} else {
		slog.Info("no config file found, skipping", "path", path)
	}

	// Environment variable
	if env := os.Getenv("ALLOWED_USERS"); env != "" {
		for _, u := range strings.Split(env, ",") {
			if u = strings.TrimSpace(u); u != "" {
				cfg.AllowedUsers = append(cfg.AllowedUsers, u)
			}
		}
		slog.Info("loaded ALLOWED_USERS from env")
	}

	// Deduplicate
	seen := make(map[string]struct{})
	unique := cfg.AllowedUsers[:0]
	for _, u := range cfg.AllowedUsers {
		if _, ok := seen[u]; !ok {
			seen[u] = struct{}{}
			unique = append(unique, u)
		}
	}
	cfg.AllowedUsers = unique

	slog.Info("config loaded", "total_users", len(cfg.AllowedUsers))
	return cfg
}
