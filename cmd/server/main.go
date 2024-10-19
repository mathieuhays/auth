package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mathieuhays/auth"
	"github.com/mathieuhays/auth/internal/templates"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

var errInvalidPort = errors.New("invalid PORT")

func run(ctx context.Context, getenv func(string) string, stdout io.Writer, stderr io.Writer) error {
	port := getenv("PORT")
	if port == "" {
		return errInvalidPort
	}

	tpl, err := auth.Templates()
	if err != nil {
		return fmt.Errorf("template engine: %s", err)
	}

	tplEngine := templates.NewEngine(tpl)

	server := &http.Server{
		Addr:              net.JoinHostPort("", port),
		Handler:           auth.NewServer(&tplEngine),
		ReadHeaderTimeout: time.Second * 5,
		WriteTimeout:      time.Second * 5,
	}

	serverWG := sync.WaitGroup{}
	serverWG.Add(1)
	serverDone := make(chan struct{}, 1)

	go func() {
		defer serverWG.Done()

		_, _ = fmt.Fprintf(stdout, "Starting server on %s\n", server.Addr)
		if err = server.ListenAndServe(); err != nil {
			_, _ = fmt.Fprintf(stderr, "listen and serve err: %s\n", err)
		}

		serverDone <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		_, _ = fmt.Fprintf(stdout, "graceful shutdown\n")
		if err = server.Shutdown(ctx); err != nil {
			panic(err)
		}
	case <-serverDone:
		_, _ = fmt.Fprintf(stdout, "server has shutdown on its own\n")
	}

	serverWG.Wait()

	return nil
}

func main() {
	fmt.Println("starting...")

	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			log.Fatalf("unexpected error with .env file: %s", err)
		}

		log.Println("no .env file found. skipping")
	}

	if err := run(context.Background(), os.Getenv, os.Stdout, os.Stderr); err != nil {
		log.Fatalf("run error: %s", err)
	}

	log.Printf("Closing...")
}
