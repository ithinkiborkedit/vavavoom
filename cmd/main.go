package main

import (
	"fmt"
	"nelly/internal/dsl"
	"nelly/internal/dslcore"
	"os"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("usage program script")
		os.Exit(1)
	}

	filename := os.Args[1]
	src, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read file: %v\n", err)
		os.Exit(1)
	}
	program, err := dsl.ParseScript(string(src))
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	state := &dslcore.ExecutionState{
		Vars: make(map[string]interface{}),
	}

	dsl.RegisterDefaultCommands()

	if err := program.Execute(state); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)

	}

}
