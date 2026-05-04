package config

import (
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AllowedUsers  []string `yaml:"allowed_users"`
	AllowedGroups []string `yaml:"allowed_groups"`
	GHToken       string   `yaml:"-"`
	KeysFile      string   `yaml:"-"`
}

func Load(path string) Config {
	cfg := Config{
		KeysFile: "/home/git/.ssh/authorized_keys",
		GHToken:  os.Getenv("GH_TOKEN"),
	}

	// YAML file
	if data, err := os.ReadFile(path); err == nil {
		var fileCfg Config
		if err := yaml.Unmarshal(data, &fileCfg); err != nil {
			slog.Warn("failed to parse config file", "path", path, "error", err)
		} else {
			cfg.AllowedUsers = append(cfg.AllowedUsers, fileCfg.AllowedUsers...)
			cfg.AllowedGroups = append(cfg.AllowedGroups, fileCfg.AllowedGroups...)
			slog.Info("loaded config file", "path", path,
				"users", len(fileCfg.AllowedUsers),
				"groups", len(fileCfg.AllowedGroups))
		}
	} else {
		slog.Info("no config file found, skipping", "path", path)
	}

	// Environment variables
	if env := os.Getenv("ALLOWED_USERS"); env != "" {
		for _, u := range strings.Split(env, ",") {
			if u = strings.TrimSpace(u); u != "" {
				cfg.AllowedUsers = append(cfg.AllowedUsers, u)
			}
		}
		slog.Info("loaded ALLOWED_USERS from env")
	}
	if env := os.Getenv("ALLOWED_GROUPS"); env != "" {
		for _, g := range strings.Split(env, ",") {
			if g = strings.TrimSpace(g); g != "" {
				cfg.AllowedGroups = append(cfg.AllowedGroups, g)
			}
		}
		slog.Info("loaded ALLOWED_GROUPS from env")
	}

	cfg.AllowedUsers = dedupe(cfg.AllowedUsers)
	cfg.AllowedGroups = dedupe(cfg.AllowedGroups)

	slog.Info("config loaded",
		"total_users", len(cfg.AllowedUsers),
		"total_groups", len(cfg.AllowedGroups))
	return cfg
}

func dedupe(s []string) []string {
	seen := make(map[string]struct{}, len(s))
	out := s[:0]
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}
