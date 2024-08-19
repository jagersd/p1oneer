package main

import (
	"log"
	"os"
	"os/signal"
	"p1oneer/internal/pparser"
	"p1oneer/internal/proc"
	"syscall"
)

var signalChannel chan os.Signal

func main() {
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	requests := pparser.ParseConfigFiles()
	for i := 0; i <= 255; i++ {
		if r, ok := requests[uint8(i)]; ok {
			log.Println("Starting", r.Title)
			startProc(r)
		}
	}

	<-signalChannel

	log.Println("Received termination signal, shutting down...")

	os.Exit(0)
}

func startProc(request pparser.StartRequest) {
	var p proc.Proc
	switch request.ReqType {
	case "long":
		go p.StartLong(request.Command, request.Args, signalChannel)
	case "once":
		go p.StartOne(request.Command, request.Args)
	}
}
