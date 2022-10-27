package main

import (
	"os"
	"os/signal"
	"spolks/internal/server"
	"syscall"
)

func main() {
	server := server.NewServer("tcp", "0.0.0.0:8080")
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		<-sigCh

		server.Shutdown()
	}()

	server.Start()

}
