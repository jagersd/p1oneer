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
	log.Println("Scanning for zombies")
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
				r.kill(ppid, p)
			}
		}
	}

	return nil
}

func (r *reaper) kill(ppid, pid int) {
	fmt.Println("Reaping:", pid)
	p, err := os.FindProcess(ppid)
	p.Signal(syscall.SIGCHLD)
	if err != nil {
		log.Println("Error finding process", err)
		return
	}
	group, err := syscall.Getpgid(ppid)
	if err != nil {
		log.Println("Error getting pgid", err)
	} else {
		log.Println("Group:", -group)
	}
	//if err := syscall.Kill(-group, 15); err != nil {
	//	log.Println("Error sending SIGTERM", err)
	//}
	zombie, err := os.FindProcess(pid)
	if err != nil {
		log.Println("Error finding zombie process", pid, err)
	}
	if _, err := zombie.Wait(); err != nil {
		log.Println("Zombie ", pid, "reaped")
	}
}

func (r *reaper) getMonitoredPids() []int {
	var pids []int
	for _, p := range hub.processes {
		pids = append(pids, p.Pid)
	}
	return pids
}
