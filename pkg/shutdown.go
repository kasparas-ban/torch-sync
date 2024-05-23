package pkg

import (
	"os"
	"os/signal"
	"syscall"
)

func GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	defer close(quit)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}