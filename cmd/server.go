package main

import (
	"fmt"
	"os"
	"os/signal"
	"spolks/internal/service"
	"syscall"
)

func main() {
	server := service.NewServer()
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		<-sigCh

		fmt.Println("Shutting down server")
		server.Shutdown()
		//todo impl it
		os.Exit(1)
	}()

	server.Start("tcp", ":8080")

}
