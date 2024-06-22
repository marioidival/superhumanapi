package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/marioidival/superhuman-api/internal/handlers"
	"github.com/marioidival/superhuman-api/pkg/database"
	"github.com/peterbourgon/ff"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err.Error())
	}
}

func run() error {
	fs := flag.NewFlagSet("api", flag.ExitOnError)

	var databaseURL string
	var clearbitAPIKey string
	var appPort int

	fs.StringVar(&databaseURL, "database-url", "", "e.g., postgres://username:password@localhost:5432/database_name")
	fs.StringVar(&clearbitAPIKey, "clearbit-api-key", "", "APIKEY for Clearbit")
	fs.IntVar(&appPort, "app-port", 3000, "Port for the application")

	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix())
	if err != nil {
		return err
	}

	ctx := context.Background()
	e := echo.New()

	dbc, err := database.Open(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer dbc.Close()

	handlers := handlers.NewServer(dbc, clearbitAPIKey)
	e.GET("/lookup", handlers.EmailLookupHandler)
	e.GET("/popularity", handlers.PopularityHandler)
	e.GET("/debug/*", echo.WrapHandler(http.DefaultServeMux))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", appPort),
		Handler: e,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	serverErrors := make(chan error, 1)
	go func() {
		log.Println("startup job processing system api", "PORT", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case serverError := <-serverErrors:
		return errors.Unwrap(serverError)

	case sig := <-quit:
		log.Println("Server is shutting down", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if shutdownErr := server.Shutdown(ctx); shutdownErr != nil {
			defer func() {
				closeErr := server.Close()
				if closeErr != nil {
					log.Fatalln("Could not close server", closeErr)
				}
			}()
			log.Fatalln("Could not gracefully shutdown the server")
		}
		close(done)
	case <-done:
		return nil
	}

	return nil
}
