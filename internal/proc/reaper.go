package proc

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"syscall"
)

type reaper struct{}

const statfmt = "%s %d"

func (r *reaper) scan() error {
	managedPids := r.getMonitoredPids()
	if len(managedPids) == 0 {
		return nil
	}
	files, err := filepath.Glob("/proc/[0-9]*")
	if err != nil {
		return err
	}

	for _, dir := range files {
		stat := filepath.Join(dir, "stat")
		data, err := os.ReadFile(stat)
		if err != nil {
			log.Println(err)
			continue
		}
		info := strings.FieldsFunc(string(data), func(r rune) bool {
			return r == '(' || r == ')'
		})
		if len(info) != 3 {
			log.Println("Invalid pid info for ", data)
			continue
		}
		var state string
		var ppid int
		if _, err := fmt.Sscanf(info[2], statfmt, &state, &ppid); err != nil {
			log.Println("fmt err")
			return err
		}
		if state == "Z" && slices.Contains(managedPids, ppid) {
			if p, err := strconv.Atoi(strings.Trim(info[0], " ")); err != nil {
				log.Println("string to pid conversion err", err)
				continue
			} else {
				r.kill(p)
			}
		}
	}

	return nil
}

func (r *reaper) kill(pid int) {
	fmt.Println("Reaping:", pid)
	p, err := os.FindProcess(pid)
	if err != nil {
		log.Println("os err pid not found", err)
	}
	if err := p.Kill(); err != nil {
		log.Println("pid could not be killed", err)
	}
	if err := p.Signal(os.Kill); err != nil {
		log.Println("Unable to send signal", err)
	}
	err = syscall.Kill(pid, syscall.SIGKILL)
	if err != nil {
		log.Printf("Unable to kill pid %d ", err)
	}
}

func (r *reaper) getMonitoredPids() []int {
	var pids []int
	for _, p := range hub.processes {
		err := syscall.Kill(p.Pid, syscall.SIGCHLD)
		if err != nil {
			fmt.Printf("Error sending SIGCHLD to process %d: %v\n", p.Pid, err)
		} else {
			fmt.Printf("Successfully sent SIGCHLD to process %d\n", p.Pid)
		}
		pids = append(pids, p.Pid)
	}
	return pids
}
