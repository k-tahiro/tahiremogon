package handler

import (
	"os/exec"

	myMw "../middleware"
)

func execCommand(mode string, command string, cc *myMw.CustomContext) ([]byte, error) {
	switch mode {
	case "ssh": return execSSHCommand(command, cc)
	default: return exec.Command("sh", "-c", command).Output()
	}
}

func execSSHCommand(command string, cc *myMw.CustomContext) ([]byte, error) {
	session, err := cc.Client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return session.Output(command)
}
