package keysync

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/cross-entropy-ai/gitbastion/internal/config"
	"github.com/cross-entropy-ai/gitbastion/internal/github"
)

type KeySync struct {
	cfg    config.Config
	client *github.Client
}

func New(cfg config.Config) *KeySync {
	return &KeySync{
		cfg:    cfg,
		client: github.NewClient(),
	}
}

func (s *KeySync) Sync() error {
	slog.Info("syncing authorized keys")
	var lines []string

	for _, user := range s.cfg.AllowedUsers {
		keys, err := s.client.UserKeys(user)
		if err != nil {
			slog.Error("failed to fetch user keys", "user", user, "error", err)
			continue
		}
		slog.Info("fetched user keys", "user", user, "count", len(keys))
		for _, k := range keys {
			lines = append(lines, k.Key+" "+user)
		}
	}

	if err := os.WriteFile(s.cfg.KeysFile, []byte(strings.Join(lines, "\n")+"\n"), 0600); err != nil {
		return fmt.Errorf("write keys file: %w", err)
	}
	if err := os.Chown(s.cfg.KeysFile, 1000, 1000); err != nil {
		return fmt.Errorf("chown keys file: %w", err)
	}

	slog.Info("key sync complete", "count", len(lines))
	return nil
}
