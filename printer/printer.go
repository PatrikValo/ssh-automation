package printer

import "fmt"

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
)

type Printer interface {
	Print()
}

type SuccessPrinter struct {
	Msg string
}

func (printer SuccessPrinter) Print() {
	fmt.Println(Green + printer.Msg + Reset)
}

type ClassicPrinter struct {
	Msg string
}

func (printer ClassicPrinter) Print() {
	fmt.Println(printer.Msg)
}

type WarningPrinter struct {
	Msg string
}

func (printer WarningPrinter) Print() {
	fmt.Println(Yellow + printer.Msg + Reset)
}

type ErrorPrinter struct {
	Err error
}

func (printer ErrorPrinter) Print() {
	fmt.Println(Red + printer.Err.Error() + Reset)
}

func PrintOnTerminal(printChan <-chan Printer) {
	for printer := range printChan {
		printer.Print()
	}
}
