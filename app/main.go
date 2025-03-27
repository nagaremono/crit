package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

func exit(commandArgs []string) {
	if len(commandArgs) >= 1 == false {
		return
	}

	exitCode, err := strconv.Atoi(commandArgs[1])
	if err != nil {
		fmt.Println(os.Stderr, "Invalid codee", err)
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func echo(commandArgs []string) {
	if !(len(commandArgs) >= 2) {
		fmt.Println("")
	}

	toPrint := commandArgs[1:]

	fmt.Println(strings.Join(toPrint, " "))
}

func typeOf(commandArgs []string) {
	if !(len(commandArgs) >= 2) {
		fmt.Println(": not found")
	}
	toCheck := commandArgs[1]

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

func run(commandName string, commandArgs []string) {
	_, err := exec.LookPath(commandName)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s: not found\n", commandName)
		return
	}

	cmd := exec.Command(commandName, commandArgs...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	fmt.Fprintf(os.Stdout, "%s", out)
}
