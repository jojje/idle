package process

import (
	"log"
	"unsafe"

	"golang.org/x/sys/windows"
)

const ProcessAllAccess = windows.STANDARD_RIGHTS_REQUIRED | windows.SYNCHRONIZE | 0xffff

func setIdleProcessPriority(pid uint32) error {
	handle, err := windows.OpenProcess(ProcessAllAccess, false, uint32(pid))
	if err != nil {
		return err
	}
	defer windows.CloseHandle(handle) // ignore failure since there's nothing that can be done about it

	return windows.SetPriorityClass(handle, windows.IDLE_PRIORITY_CLASS)
}

func hasIdlePriority(pid uint32) bool {
	handle, err := windows.OpenProcess(ProcessAllAccess, false, uint32(pid))
	if err != nil {
		if isAlive(pid) {
			log.Printf("[WARNING] failed to get a handle to process with pid=%d", pid)
		}
		return false
	}
	defer windows.CloseHandle(handle) // ignore failure since there's nothing that can be done about it

	prio, err := windows.GetPriorityClass(handle)
	if err != nil {
		if isAlive(pid) {
			log.Printf("[WARNING] failed to get process priority for pid=%d", pid)
		}
		return false
	}

	return prio == windows.IDLE_PRIORITY_CLASS
}

func listProcesses() ([]*Process, error) {
	h, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(h) // ignore error

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))

	// get first process
	err = windows.Process32First(h, &pe)
	if err != nil {
		return nil, err
	}

	var out []*Process

	for {
		// process current entry
		out = append(out, &Process{
			Pid:  pe.ProcessID,
			Name: windows.UTF16ToString(pe.ExeFile[:]),
		})

		err = windows.Process32Next(h, &pe) // iterate to next process
		if err != nil {
			if err == windows.ERROR_NO_MORE_FILES {
				break // done
			}
			return nil, err
		}
	}
	return out, nil
}
