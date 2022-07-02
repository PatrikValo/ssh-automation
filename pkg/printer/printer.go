package printer

import "fmt"

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
)

func structuredPrint(host string, msg string, color string) {
	if host != "" {
		fmt.Println(color + "HOST: [" + host + "]" + reset)
	}
	fmt.Println(color + msg + reset)
}

type Printer interface {
	Print()
}

type ColorPrinter struct {
	Host  string
	Msg   string
	Color string
}

func (printer ColorPrinter) Print() {
	structuredPrint(printer.Host, printer.Msg, printer.Color)
}

func Red(host, msg string) ColorPrinter {
	return ColorPrinter{Host: host, Msg: msg, Color: red}
}

func Green(host, msg string) ColorPrinter {
	return ColorPrinter{Host: host, Msg: msg, Color: green}
}

func Yellow(host, msg string) ColorPrinter {
	return ColorPrinter{Host: host, Msg: msg, Color: yellow}
}

type ErrorPrinter struct {
	Host string
	Msg  string
}

func (printer ErrorPrinter) Print() {
	structuredPrint(printer.Host, printer.Msg, red)
	fmt.Println("")
}

func RegisterPrinter(printChan <-chan Printer) {
	for printer := range printChan {
		printer.Print()
	}
}
