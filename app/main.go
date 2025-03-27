package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var builtIns = []string{"echo", "exit", "type"}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(os.Stderr, "Error reading command:", err)
			os.Exit(1)
		}

		input = strings.Trim(input, "\n")
		command := strings.Split(input, " ")

		switch command[0] {
		case "exit":
			exit(command)
		case "echo":
			echo(command)
		case "type":
			typeOf(command)

		default:
			run(command)
		}
	}
}
