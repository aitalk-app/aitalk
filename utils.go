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

func saveTalk(topic, content string) {
	outputDir := makeDir("talks")
	filename := topic + ".txt"
	path := filepath.Join(outputDir, filename)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("failed to open file to save content ", filename, err)
		return
	}
	defer f.Close()
	_, _ = f.WriteString(content)
	fmt.Println("The output is saved in ", path)
}
