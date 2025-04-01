package main

import (
	"bufio"
	"errors"
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

		parsed := parseCmdArgs(strings.TrimSpace(input))

		commandName := parsed[0]
		commandArgs := parsed[1:]

		err, output := execCommand(commandName, commandArgs)

		if err != nil {
			fmt.Println(os.Stdout, err)
		}
		fmt.Print(os.Stdout, output)
	}
}

func execCommand(name string, args []string) (error, string) {
	var err error
	var output string

	switch name {
	case "exit":
		exit(args)
	case "echo":
		err, output = echo(args)
	case "type":
		err, output = typeOf(args)
	case "pwd":
		err, output = pwd()
	case "cd":
		cd(args)

	default:
		run(name, args)
	}

	return err, output
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

func echo(commandArgs []string) (error, string) {
	output := fmt.Sprintf("%s\n", strings.Join(commandArgs, " "))
	return nil, output
}

func typeOf(commandArgs []string) (error, string) {
	if len(commandArgs) == 0 {
		return nil, fmt.Sprintln(": not found")
	}
	toCheck := commandArgs[0]

	if slices.Contains(builtIns, toCheck) {
		return nil, fmt.Sprintln(toCheck + " is a shell builtin")
	}

	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		location := filepath.Join(dir, toCheck)
		_, err := os.Stat(location)
		if err == nil {
			return nil, fmt.Sprintf("%s is %s\n", toCheck, location)
		}
	}

	return nil, fmt.Sprintf("%s: not found\n", toCheck)
}

func run(commandName string, commandArgs []string) (error, string) {
	_, err := exec.LookPath(commandName)
	if err != nil {
		output := fmt.Sprintf("%s: not found\n", commandName)
		return err, output
	}

	cmd := exec.Command(commandName, commandArgs...)
	out, err := cmd.Output()
	if err != nil {
		output := fmt.Sprintf("%s\n", err)
		return err, output
	}
	output := fmt.Sprintf("%s\n", out)
	return nil, output
}

func pwd() (error, string) {
	wd, err := os.Getwd()
	if err != nil {
		output := fmt.Sprintf("%s\n", err.Error())
		return err, output
	}
	output := fmt.Sprintf("%s\n", wd)
	return nil, output
}

func cd(commandArgs []string) error {
	var path string

	if len(commandArgs) == 0 || commandArgs[0] == "~" {
		path = os.Getenv("HOME")
	} else {
		path = commandArgs[0]
	}
	err := os.Chdir(path)
	if err != nil {
		msg := fmt.Sprintf("cd: %s: No such file or directory\n", path)
		errors.New(msg)
	}
	return nil
}

var doubleQuoteExc = []rune{
	'\\', '$', '"', '\n',
}

func parseCmdArgs(args string) []string {
	args = args + " "
	var commandArgs []string
	var tmp string
	var inSingleQuote bool
	var inDoubleQuote bool

	for index := 0; index < len(args); index++ {
		char := rune(args[index])
		if char == '\\' {
			nextChar := rune(args[index+1])
			if (inSingleQuote && index+1 != len(args)) ||
				(inDoubleQuote && !slices.Contains(doubleQuoteExc, nextChar)) {
				tmp = tmp + string('\\') + string(nextChar)
			} else {
				tmp = tmp + string(nextChar)
			}
			index++
		} else if char == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
		} else if char == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
		} else if char == ' ' {
			if !inSingleQuote && !inDoubleQuote {
				if len(tmp) > 0 {
					commandArgs = append(commandArgs, tmp)
					tmp = ""
				}
			} else {
				tmp = tmp + string(char)
			}
		} else {
			tmp = tmp + string(char)
		}
	}

	return commandArgs
}
