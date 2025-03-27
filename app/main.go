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

		command := strings.Split(strings.TrimSpace(input), " ")
		commandName := command[0]
		commandArgs := command[1:]

		switch commandName {
		case "exit":
			exit(commandArgs)
		case "echo":
			echo(commandArgs)
		case "type":
			typeOf(commandArgs)

		default:
			run(commandName, commandArgs)
		}
	}
}
