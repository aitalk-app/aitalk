package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
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

func printReply(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(time.Duration(rand.Intn(42)+42) * time.Millisecond)
	}
	fmt.Println()
}
