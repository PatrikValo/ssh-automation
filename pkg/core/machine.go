package core

import (
	"net"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/PatrikValo/ssh-automation/pkg/printer"
	"github.com/PatrikValo/ssh-automation/pkg/program"
)

type machine struct {
	config     *program.Config
	connection *ssh.Client
	host       string
	terminal   chan<- printer.Printer
}

func (machine *machine) selectTimeout(connChan <-chan *ssh.Client, errorChan <-chan error) bool {
	if machine.config.ConnectionTimeout <= 0 {
		select {
		case conn := <-connChan:
			machine.connection = conn
			return true
		case err := <-errorChan:
			machine.terminal <- printer.ErrorPrinter{Host: machine.host, Msg: err.Error()}
			return false
		}
	} else {
		select {
		case conn := <-connChan:
			machine.connection = conn
			machine.terminal <- printer.Green(machine.host, "Connection is established.\n")
			return true
		case err := <-errorChan:
			machine.terminal <- printer.ErrorPrinter{Host: machine.host, Msg: err.Error()}
			return false
		case <-time.After(time.Duration(machine.config.ConnectionTimeout) * time.Millisecond):
			machine.terminal <- printer.ErrorPrinter{
				Host: machine.host,
				Msg:  "TIMEOUT: Connection to machine takes too long.",
			}
			return false
		}
	}
}

func (machine *machine) connect(user string, authMethod *ssh.AuthMethod, connected chan<- bool) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			*authMethod,
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	connChan := make(chan *ssh.Client)
	errorChan := make(chan error)

	go func() {
		conn, err := ssh.Dial("tcp", machine.host, sshConfig)
		if err != nil {
			errorChan <- err
		} else {
			connChan <- conn
		}
	}()

	connected <- machine.selectTimeout(connChan, errorChan)
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
