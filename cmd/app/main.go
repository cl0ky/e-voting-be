package main

import (
	"context"
	"fmt"
	"github/com/cl0ky/e-voting-be/config"
	"github/com/cl0ky/e-voting-be/config/seed"
	"github/com/cl0ky/e-voting-be/env"
	"github/com/cl0ky/e-voting-be/server/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadEnv()
	env.GetEnv()
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	gin.SetMode(env.GinMode)
	port := env.Port
	app := gin.Default()

	db := config.NewDB(ctx)
	defer db.Close()

	if err := seed.SeedRT(db.Instance()); err != nil {
		log.Fatalf("Gagal seeding RT: %v", err)
	}

	router.SetupRoutes(router.SetupRoutesConfig{
		Router: app,
		DB:     db.Instance(),
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Server running on PORT %d", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-ctx.Done()
	log.Println("Shutting down server...")
	server.Shutdown(context.Background())
}
