package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const FileHeader = `...
%s
...


`

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

func saveTalk(model, lang string, roles roles, topic, content string) {
	outputDir := makeDir("talks")
	filename := strings.ReplaceAll(strings.ToLower(topic), " ", "_")
	path := filepath.Join(outputDir, filename) + ".txt"
	if _, err := os.Stat(path); err == nil {
		// file existed
		number := 2
		parts := strings.Split(filename, "_")
		oldNumber, err := strconv.Atoi(parts[len(parts)-1])
		if err == nil {
			number = oldNumber + 1
		}
		path = filepath.Join(outputDir, fmt.Sprintf("%s_%d", filename, number)) + ".txt"
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open file to save content ", filename, err)
		return
	}
	defer f.Close()
	header := fmt.Sprintf("model: %s\nlang: %s\ntopic: %s", model, lang, topic)
	if len(roles) == 2 {
		header = fmt.Sprintf("%s\nA: %s\nB: %s", header, roles[0], roles[1])
	}
	_, _ = f.WriteString(fmt.Sprintf(FileHeader, header))
	_, _ = f.WriteString(content)
	fmt.Println("The output is saved in ", path)
}
