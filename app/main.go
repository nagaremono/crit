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

var builtIns = []string{"echo", "exit", "type", "pwd"}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(os.Stderr, "Error reading command:", err)
			os.Exit(1)
		}

		command, args, found := strings.Cut(strings.TrimSpace(input), " ")
		commandName := command
		var commandArgs []string
		if found {
			commandArgs = parseCmdArgs(args)
		}

		switch commandName {
		case "exit":
			exit(commandArgs)
		case "echo":
			echo(commandArgs)
		case "type":
			typeOf(commandArgs)
		case "pwd":
			pwd()
		case "cd":
			cd(commandArgs)

		default:
			run(commandName, commandArgs)
		}
	}
}

func exit(commandArgs []string) {
	var exitCode int

	if len(commandArgs) == 0 {
		exitCode = 0
	} else {
		var err error
		exitCode, err = strconv.Atoi(commandArgs[0])
		if err != nil {
			fmt.Println(os.Stderr, "Invalid codee", err)
			os.Exit(1)
		}
	}

	os.Exit(exitCode)
}

func echo(commandArgs []string) {
	fmt.Println(strings.Join(commandArgs, " "))
}

func typeOf(commandArgs []string) {
	if len(commandArgs) == 0 {
		fmt.Println(": not found")
	}
	toCheck := commandArgs[0]

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

func pwd() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	fmt.Fprintf(os.Stdout, "%s\n", wd)
}

func cd(commandArgs []string) {
	var path string

	if len(commandArgs) == 0 || commandArgs[0] == "~" {
		path = os.Getenv("HOME")
	} else {
		path = commandArgs[0]
	}
	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", path)
	}
}

func parseCmdArgs(args string) []string {
	var commandArgs []string
	var tmp string
	var inSingleQuote bool
	var inDoubleQuote bool

	for _, i := range args + " " {
		if i == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if i == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		} else if i == ' ' {
			if !inSingleQuote && !inDoubleQuote {
				if len(tmp) > 0 {
					commandArgs = append(commandArgs, tmp)
					tmp = ""
				}
			} else {
				tmp = tmp + string(i)
			}
		} else {
			tmp = tmp + string(i)
		}
	}

	return commandArgs
}
