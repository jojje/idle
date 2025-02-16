package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jojje/idle/pattern"
	"github.com/jojje/idle/process"
	"github.com/jojje/idle/util"
)

var (
	watch           bool
	caseInsensitive bool
	pollInterval    int64
	excludedPids    string
	showVersion     bool
	version         = "dev"
)

func usage() {
	pgm := util.BaseName(os.Args[0])
	fmt.Fprintf(os.Stderr, "usage: %v process [process ...]\n\n%s", pgm,
		"Changes priority to Idle for matching processes.\n\n"+

			"Process is treated as a case-sensitive substring to match running processes against. "+
			"When prefixed and suffixed (surrounded by) '/' characters, the pattern is interpreted "+
			"as a regular expression.\n\n"+

			"Options:\n")

	flag.PrintDefaults()
	os.Exit(2)
}

func parseArgs() []string {
	flag.BoolVar(&watch, "w", false, "Watch processes and lower the priority of any matches the process expression(s)")
	flag.BoolVar(&caseInsensitive, "i", false, "Use case insensitive process name matching")
	flag.StringVar(&excludedPids, "e", "", "Exclude specific pids from being idled. Argument is specified as a csv string: pid[,...]")
	flag.Int64Var(&pollInterval, "p", 500, "Process poll interval (in milliseconds)")
	flag.BoolVar(&showVersion, "V", false, "Show version of this program")
	flag.Usage = usage

	if len(os.Args) < 2 {
		usage()
	}
	flag.Parse()
	return flag.Args()
}

func main() {
	os.Args = append(os.Args, "-V")
	args := parseArgs()

	if showVersion {
		fmt.Printf("Version: %s", version)
		return
	}

	var matchers []pattern.Matcher
	for _, expr := range args {
		m, err := pattern.NewMatcher(expr, caseInsensitive)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid process expression: %s", expr)
			os.Exit(3)
		}
		matchers = append(matchers, m)
	}

	excluded := util.ParseExcludedPids(excludedPids)

	lowerPrio := func() {
		procs := process.FindProcesses(matchers)
		procs = util.ExcludeProcesses(procs, excluded)
		process.LowerPriorities(procs)
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	lowerPrio()

	if watch {
		startWatching(lowerPrio, time.Duration(pollInterval)*time.Millisecond)
	}
}

func startWatching(lowerPrio func(), pollInterval time.Duration) {
	for {
		time.Sleep(pollInterval)
		lowerPrio()
	}
}
