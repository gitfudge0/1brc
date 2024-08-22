package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/gitfudge0/1brc.git/internal/creator"
)

func main() {
	fmt.Println("What do you want to do?")
	fmt.Println("1. Create data")
	fmt.Println("2. Solve problem")
	fmt.Println("Enter choice: ")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)

	switch input, _ := reader.ReadString('\n'); input {
	case "1":
		startTime := time.Now()

		fmt.Println("Hold on while the file is being generated...")
		creator.CreateData()

		fmt.Println("done")
		fmt.Printf("Time taken: %s\n", time.Since(startTime))

	case "2":
		fmt.Println("Will solve")
		os.Exit(0)
	default:
		fmt.Println("lol")
		os.Exit(1)
	}
}
