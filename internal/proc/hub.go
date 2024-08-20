package proc

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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
		signalChannel:  make(chan os.Signal),
		processChannel: make(chan *os.Process),
	}
	signal.Notify(hub.signalChannel, syscall.SIGINT, syscall.SIGTERM)
}

func Monitor() {
	startMonitorRoutine()
}

func startMonitorRoutine() {
	for {
		select {
		case proc := <-hub.processChannel:
			hub.processes = append(hub.processes, proc)
			fmt.Println("Pid started:", proc.Pid)
		case <-hub.signalChannel:
			fmt.Println("Received signal, stopping all processes")
			hub.stopAllProcesses()
		}
	}
}
