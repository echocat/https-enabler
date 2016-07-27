// +build windows

package main

import (
	"os"
	"os/signal"
	"syscall"
	"log"
)

var specialSignalHandling = false

func (instance *execution) createSignalChannel() chan os.Signal {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGILL,
		syscall.SIGTRAP,
		syscall.SIGABRT,
		syscall.SIGBUS,
		syscall.SIGFPE,
		syscall.SIGKILL,
		syscall.SIGSEGV,
		syscall.SIGPIPE,
		syscall.SIGALRM,
		syscall.SIGTERM,
	)
	return signalChannel
}

func sendSignal(process *os.Process, what syscall.Signal) {
	if process == nil {
		return
	}
	switch what {
	case syscall.SIGTERM:
		sendSpecialSignal(process, syscall.CTRL_BREAK_EVENT)
	case syscall.SIGINT:
		sendSpecialSignal(process, syscall.CTRL_C_EVENT)
	default:
		process.Signal(what)
	}
}

func sendSpecialSignal(process *os.Process, what uintptr) {
	if specialSignalHandling {
		pid := process.Pid
		d, e := syscall.LoadDLL("kernel32.dll")
		if e != nil {
			log.Fatalf("Could not signal %v to #%v. Cause: %v", what, pid, e)
		}
		p, e := d.FindProc("GenerateConsoleCtrlEvent")
		if e != nil {
			log.Fatalf("Could not signal %v to #%v. Cause: %v", what, pid, e)
		}
		p.Call(what, uintptr(pid))
	} else {
		process.Signal(syscall.SIGKILL)
	}
}

func createSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}
