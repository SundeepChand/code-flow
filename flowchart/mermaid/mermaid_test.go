package mermaid_test

import (
	"testing"

	"github.com/SundeepChand/code-flow/flowchart/mermaid"
	"github.com/SundeepChand/code-flow/parser/goparser"
	"github.com/stretchr/testify/require"
)

func TestMermaidGenerator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		inputProgram   string
		expectedOutput string
	}{
		{
			name: "Generate flow chart for fibonacci series program",
			inputProgram: `package main

			import "fmt"
			
			func NthFibonacciNumber(n int) {
				f, s := 1, 1
				if n > 0 {
					fmt.Println(f)
				} else if n == 0 {
					fmt.Println("no series")
				} else {
					fmt.Println("Invalid")
				}
				if n > 1 {
					fmt.Println(s)
				}
			
				for i := 0; i < (n - 2); i++ {
					t := f + s
					fmt.Println(t)
					f, s = s, t
				}
			}`,
			expectedOutput: `---
Title: NthFibonacciNumber
---
flowchart TD


	start(["
		START
	"])

	stop(["
		STOP
	"])

	id:1["
		f, s := 1, 1
	"]

	id:2{"
		n > 0
	"}

	id:3["
		fmt.Println(f)
	"]

	id:4{"
		n == 0
	"}

	id:5["
		fmt.Println(no series)
	"]

	id:6["
		fmt.Println(Invalid)
	"]

	id:7{"
		n > 1
	"}

	id:8["
		fmt.Println(s)
	"]

	id:9["
		i := 0
	"]

	id:10{"
		i < (n - 2)
	"}

	id:11["
		t := f + s
	"]

	id:12["
		fmt.Println(t)
	"]

	id:13["
		f, s = s, t
	"]

	id:14["
		i++
	"]

	id:2 --Yes--> id:3
	
	id:4 --Yes--> id:5
	
	id:4 --No--> id:6
	
	id:2 --No--> id:4
	
	id:1 --> id:2
	
	id:7 --Yes--> id:8
	
	id:3 --> id:7
	
	id:5 --> id:7
	
	id:6 --> id:7
	
	id:9 --> id:10
	
	id:11 --> id:12
	
	id:12 --> id:13
	
	id:10 --Yes--> id:11
	
	id:14 --> id:10
	
	id:13 --> id:14
	
	id:8 --> id:9
	
	id:7 --> id:9
	
	start --> id:1
	
	id:10 --> stop
	`,
		},
	}

	for _, tc := range tests {

		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			goParser := goparser.New()

			parsedAst, err := goParser.Parse(tc.inputProgram)
			require.NoError(t, err)

			m := mermaid.New()
			m.FromAst(parsedAst)
			got, err := m.String()
			if tc.expectedOutput != "" {
				require.Equal(t, tc.expectedOutput, got)
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}

}
