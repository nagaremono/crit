package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

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

func typeOf(command []string) {
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

func run(command []string) {
	_, err := exec.LookPath(command[0])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s: not found\n", command[0])
		return
	}

	cmd := exec.Command(command[0], command[1:]...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	fmt.Fprintf(os.Stdout, "%s", out)
}
