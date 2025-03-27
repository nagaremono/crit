package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
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
			checkType(command)

		default:
			fmt.Println(command[0] + ": command not found")
		}
	}
}

func exit(command []string) {
	if len(command) >= 2 == false {
		return
	}

	exitCode, err := strconv.Atoi(command[1])
	if err != nil {
		fmt.Println(os.Stderr, "Invalid codee", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func echo(command []string) {
	if !(len(command) >= 2) {
		fmt.Println("")
	}

	toPrint := command[1:]

	fmt.Println(strings.Join(toPrint, " "))
}

func checkType(command []string) {
	if !(len(command) >= 2) {
		fmt.Println(": not found")
	}
	toCheck := command[1]

	if slices.Contains(builtIns, toCheck) {
		fmt.Println(toCheck + " is a shell builtin")
		return
	}

	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		location := filepath.Join(dir, toCheck)
		_, err := os.Stat(location)
		if err == nil {
			fmt.Fprintf(os.Stdout, "%s is %s\n", toCheck, location)
			return
		}
	}

	fmt.Fprintf(os.Stdout, "%s: not found\n", toCheck)
}
