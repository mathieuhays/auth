package main

import (
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
	"time"
)

var errInvalidPort = errors.New("invalid PORT")

func run(getenv func(string) string, stdout io.Writer) error {
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

	_, _ = fmt.Fprintf(stdout, "Starting server on %s\n", server.Addr)
	if err = server.ListenAndServe(); err != nil {
		return err
	}

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

	if err := run(os.Getenv, os.Stdout); err != nil {
		log.Fatalf("run error: %s", err)
	}

	log.Printf("Closing...")
}
