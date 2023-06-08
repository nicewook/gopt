package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		fmt.Println("You entered: ", text)
	}
}
