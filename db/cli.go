package main

import (
	"bufio"
	"fmt"
	"os"
)

func validateMetaCommand(cmd string) error {
	switch cmd {
	case ".exit":
		return nil
	}

	return fmt.Errorf("unrecognized meta command: %s", cmd)
}

func doMetaCommand(cmd string) {
	switch cmd {
	case ".exit":
		fmt.Println("adios!")
		os.Exit(0)
	}
}

func validateStatement(cmd string) error {
	switch cmd {
	case "select":
		return nil
	case "insert":
		return nil
	}

	return fmt.Errorf("unrecognized meta command: %s", cmd)
}

func doStatement(cmd string) {
	switch cmd {
	case "select":
		fmt.Println("TODO: select handling goes here!")
	case "insert":
		fmt.Println("TODO: insert handling goes here!")
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("db > ")
		scanner.Scan()
		input := scanner.Text()

		if input[0] == '.' {
			err := validateMetaCommand(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
			doMetaCommand(input)
			continue
		}

		err := validateStatement(input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		doStatement(input)
	}
}
