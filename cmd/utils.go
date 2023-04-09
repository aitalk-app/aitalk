package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/shellfly/aoi/pkg/color"
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
	topicSlug := strings.ReplaceAll(strings.ToLower(topic), " ", "_")
	filename := topicSlug
	var path string
	for {
		path = filepath.Join(outputDir, filename) + ".txt"
		if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
			break
		}

		// file existed
		number := 2
		parts := strings.Split(filename, "_")
		oldNumber, err := strconv.Atoi(parts[len(parts)-1])
		if err == nil {
			number = oldNumber + 1
		}
		filename = fmt.Sprintf("%s_%d", topicSlug, number)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		printInfo("failed to open file to save content %s, %v", filename, err)
		return
	}
	defer f.Close()
	header := fmt.Sprintf("model: %s\nlang: %s\ntopic: %s", model, lang, topic)
	if len(roles) == 2 {
		header = fmt.Sprintf("%s\nA: %s\nB: %s", header, roles[0], roles[1])
	}
	_, _ = f.WriteString(fmt.Sprintf(FileHeader, header))
	_, _ = f.WriteString(content)
	printInfo("The output is saved in %s", path)
}

func printInfo(f string, args ...any) {
	fmt.Println(color.Green(fmt.Sprintf(f, args...)))
}
