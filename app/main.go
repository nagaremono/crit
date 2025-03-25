package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

		if command[0] != "exit" {
			fmt.Println(command[0] + ": command not found")
		}

		if len(command) >= 2 == false {
			continue
		}

		exitCode, err := strconv.Atoi(command[1])
		if err != nil {
			fmt.Println(os.Stderr, "Invalid codee", err)
			os.Exit(1)
		}

		os.Exit(exitCode)
	}
}
