package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	t.Run("fail on missing PORT", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(1)

		errCh := make(chan error)

		go func() {
			defer wg.Done()

			err := run(ctx, func(s string) string {
				return ""
			}, os.Stdout, os.Stderr)
			errCh <- err
		}()

		var err error

		select {
		case <-time.Tick(time.Millisecond * 10):
			cancel()
			t.Fatalf("test timeout, expected error")
		case err = <-errCh:
		}

		if err == nil {
			t.Fatalf("expected error but none were returned")
		}

		if !errors.Is(err, errInvalidPort) {
			t.Fatalf("unexpected error. expected: %s. got: %s", errInvalidPort, err)
		}
	})

	t.Run("serves requests", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		wg := sync.WaitGroup{}
		wg.Add(1)

		errCh := make(chan error)
		port := "12345"

		go func() {
			defer wg.Done()

			err := run(ctx, func(s string) string {
				return port
			}, os.Stdout, os.Stderr)
			errCh <- err
		}()

		// wait for server to be ready
		// there might be a better way to do this?
		time.Sleep(time.Millisecond * 500)

		var err error

		baseUrl := fmt.Sprintf("http://127.0.0.1:%s", port)
		assertRequestStatusCode(t, baseUrl+"/", http.StatusOK)
		assertRequestStatusCode(t, baseUrl+"/lksjdf", http.StatusNotFound)

		cancel()

		select {
		case <-time.Tick(time.Millisecond * 300):
			panic("graceful shutdown failed?")
		case err = <-errCh:
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				t.Fatalf("unexpected server err: %s", err)
			}
		}

		log.Println("shutting down")
	})
}

func assertRequestStatusCode(t testing.TB, url string, statusCode int) {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatalf("error initialising test request: %s", err)
	}

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("error executing test request: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != statusCode {
		t.Errorf("unexpected status code. expected: %d. got: %d", statusCode, res.StatusCode)
	}
}
