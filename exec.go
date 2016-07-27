package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type execution struct {
	command *exec.Cmd
}

func runEnclosedCommand(command string, arguments []string) {
	execution := newExecution(command, arguments)
	signalChannel := execution.createSignalChannel()

	go execution.listenToSignals(signalChannel)
	go execution.runChecked()
	time.Sleep(1 * time.Second)
	sendSignal(execution.command.Process, syscall.SIGKILL)
}

func (instance *execution) listenToSignals(signalChannel chan os.Signal) {
	for {
		osSignal, channelReady := <-signalChannel
		if channelReady {
			process := instance.command.Process
			if process != nil {
				signal := osSignal.(syscall.Signal)
				log.Printf("Forward signal %v to %v (#%v)...", signal, instance.command.Args[0], process.Pid)
				sendSignal(instance.command.Process, signal)
			}
		} else {
			break
		}
	}
}

func (instance *execution) runChecked() {
	exitCode, err := instance.run()
	if err != nil {
		log.Fatalf("Cannot execute %v. Cause: %v", instance.commandLine(), err)
	}
	log.Printf("%v exited with %v", instance.commandLine(), exitCode)
	os.Exit(exitCode)
}

func newExecution(command string, arguments []string) *execution {
	cmd := exec.Command(command, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = createSysProcAttr()

	return &execution{
		command: cmd,
	}
}

func (instance *execution) run() (int, error) {
	log.Printf("Start enclosed process using: %v", instance.commandLine())
	var waitStatus syscall.WaitStatus
	if err := instance.command.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			exitSignal := waitStatus.Signal()
			if exitSignal > 0 {
				return (int(exitSignal) + 128), nil
			}
			return waitStatus.ExitStatus(), nil
		}
		return 0, err
	}
	waitStatus = instance.command.ProcessState.Sys().(syscall.WaitStatus)
	return waitStatus.ExitStatus(), nil
}

func (instance *execution) commandLine() string {
	result := ""
	for i, arg := range instance.command.Args {
		if i != 0 {
			result += " "
		}
		if strings.Contains(arg, "\"") || strings.Contains(arg, "\\") || strings.Contains(arg, " ") {
			result += strconv.Quote(arg)
		} else {
			result += arg
		}
	}
	return result
}
