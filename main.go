package main

import (
	"log/slog"
	"os"

	"github.com/cross-entropy-ai/gitbastion/internal/config"
	"github.com/cross-entropy-ai/gitbastion/internal/keysync"
	"github.com/cross-entropy-ai/gitbastion/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfgPath := os.Getenv("CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "/etc/gitbastion/config.yaml"
	}
	cfg := config.Load(cfgPath)

	srv := server.New(keysync.New(cfg))
	if err := srv.Run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}
