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
		client: github.NewClient(cfg.GHToken),
	}
}

func (s *KeySync) Sync() error {
	slog.Info("syncing authorized keys")
	users := make(map[string]struct{})

	// Collect users from allowed_users
	for _, u := range s.cfg.AllowedUsers {
		users[u] = struct{}{}
	}

	// Collect users from allowed_groups (org/team-slug)
	for _, group := range s.cfg.AllowedGroups {
		org, slug, ok := strings.Cut(group, "/")
		if !ok {
			slog.Error("invalid group format, expected org/team-slug", "group", group)
			continue
		}
		if s.cfg.GHToken == "" {
			slog.Warn("GH_TOKEN required for team member lookup, skipping", "group", group)
			continue
		}
		members, err := s.client.TeamMembers(org, slug)
		if err != nil {
			slog.Error("failed to fetch team members", "group", group, "error", err)
			continue
		}
		slog.Info("fetched team members", "group", group, "count", len(members))
		for _, m := range members {
			users[m] = struct{}{}
		}
	}

	// Fetch keys for all unique users
	var lines []string
	for user := range users {
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
