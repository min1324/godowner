package main

import (
	"fmt"
)

func main() {
	// url := "https://dl.google.com/go/go1.11.1.src.tar.gz"
	err := CMD()
	if err != nil {
		fmt.Println(err)
	}
}
