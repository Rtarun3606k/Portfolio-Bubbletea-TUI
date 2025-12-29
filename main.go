package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portfolioTUI/config"
	"portfolioTUI/database"
	"portfolioTUI/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

const (
	host = "0.0.0.0"
	port = 23234
)

func main() {
	// 1. Infrastructure
	config.LoadEnv()
	database.ConnectToDataBase()
	go tui.GetOrFetchData()

	// 2. SSH Keys
	keyPath := ".ssh/term_info_ed25519"
	if _, err := os.Stat(".ssh"); os.IsNotExist(err) {
		_ = os.Mkdir(".ssh", 0700)
	}

	// 3. Configure Server
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, fmt.Sprintf("%d", port))),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			// --- SIMPLIFIED COLOR FIX ---
			bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
				// 1. Force Lipgloss to use 256 colors
				// This overrides the auto-detection which fails over SSH
				lipgloss.SetColorProfile(termenv.ANSI256)

				// 2. Call your standard handler
				return tui.TeaHandler(s)
			}),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// 4. Start
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Starting SSH server on %s:%d", host, port)

	go func() {
		if err = s.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	<-done
	log.Println("Stopping SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}

