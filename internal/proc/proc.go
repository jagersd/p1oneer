package proc

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Proc struct {
	reportTitle string
	commandName string
	cmd         *exec.Cmd
}

func NewProcess(title string, command string, args []string) *Proc {
	p := Proc{
		commandName: command,
		cmd:         exec.Command(command, args...),
	}
	p.cmd.Stdout = os.Stdout
	p.cmd.Stderr = os.Stderr

	return &p
}

func (P *Proc) StartLong() {
	if err := P.cmd.Start(); err != nil {
		log.Fatalf("Failed to start %s with Err: %v", P.reportTitle, err)
	}

	hub.processChannel <- P.cmd.Process

	go func() {
		if err := P.cmd.Wait(); err != nil {
			log.Printf("%s exited with error:%d %v", P.reportTitle, err.(*exec.ExitError).ExitCode(), err)
		} else {
			log.Println(P.commandName, " exited successfully")
		}
		hub.signalChannel <- syscall.SIGTERM
	}()
}

func (P *Proc) StartOne() {
	if err := P.cmd.Start(); err != nil {
		log.Fatalf("command %s failed with error: %v", P.reportTitle, err)
	}

	go func(process *os.Process) {
		if state, err := process.Wait(); err != nil {
			log.Printf("Single execution process failed with exit code %d error: %v", state.ExitCode(), err)
			hub.signalChannel <- syscall.SIGTERM
		} else {
			fmt.Println("Single execution process done. Reaping pid: ", state.Pid())
		}

	}(P.cmd.Process)
}

func (P *Proc) StartBefore() {
	if err := P.cmd.Run(); err != nil {
		log.Fatalf("command %s failed with error: %v", P.reportTitle, err)
	}
}
