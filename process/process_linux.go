package process

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func setIdleProcessPriority(pid uint32) error {
	// use a shell utility (like on mac),  difference being that 19 is the lowest priority on linux
	cmd := exec.Command("renice", "19", "-p", strconv.FormatUint(uint64(pid), 10))
	return cmd.Run()
}

func hasIdlePriority(pid uint32) bool {
	path := fmt.Sprintf("/proc/%d/stat", pid)
	content, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	fields := strings.Fields(string(content))
	if len(fields) < 19 {
		return false
	}

	nice, err := strconv.Atoi(fields[18])
	if err != nil {
		return false
	}

	return nice == 19
}

func listProcesses() ([]*Process, error) {
	var processes []*Process

	dirs, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		// consider only number-dirs for pids
		pid, err := strconv.ParseUint(dir.Name(), 10, 32)
		if err != nil {
			continue
		}

		// get the process name from /proc/[pid]/comm
		commPath := filepath.Join("/proc", dir.Name(), "comm")
		comm, err := os.ReadFile(commPath)
		if err != nil {
			continue
		}

		processes = append(processes, &Process{
			Pid:  uint32(pid),
			Name: strings.TrimSpace(string(comm)),
		})
	}

	return processes, nil
}
