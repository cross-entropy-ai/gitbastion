package server

import (
	"log/slog"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/cross-entropy-ai/gitbastion/internal/keysync"
)

const syncInterval = 5 * time.Minute

type Server struct {
	syncer *keysync.KeySync
}

func New(syncer *keysync.KeySync) *Server {
	return &Server{syncer: syncer}
}

func (s *Server) Run() error {
	if err := s.syncer.Sync(); err != nil {
		return err
	}

	sshd := exec.Command("/usr/sbin/sshd", "-D", "-e")
	sshd.Stdout = os.Stdout
	sshd.Stderr = os.Stderr
	if err := sshd.Start(); err != nil {
		return err
	}
	slog.Info("sshd started")

	go s.syncLoop()

	if err := sshd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.Sys().(syscall.WaitStatus).ExitStatus())
		}
		return err
	}
	return nil
}

func (s *Server) syncLoop() {
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()
	for range ticker.C {
		if err := s.syncer.Sync(); err != nil {
			slog.Error("key sync failed", "error", err)
		}
	}
}
