package process

import (
	"os/exec"
	"strconv"
	"strings"
)

func setIdleProcessPriority(pid uint32) error {
	// use nice utility to set lowest priority (20)
	cmd := exec.Command("renice", "20", "-p", strconv.Itoa(int(pid)))
	return cmd.Run()
}

func hasIdlePriority(pid uint32) bool {
	cmd := exec.Command("ps", "-o", "nice=", "-p", strconv.Itoa(int(pid)))
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	nice, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return false
	}

	return nice == 20
}

func listProcesses() ([]*Process, error) {
	cmd := exec.Command("ps", "-e", "-o", "pid=,comm=")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var out []*Process
	for _, line := range strings.Split(string(output), "\n") {
		if line = strings.TrimSpace(line); line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}

		pid, err := strconv.ParseUint(fields[0], 10, 32)
		if err != nil {
			continue
		}

		out = append(out, &Process{
			Pid:  uint32(pid),
			Name: fields[1],
		})
	}

	return out, nil
}
