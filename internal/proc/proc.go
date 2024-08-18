package proc

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type Proc struct {
	commandName string
	args        []string
}

func (P *Proc) StartLong(command string, signalChannel chan (os.Signal)) {
	P.generateStartCommand(command)
	cmd := exec.Command(P.commandName, P.args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start %s with Err: %v", P.commandName, err)
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf(P.commandName+" exited with error: %v", err)
		} else {
			log.Println(P.commandName, " exited successfully")
		}
		os.Exit(0)
	}()

	<-signalChannel

	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		log.Printf("Failed to send SIGTERM to %s: %v", P.commandName, err)
	}

	syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	os.Exit(0)
}

func (P *Proc) generateStartCommand(s string) {
	split := strings.Split(s, " ")
	P.commandName = split[0]
	if len(split) > 1 {
		P.args = split[1:]
	} else {
		P.args = []string{""}
	}

	if P.commandName == "/usr/sbin/nginx" {
		P.args = []string{"-g", "daemon off;"}
	}
}
