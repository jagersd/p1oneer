package proc

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Proc struct {
	commandName string
	args        []string
}

func (P *Proc) StartLong(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start %s with Err: %v", P.commandName, err)
	}

	hub.processChannel <- cmd.Process

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf("%s exited with error:%d %v", P.commandName, err.(*exec.ExitError).ExitCode(), err)
		} else {
			log.Println(P.commandName, " exited successfully")
		}
		hub.signalChannel <- syscall.SIGTERM
	}()
}

func (P *Proc) StartOne(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("command %s failed with error: %v", command, err)
	}
}
