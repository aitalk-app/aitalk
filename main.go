package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/shellfly/aoi/pkg/chatgpt"
	"github.com/shellfly/aoi/pkg/color"
)

var languagePrompts = map[string]string{
	"en":    "",
	"cn":    "Please reply in Chinese",
	"jp":    "Please reply in Japanese",
	"fr":    "Please reply in France",
	"de":    "Please reply in German",
	"ru":    "Please reply in Russia",
	"zh-tw": "Please reply in Traditional Chinese",
}

func main() {
	var model, openaiAPIKey string
	var lang, roleA, roleB, topic string
	flag.StringVar(&model, "model", "gpt-3.5-turbo", "model to use")
	flag.StringVar(&openaiAPIKey, "openai_api_key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	flag.StringVar(&lang, "lang", "en", "language")
	flag.StringVar(&roleA, "roleA", "", "role A")
	flag.StringVar(&roleB, "roleB", "", "role B")
	flag.StringVar(&topic, "topic", "", "topic")
	flag.Parse()

	if roleA == "" || roleB == "" || topic == "" {
		fmt.Println("Usage:")
		fmt.Println("talkgpt -roleA {roleA} -roleB {roleB} -topic {topic}")
	}

	// Create an AI
	aiA, err := chatgpt.NewAI(openaiAPIKey, model)
	if err != nil {
		fmt.Println("create ai1 error: ", err)
		return
	}
	aiA.SetSystem(systemPrompt(roleA))

	aiB, err := chatgpt.NewAI(openaiAPIKey, model)
	if err != nil {
		fmt.Println("create ai2 error: ", err)
		return
	}
	aiB.SetSystem(systemPrompt(roleB))

	// Setup signal handler
	discussions := []string{}
	filename := topic + ".txt"
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open ", filename, err)
		return
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	go func() {
		<-sigChan
		os.Stdout.Sync()
		_, _ = f.WriteString(strings.Join(discussions, "\n\n"))
		os.Exit(0)
	}()

	// prepare first prompt
	var (
		promptsA, promptsB []string
		replyA, replyB     string
		spinnerA           = spinner.New(spinner.CharSets[14], 200*time.Millisecond, spinner.WithColor("yellow"), spinner.WithSuffix(" thinking..."))
		spinnerB           = spinner.New(spinner.CharSets[14], 200*time.Millisecond, spinner.WithColor("green"), spinner.WithSuffix(" thinking..."))
	)
	promptsA = []string{topicPrompt(topic), "You start first"}
	langPrompt := languagePrompts[lang]
	if langPrompt != "" {
		promptsA = append(promptsA, langPrompt)
	}

	fmt.Printf("A: %s\nB: %s\ntopic: %s\n", roleA, roleB, topic)
	fmt.Println(color.Yellow("A:"))
	spinnerA.Start()
	replyA, err = aiA.Query(promptsA)
	if err != nil {
		fmt.Println("failed to get reply from A: ", err)
		return
	}
	spinnerA.Stop()
	fmt.Println(replyA)
	fmt.Println()
	promptsB = []string{topicPrompt(topic), replyA}
	for {
		discussions = append(discussions, "A: "+replyA)

		fmt.Println(color.Green("B: "))
		spinnerB.Start()
		replyB, err = aiB.Query(promptsB)
		spinnerB.Stop()
		if err != nil {
			fmt.Println("failed to get reply from B: ", err)
			return
		}
		discussions = append(discussions, "B: "+replyB)
		fmt.Println(replyB)
		fmt.Println()

		promptsA = []string{replyB}
		fmt.Println(color.Yellow("A: "))
		spinnerA.Start()
		replyA, err = aiA.Query(promptsA)
		spinnerA.Stop()
		if err != nil {
			fmt.Println("failed to get reply from B: ", err)
			return
		}
		fmt.Println(replyA)
		fmt.Println()

		promptsB = []string{replyA}
	}

}

func systemPrompt(role string) string {
	return fmt.Sprintf("You are %s, You don't agree anything I say to you and give opposite argument", role)
}

func topicPrompt(topic string) string {
	return fmt.Sprintf(`I want to have a discussion with you on "%s".
	Please be specific as much as possible, don't use empty or ambiguity words`, topic)
}
