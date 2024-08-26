package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gitfudge0/1brc.git/internal/creator"
	"github.com/gitfudge0/1brc.git/internal/solver"
	"github.com/gitfudge0/1brc.git/internal/utils"
)

func main() {
	utils.ClearScreen()
	fmt.Println("What do you want to do?")
	fmt.Println("1. Create data")
	fmt.Println("2. Solve problem")
	fmt.Println("Enter choice: ")
	fmt.Println("")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')

	switch strings.TrimSpace(input) {
	case "1":
		utils.ClearScreen()
		startTime := time.Now()

		fmt.Println("Hold on while the file is being generated...")
		creator.CreateData()

		fmt.Println("done")
		fmt.Printf("Time taken: %s\n", time.Since(startTime))
	case "2":
		utils.ClearScreen()
		startTime := time.Now()

		fmt.Println("Solving problem...")
		solver.Solve()
		fmt.Printf("Time taken: %s\n", time.Since(startTime))
	default:
		utils.CowSays()
		os.Exit(1)
	}
}
