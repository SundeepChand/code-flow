package dummyreader

const simpleProgram = `
package main

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
}
`

type DummyReader struct {
}

func New() *DummyReader {
	return &DummyReader{}
}

func (r *DummyReader) Read() (string, error) {
	return simpleProgram, nil
}
