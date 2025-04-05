package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	builtIns      = []string{"echo", "exit", "type", "pwd"}
	redirectOpPtn = "^(1|2)?>$"
	redirectOpReg = regexp.MustCompile(redirectOpPtn)
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
			os.Exit(1)
		}

		parsed := parseCmdArgs(strings.TrimSpace(input))
		redirectOpIndex := slices.IndexFunc(parsed, func(s string) bool {
			matched, _ := regexp.MatchString(redirectOpPtn, s)
			return matched
		})

		var output string
		commandName := parsed[0]
		var commandArgs []string

		if redirectOpIndex != -1 {
			commandArgs = parsed[1:redirectOpIndex]
		} else {
			commandArgs = parsed[1:]
		}

		err, output = execCommand(commandName, commandArgs)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}

		if redirectOpIndex != -1 {
			writeToFile(parsed[redirectOpIndex+1], output)
		} else {
			fmt.Fprint(os.Stdout, output)
		}

	}
}

func execCommand(name string, args []string) (error, string) {
	var err error
	var output string

	switch name {
	case "exit":
		exit(args)
	case "echo":
		output, err = echo(args)
	case "type":
		output, err = typeOf(args)
	case "pwd":
		output, err = pwd()
	case "cd":
		err = cd(args)

	default:
		output, err = run(name, args)
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

func echo(commandArgs []string) (string, error) {
	output := fmt.Sprintf("%s\n", strings.Join(commandArgs, " "))
	return output, nil
}

func typeOf(commandArgs []string) (string, error) {
	if len(commandArgs) == 0 {
		return fmt.Sprintln(": not found"), nil
	}
	toCheck := commandArgs[0]

	if slices.Contains(builtIns, toCheck) {
		return fmt.Sprintln(toCheck + " is a shell builtin"), nil
	}

	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		location := filepath.Join(dir, toCheck)
		_, err := os.Stat(location)
		if err == nil {
			return fmt.Sprintf("%s is %s\n", toCheck, location), nil
		}
	}

	return fmt.Sprintf("%s: not found\n", toCheck), nil
}

func run(commandName string, commandArgs []string) (string, error) {
	_, err := exec.LookPath(commandName)
	if err != nil {
		output := fmt.Sprintf("%s: command not found\n", commandName)
		return output, nil
	}

	cmd := exec.Command(commandName, commandArgs...)
	out, err := cmd.Output()
	if err != nil {
		switch e := err.(type) {
		case *exec.ExitError:
			return string(out), errors.New(string(e.Stderr))
		default:
			return string(out), err
		}
	}
	output := fmt.Sprintf("%s", out)
	return output, nil
}

func pwd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		output := fmt.Sprintf("%s\n", err.Error())
		return output, err
	}
	output := fmt.Sprintf("%s\n", wd)
	return output, err
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
		return errors.New(msg)
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

func writeToFile(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil
}
