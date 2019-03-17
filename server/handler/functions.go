package handler

import (
	"os/exec"

	"golang.org/x/crypto/ssh"

	myMw "../middleware"
)

func execCommand(mode string, command string, cc *myMw.CustomContext) error {
	switch mode {
	case "ssh": err := execSSHCommand(command, cc)
	default: err := exec.Command("sh", "-c", command).Run()
	}
	return err
}

func execSSHCommand(command string, cc *myMw.CustomContext) error {
	session, err := cc.Client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	return session.Run(command)
}
