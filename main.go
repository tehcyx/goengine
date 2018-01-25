package main

import (
	"fmt"
)

func main() {
	printBanner()
}

func printBanner() {
	// http://patorjk.com/software/taag/#p=display&f=Rectangles&t=gomanager
	fmt.Println()
	fmt.Printf(`
	Welcome to
                                _         
            ___ ___ ___ ___ ___|_|___ ___ 
           | . | . | -_|   | . | |   | -_|
           |_  |___|___|_|_|_  |_|_|_|___|
           |___|           |___|          `)
	fmt.Println()
	fmt.Println()
	fmt.Println()
}
