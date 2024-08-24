package proc

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type channelHub struct {
	signalChannel  chan (os.Signal)
	processChannel chan (*os.Process)
	processes      []*os.Process
}

var hub channelHub

func (hub *channelHub) stopAllProcesses() {
	for _, proc := range hub.processes {
		log.Println("Stopping ", proc.Pid)
		if err := proc.Signal(syscall.SIGTERM); err != nil {
			if err == os.ErrProcessDone {
				log.Printf("%d has already exited", proc.Pid)
			} else {
				log.Printf("Failed to send SIGTERM to %d: %v", proc.Pid, err)
			}
		}
		syscall.Kill(proc.Pid, syscall.SIGTERM)
	}
	os.Exit(0)
}

func StartProcessHub() {
	hub = channelHub{
		signalChannel:  make(chan os.Signal, 1),
		processChannel: make(chan *os.Process, 1),
	}
	signal.Notify(hub.signalChannel, syscall.SIGINT, syscall.SIGTERM)
}

func Monitor() {
	startMonitorRoutine()
}

func startMonitorRoutine() {
	go reapZombies()
	for {
		select {
		case proc := <-hub.processChannel:
			hub.processes = append(hub.processes, proc)
			log.Println("Pid started:", proc.Pid)
		case <-hub.signalChannel:
			log.Println("Received signal, stopping all processes")
			hub.stopAllProcesses()
		}
	}
}

func reapZombies() {
	reaper := reaper{}
	for {
		if err := reaper.scan(); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Second * 3)
	}
}
