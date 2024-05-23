package main

import (
	"fmt"
	"os"
	"time"
	"torch/torch-sync/config"
	"torch/torch-sync/pkg"
	"torch/torch-sync/storage"
	"torch/torch-sync/websockets"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// setup exit code for graceful shutdown
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	// load config
	env, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	// run the server
	cleanup, err := run(env)

	defer cleanup()
	if err != nil {
		fmt.Printf("error: %v", err)
		exitCode = 1
		return
	}

	pkg.GracefulShutdown()
}

func run(env config.EnvVars) (func(), error) {
	app, cleanup, err := buildServer(env)
	if err != nil {
		return nil, err
	}

	// start the server
	go func() {
		app.Listen(":" + env.PORT)
	}()

	return func() {
		cleanup()
		app.Shutdown()
	}, nil
}

func buildServer(env config.EnvVars) (*fiber.App, func(), error) {
	// init storage
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", env.DB_USERNAME, env.DB_PASSWORD, env.DB_HOSTNAME, env.DB_PORT, env.DB_NAME)
	db, err := storage.InitDB(connStr, 30 * time.Second)
	if err != nil {
		return nil, nil, err
	}

	// create fiber app
	app := fiber.New()

	// TODO: add auth middleware
	app.Use("/sync", websockets.WebsocketsMiddleware)

	// add health check
	app.Get("/health", func (c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	app.Get("/sync", websocket.New(websockets.SyncHandler))

	return app, func() {
		storage.CloseDB(db)
	}, nil
}