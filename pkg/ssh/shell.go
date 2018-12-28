package ssh

import (
	"os"
	"os/exec"
	"syscall"
)

//ExecCmd executes directly via shell command
func ExecCmd(user string, port string, ipAddress string, command string) error {

	sshCommand := exec.Command("ssh", "-oStrictHostKeyChecking=no", "-l", user, "-p", port, ipAddress, command)
	sshCommand.Stdin = os.Stdin
	sshCommand.Stdout = os.Stdout
	sshCommand.Stderr = os.Stderr

	if err := sshCommand.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		} else {
			return err
		}
	}

	return nil
}

//ExecCmdLocal executes directly via shell command
func ExecCmdLocal(cmd string, args ...string) error {

	sshCommand := exec.Command(cmd, args...)
	sshCommand.Stdin = os.Stdin
	sshCommand.Stdout = os.Stdout
	sshCommand.Stderr = os.Stderr

	if err := sshCommand.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			os.Exit(waitStatus.ExitStatus())
		} else {
			return err
		}
	}

	return nil
}