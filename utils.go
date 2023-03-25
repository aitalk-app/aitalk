package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func makeDir(dirName string) string {
	dirName = filepath.Join(dirName)
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
		}
	}
	return dirName
}
