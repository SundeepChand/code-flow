package generator

import (
	"github.com/SundeepChand/code-flow/flowchart"
	"github.com/SundeepChand/code-flow/parser"
	"github.com/SundeepChand/code-flow/reader"
	"github.com/SundeepChand/code-flow/writer"
	"github.com/pkg/errors"
)

type Generator struct {
	readerImpl    reader.Reader
	parserImpl    parser.Parser
	flowChartImpl flowchart.FlowChart
	writerImpl    writer.Writer
}

func New(r reader.Reader, p parser.Parser, f flowchart.FlowChart, w writer.Writer) *Generator {
	return &Generator{
		readerImpl:    r,
		parserImpl:    p,
		flowChartImpl: f,
		writerImpl:    w,
	}
}

func (g *Generator) Generate() error {
	input, err := g.readerImpl.Read()
	if err != nil {
		return errors.Wrap(err, "could not load input")
	}
	parsedCode, err := g.parserImpl.Parse(input)
	if err != nil {
		return errors.Wrap(err, "error in parsing input code")
	}
	g.flowChartImpl.FromAst(parsedCode)
	flowchartStr, err := g.flowChartImpl.String()
	if err != nil {
		return errors.Wrap(err, "error in converting flowchart to string")
	}
	err = g.writerImpl.Write(flowchartStr)
	if err != nil {
		return errors.Wrap(err, "error in writing generated flowchart into file")
	}
	return nil
}
