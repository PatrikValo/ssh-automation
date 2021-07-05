package core

import (
	"github.com/PatrikValo/ssh-automation/printer"
	"golang.org/x/crypto/ssh"
	"net"
)

type machine struct {
	connection *ssh.Client
	host       string
	terminal   chan<- printer.Printer
}

func (machine *machine) connect(user string, authMethod *ssh.AuthMethod, connected chan<- bool) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			*authMethod,
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	conn, err := ssh.Dial("tcp", machine.host, sshConfig)

	if err != nil {
		machine.terminal <- printer.ErrorPrinter{Host: machine.host, Msg: err.Error()}
		connected <- false
		return
	}

	machine.connection = conn
	connected <- true
}

func (machine *machine) execCmd(cmd string, out bool, executed chan<- bool) {
	session, err := machine.connection.NewSession()

	if err != nil {
		executed <- false
		return
	}

	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	output, err := sendRequest(session, cmd)

	if err != nil {
		machine.terminal <- printer.ErrorPrinter{Host: machine.host, Msg: string(output) + err.Error()}
		executed <- false
		return
	}

	if out {
		machine.terminal <- printer.Yellow(machine.host, string(output))
	}

	executed <- true
}

func sendRequest(session *ssh.Session, cmd string) ([]byte, error) {
	modes := ssh.TerminalModes{
		ssh.ECHO: 0,
	}

	err := session.RequestPty("xterm", 24, 80, modes)

	if err != nil {
		return make([]byte, 0), err
	}

	output, err := session.Output(cmd)

	if err != nil {
		return output, err
	}

	return output, nil
}