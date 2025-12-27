package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portfolioTUI/config"
	"portfolioTUI/database"
	"portfolioTUI/tui" // Import your local TUI package

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	bm "github.com/charmbracelet/wish/bubbletea"
	lm "github.com/charmbracelet/wish/logging"
)

const (
	host = "0.0.0.0"
	port = 23234
)

func main() {
	// 1. Setup Infrastructure
	config.LoadEnv()
	database.ConnectToDataBase()

	// 2. Configure SSH Server
	s, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", host, port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bm.Middleware(tui.TeaHandler), // Connects Bubble Tea to SSH
			lm.Middleware(),               // Adds logging
		),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// 3. Handle Graceful Shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Starting SSH server on %s:%d\n", host, port)

	// 4. Start Server in a Goroutine
	go func() {
		if err = s.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	// 5. Block until Ctrl+C is pressed
	<-done

	fmt.Println("Stopping SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalln(err)
	}
}
