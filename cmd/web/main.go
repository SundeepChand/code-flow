package main

import (
	"fmt"
	"syscall/js"

	mermaidflowchart "github.com/SundeepChand/code-flow/flowchart/mermaid"
	"github.com/SundeepChand/code-flow/parser/goparser"
	"github.com/pkg/errors"
)

func genWrapper() js.Func {
	goParser := goparser.New()
	flowChart := mermaidflowchart.New()

	generatorFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "invalid number of arguments passed"
		}
		input := args[0].String()

		parsedCode, err := goParser.Parse(input)
		if err != nil {
			return errors.Wrap(err, "error in parsing input code").Error()
		}

		flowChart.Clear()
		flowChart.FromAst(parsedCode)
		flowchartStr, err := flowChart.String()
		if err != nil {
			return errors.Wrap(err, "error in converting flowchart to string").Error()
		}
		return flowchartStr
	})
	return generatorFunc
}

// Export method to javascript, so that they can be called from within JS.
func main() {
	fmt.Println("Init go webassembly")
	js.Global().Set("generateMermaidCode", genWrapper())
	<-make(chan bool)
}
