package main

import (
	"log"
	"os"
	"p1oneer/internal/pparser"
	"p1oneer/internal/proc"
)

func main() {
	proc.StartProcessHub()

	requests := pparser.ParseConfigFiles()
	for i := 0; i <= 255; i++ {
		if r, ok := requests[uint8(i)]; ok {
			log.Println("Starting", r.Title)
			startProc(r)
		}
	}

	proc.Monitor()

	log.Println("Received termination signal, shutting down...")

	os.Exit(0)
}

func startProc(request pparser.StartRequest) {
	var p proc.Proc
	switch request.ReqType {
	case "long":
		go p.StartLong(request.Command, request.Args)
	case "once":
		go p.StartOne(request.Command, request.Args)
	}
}
