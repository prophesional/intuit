package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	path := filepath.Base("./test.sql")

	if _, noPathErr := os.Stat(path); noPathErr != nil {
		fmt.Println(noPathErr)
		if mkdirErr := os.MkdirAll(path, 0755); mkdirErr != nil {
			os.Exit(1)
			fmt.Println(mkdirErr)
		}
	}
}
