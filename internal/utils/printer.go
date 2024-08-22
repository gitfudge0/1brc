package utils

import "fmt"

func CowSays() {
	str := `
 _______________________ 
           lol
 ----------------------- 
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
  `

	fmt.Println(str)
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
