package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("db > ")
		scanner.Scan()
		input := scanner.Text()

		if input == ".exit" {
			fmt.Println("adios!")
			break
		} else {
			fmt.Println("TODO: Parse the following command: ", input)
		}
	}
}
