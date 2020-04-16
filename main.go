package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

const usage = "Usage: runwith <env-file> <arguments>"

func main() {
	args := os.Args[1:]

	// Make sure we have enough args
	if len(args) < 2 {
		fmt.Println("Not enough arguments.")
		fmt.Println(usage)
		return
	}

	envFile := args[0]
	cmd := args[1:]

	// Check that env file exists
	if !fileExists(envFile) {
		fmt.Println("Env file is not valid.")
		fmt.Println(usage)
		return
	}

	envVars, err := parseEnvFile(envFile)
	if err != nil {
		fmt.Printf("Error parsing env file: %v\n", err)
		return
	}

	resp, err := runCmd(cmd, envVars)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(resp))

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func parseEnvFile(filename string) ([]string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read the envfile")
	}
	return strings.Split(strings.TrimSuffix(string(file), "\n"), "\n"), nil
}

func runCmd(cmd []string, envVars []string) ([]byte, error) {
	var args []string

	fmt.Println("command y'all")
	fmt.Println(cmd)
	// Check if the cmd has args
	if len(cmd) > 1 {
		args = cmd[1:]
		fmt.Println("We got args y'all")
		fmt.Println(args)
	} else {
		args = []string{}
	}

	// Prep the command
	command := exec.Command(cmd[0], args...)
	command.Env = os.Environ()

	for _, envVar := range envVars {
		command.Env = append(command.Env, envVar)
	}

	resp, err := command.Output()
	if err != nil {
		return resp, err
	}

	return resp, nil
}
