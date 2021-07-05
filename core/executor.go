package core

import (
	"errors"
	"github.com/PatrikValo/ssh-automation/printer"
	"github.com/PatrikValo/ssh-automation/program"
	"golang.org/x/crypto/ssh"
	"strconv"
)

type Executor struct {
	machines []*machine
	user     string
	auth     *ssh.AuthMethod
	terminal chan<- printer.Printer
}

func (executor *Executor) init() error {
	executor.printHeader("Try to connect machines")

	connected := make(chan bool)
	for _, mach := range executor.machines {
		go mach.connect(executor.user, executor.auth, connected)
	}

	success, fail := executor.getResult(connected)
	executor.printResult(success, fail)

	if fail != 0 {
		return errors.New("during initialization error was occurred")
	}
	return nil
}

func (executor *Executor) execTask(task *program.Task) error {
	executor.printHeader(task.Name)

	executed := make(chan bool)
	for _, mach := range executor.machines {
		go mach.execCmd(task.Cmd, task.Out, executed)
	}

	success, fail := executor.getResult(executed)
	executor.printResult(success, fail)

	if fail != 0 {
		return errors.New("during execution of the task error was occurred")
	}
	return nil
}

func (executor *Executor) getResult(okChan <-chan bool) (int, int) {
	success := 0

	for range executor.machines {
		if ok := <-okChan; ok {
			success++
		}
	}

	return success, len(executor.machines) - success
}

func (executor *Executor) printHeader(name string) {
	executor.terminal <- printer.ColorPrinter{Msg: "TASK [" + name + "]"}
	executor.terminal <- printer.ColorPrinter{Msg: "*******************************************"}
}

func (executor *Executor) printResult(success int, fail int) {
	if fail == 0 {
		executor.terminal <- printer.Green("", "|-+ SUCCESS: "+strconv.Itoa(success))
		executor.terminal <- printer.Green("", "|-+ FAIL: "+strconv.Itoa(fail))
	} else {
		executor.terminal <- printer.Red("", "|-+ SUCCESS: "+strconv.Itoa(success))
		executor.terminal <- printer.Red("", "|-+ FAIL: "+strconv.Itoa(fail))
	}
	executor.terminal <- printer.ColorPrinter{Msg: ""}
}

func (executor *Executor) Execute(program *program.Program) error {
	err := executor.init()
	if err != nil {
		return err
	}

	for _, task := range program.Tasks {
		err = executor.execTask(&task)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExecutorFactory(hosts []string, user string, auth *ssh.AuthMethod, printChan chan<- printer.Printer) Executor {
	workers := make([]*machine, 0, len(hosts))

	for _, host := range hosts {
		workers = append(workers, &machine{connection: nil, host: host, terminal: printChan})
	}

	return Executor{machines: workers, user: user, auth: auth, terminal: printChan}
}
