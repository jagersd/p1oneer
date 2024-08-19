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

func (P *Proc) StartLong(command string, args []string, signalChannel chan (os.Signal)) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start %s with Err: %v", P.commandName, err)
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Printf(P.commandName+" exited with error: %v", err)
			os.Exit(err.(*exec.ExitError).ExitCode())
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

func (P *Proc) StartOne(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("command %s failed with error: %v", command, err)
	}
}
