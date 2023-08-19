package main

import (
	"log"

	mermaidflowchart "github.com/SundeepChand/code-flow/flowchart/mermaid"
	"github.com/SundeepChand/code-flow/generator"
	"github.com/SundeepChand/code-flow/parser/goparser"
	"github.com/SundeepChand/code-flow/reader/dummyreader"
	"github.com/SundeepChand/code-flow/writer/consolewriter"
)

func main() {
	dummyReader := dummyreader.New()
	goParser := goparser.New()
	flowChart := mermaidflowchart.New()
	consoleWriter := consolewriter.New()

	mermaidFlowChartGenerator := generator.New(dummyReader, goParser, flowChart, consoleWriter)
	err := mermaidFlowChartGenerator.Generate()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully generated flow chart code")
}
