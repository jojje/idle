package process

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jojje/idle/pattern"
)

func TestProcessListing(t *testing.T) {
	procs, err := listProcesses()
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, p := range procs {
		if p.Name[:2] == "go" && p.Pid > 0 {
			found = true
		}
	}

	if !found {
		t.Error("Failed to find own go process")
	}
}

func TestLowerPriorities(t *testing.T) {
	pgm := filepath.Base(os.Args[0])
	match, err := pattern.NewMatcher(pgm, false)
	if err != nil {
		t.Fatal(err)
	}
	procs := FindProcesses([]pattern.Matcher{match})

	if len(procs) == 0 {
		t.Fatal("failed to find current go process")
	}

	LowerPriorities(procs)

	allLowered := true
	for _, p := range procs {
		if !hasIdlePriority(p.Pid) {
			allLowered = false
		}
	}

	if !allLowered {
		t.Errorf("Failed to lower priorities for go process %v", procs)
	}
}
