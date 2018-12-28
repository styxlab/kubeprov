package ssh

import (
	"os"
	"os/exec"
	"syscall"
)

//ExecCmd executes directly via shell command
func ExecCmd(user string, port string, ipAddress string, cmd string) error {

	sshCommand := exec.Command("ssh -oStrictHostKeyChecking=no", "-l", user, "-p", port, ipAddress, cmd)
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