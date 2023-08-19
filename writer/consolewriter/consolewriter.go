package consolewriter

import "fmt"

type ConsoleWriter struct {
}

func New() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (c *ConsoleWriter) Write(result string) error {
	fmt.Println(result)
	return nil
}
