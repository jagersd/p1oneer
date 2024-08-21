package main

import (
	"log"
	"os"

	"github.com/jagersd/p1oneer/internal/pparser"
	"github.com/jagersd/p1oneer/internal/proc"
)

func main() {
	proc.StartProcessHub()

	requests := pparser.ParseConfigFiles()
	for i := 0; i <= 255; i++ {
		if r, ok := requests[uint8(i)]; ok {
			log.Println("Starting", r.Title)
			startProcess(r)
		}
	}

	proc.Monitor()

	log.Println("Received termination signal, shutting down p1oneer")

	os.Exit(0)
}

func startProcess(request pparser.StartRequest) {
	p := proc.NewProcessWrapper(
		request.Title,
		request.Command,
		request.Args,
	)

	switch request.ReqType {
	case "long":
		go p.StartLong()
	case "once":
		p.StartOne()
	case "before":
		p.StartBefore()
	}
}
