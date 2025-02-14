package util

import (
	"reflect"
	"testing"

	"github.com/jojje/idle/process"
)

func Test_parseExcludedPids(t *testing.T) {
	tests := []struct {
		name string
		csv  string
		want []uint32
	}{
		{name: "single pid", csv: "123", want: []uint32{123}},
		{name: "two pids", csv: "11,2", want: []uint32{11, 2}},
		{name: "two pids plus invalid item", csv: "11,foo,2", want: []uint32{11, 2}},
		{name: "invalid item", csv: "11 2", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseExcludedPids(tt.csv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseExcludedPids() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExcludeProcesses(t *testing.T) {
	type args struct {
		procs         []*process.Process
		pidsToExclude []uint32
	}

	procs := []*process.Process{
		{Pid: 1, Name: "a"},
		{Pid: 2, Name: "b"},
		{Pid: 3, Name: "c"},
	}

	tests := []struct {
		name string
		args args
		want []*process.Process
	}{
		{"no excludes", args{procs, []uint32{}}, procs},
		{"exclude one", args{procs, []uint32{2}}, []*process.Process{procs[0], procs[2]}},
		{"exclude two", args{procs, []uint32{1, 2}}, []*process.Process{procs[2]}},
		{"exclude all", args{procs, []uint32{1, 2, 3}}, nil},
		{"exclude other", args{procs, []uint32{4}}, procs},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExcludeProcesses(tt.args.procs, tt.args.pidsToExclude); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExcludeProcesses() = %v, want %v", got, tt.want)
			}
		})
	}
}
