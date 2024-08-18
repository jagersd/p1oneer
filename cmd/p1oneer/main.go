package main

import (
	"log"
	"os"
	"os/signal"
	"p1oneer/internal/pparser"
	"p1oneer/internal/proc"
	"syscall"
)

func main() {
	signalChannel := make(chan os.Signal, 10)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	requests := pparser.ParseConfigFiles()
	for _, r := range requests {
		var p proc.Proc
		go p.StartLong(r.Command, signalChannel)
	}

	<-signalChannel

	log.Println("Received termination signal, shutting down...")

	os.Exit(0)
}
