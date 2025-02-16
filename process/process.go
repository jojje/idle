package process

import (
	"log"
	"os"

	"github.com/jojje/idle/pattern"
)

type Process struct {
	Pid  uint32
	Name string
}

func LowerPriorities(procs []*Process) {
	for _, p := range procs {
		if hasIdlePriority(p.Pid) {
			continue
		}

		err := setIdleProcessPriority(p.Pid)
		if err != nil { // non-fatal, just inform, not worth dying over
			if !isAlive(p.Pid) {
				continue
			}
			log.Printf("[WARNING] failed to lower process priority for: %s (pid=%d), cause: %s", p.Name, p.Pid, err)
		}

		if !hasIdlePriority(p.Pid) {
			log.Printf("[WARNING] failed to lower process priority for: %s (pid=%d), cause: Unknown", p.Name, p.Pid)
		} else {
			log.Printf("[IDLED] %s (pid=%d)", p.Name, p.Pid)
		}
	}
}

func FindMatchedProcesses(matchers []pattern.Matcher) []*Process {
	var out []*Process

	procs, err := listProcesses()
	if err != nil {
		log.Printf("[WARNING] failed to enumerate processes, cause: %s", err)
		return out
	}

	for _, p := range procs {
		for _, match := range matchers {
			if match(p.Name) {
				out = append(out, p)
			}
		}
	}
	return out
}

func isAlive(pid uint32) bool {
	_, err := os.FindProcess(int(pid))
	return err == nil
}
