package proc

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type channelHub struct {
	SignalChannel  chan (os.Signal)
	ProcessChannel chan (*os.Process)
}

var hub channelHub

func StartProcessHub() {
	hub = channelHub{
		SignalChannel:  make(chan os.Signal),
		ProcessChannel: make(chan *os.Process),
	}
	signal.Notify(hub.SignalChannel, syscall.SIGINT, syscall.SIGTERM)
	go startMonitorRoutine()
}

func Monitor() {
	<-hub.SignalChannel
}

func startMonitorRoutine() {
	for {
		select {
		case proc := <-hub.ProcessChannel:
			fmt.Println("Pid started:", proc.Pid)
		case <-hub.SignalChannel:
			fmt.Println("Signal received")
		}
	}
}
