package util

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jojje/idle/process"
)

func BaseName(filename string) string {
	base := filepath.Base(os.Args[0])
	return strings.TrimSuffix(base, filepath.Ext(base))
}

func Includes(items []uint32, item uint32) bool {
	for _, candidate := range items {
		if item == candidate {
			return true
		}
	}
	return false
}

func ParseExcludedPids(csv string) []uint32 {
	var out []uint32
	for _, pidStr := range strings.Split(csv, ",") {
		pid, err := strconv.ParseUint(pidStr, 10, 32)
		if err != nil {
			continue
		}
		out = append(out, uint32(pid))
	}
	return out
}

func ExcludeProcesses(procs []*process.Process, pidsToExclude []uint32) []*process.Process {
	var out []*process.Process

	for _, p := range procs {
		skip := false
		for _, pid := range pidsToExclude {
			if pid == p.Pid {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		out = append(out, p)
	}

	return out
}
