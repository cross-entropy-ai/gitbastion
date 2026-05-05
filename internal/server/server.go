package server

import (
	"bufio"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"strings"
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
	stderrPipe, err := sshd.StderrPipe()
	if err != nil {
		return err
	}
	if err := sshd.Start(); err != nil {
		return err
	}
	go forwardLogs(stderrPipe)
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

func forwardLogs(r io.Reader) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With("source", "sshd")
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		logger.Info(line)
	}
	if err := scanner.Err(); err != nil {
		slog.Error("sshd log forwarding failed", "error", err)
	}
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
