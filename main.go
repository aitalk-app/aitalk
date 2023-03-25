package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/shellfly/aoi/pkg/chatgpt"
	"github.com/shellfly/aoi/pkg/color"
	"github.com/shellfly/aoi/pkg/command"
)

type roles []string

func (rs *roles) String() string {
	return strings.Join(*rs, ",")
}

func (rs *roles) Set(value string) error {
	*rs = append(*rs, value)
	return nil
}

func main() {
	var (
		model, openaiAPIKey string
		lang, topic         string
		roles               roles
		duration            time.Duration
	)
	flag.DurationVar(&duration, "timeout", 2*time.Minute, "timeout for the talk")
	flag.StringVar(&model, "model", "gpt-3.5-turbo", "model to use")
	flag.StringVar(&openaiAPIKey, "openai_api_key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	flag.StringVar(&lang, "lang", "en", "language")
	flag.StringVar(&topic, "topic", "", "topic")
	flag.Var(&roles, "role", "list of roles")
	flag.Parse()
	if topic == "" {
		usageExit()
	}
	if len(roles) <= 1 {
		AIToHuman(openaiAPIKey, model, topic, lang, roles)
		return
	} else if len(roles) == 2 {
		AIToAI(openaiAPIKey, model, topic, lang, roles, duration)
		return
	}
	usageExit()
}

func AIToHuman(apiKey, model, topic, lang string, roles roles) {
	aiA, err := chatgpt.NewAI("", apiKey, model)
	if err != nil {
		fmt.Println("create ai1 error: ", err)
		return
	}
	if len(roles) == 1 {
		aiA.SetSystem(roles[0])
	}
	outputDir := makeDir("talks")
	// Setup signal handler
	discussions := []string{}
	filename := topic + ".txt"
	path := filepath.Join(outputDir, filename)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
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
		f.Close()
		fmt.Println("The output is saved in ", path)
		os.Exit(0)
	}()

	// prepare first prompt
	var (
		promptsA []string
		replyA   string
		spinnerA = spinner.New(spinner.CharSets[14], 200*time.Millisecond, spinner.WithColor("yellow"), spinner.WithSuffix(" thinking..."))
	)
	promptsA = []string{topicPrompt(topic)}
	language := command.Languages[lang]
	if language != "" {
		promptsA = append(promptsA, "Please replay in "+language)
	}

	fmt.Printf("topic: %s\n", topic)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(color.Yellow("AI:"))
		spinnerA.Start()
		replyA, err = aiA.Query(promptsA)
		if err != nil {
			fmt.Println("failed to get reply from A: ", err)
			return
		}
		spinnerA.Stop()
		fmt.Println(replyA)
		fmt.Println()
		discussions = append(discussions, "AI: \n"+replyA)

		fmt.Println(color.Green("You:"))
		input, _ := reader.ReadString('\n')
		discussions = append(discussions, "You: \n"+input)

		promptsA = []string{input}
	}
}

func AIToAI(apiKey, model string, topic, lang string, roles roles, timeout time.Duration) {
	roleA, roleB := roles[0], roles[1]
	// Create an AI
	aiA, err := chatgpt.NewAI("", apiKey, model)
	if err != nil {
		fmt.Println("create ai1 error: ", err)
		return
	}
	aiA.SetSystem(systemPrompt(roleA))

	aiB, err := chatgpt.NewAI("", apiKey, model)
	if err != nil {
		fmt.Println("create ai2 error: ", err)
		return
	}
	aiB.SetSystem(systemPrompt(roleB))

	// Setup signal handler
	outputDir := makeDir("talks")
	// Setup signal handler
	discussions := []string{}
	filename := topic + ".txt"
	path := filepath.Join(outputDir, filename)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("failed to open ", filename, err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	go func() {
		select {
		case <-sigChan:
		case <-time.After(timeout):
			fmt.Println("timeout")
		}

		os.Stdout.Sync()
		_, _ = f.WriteString(strings.Join(discussions, "\n\n"))
		f.Close()
		fmt.Println("The output is saved in ", path)
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
	language := command.Languages[lang]
	if language != "" {
		promptsA = append(promptsA, "Please replay in "+language)
	}

	fmt.Printf("A: %s\nB: %s\ntopic: %s\n", roleA, roleB, topic)

	for {
		fmt.Println(color.Yellow("A:"))
		spinnerA.Start()
		replyA, err = aiA.Query(promptsA)
		if err != nil {
			fmt.Println("failed to get reply from A: ", err)
			return
		}
		spinnerA.Stop()
		discussions = append(discussions, "A: \n"+replyA)
		printReply(replyA)

		promptsB = []string{topicPrompt(topic), replyA}
		fmt.Println(color.Green("B: "))
		spinnerB.Start()
		replyB, err = aiB.Query(promptsB)
		spinnerB.Stop()
		if err != nil {
			fmt.Println("failed to get reply from B: ", err)
			return
		}
		discussions = append(discussions, "B: \n"+replyB)
		printReply(replyB)

		promptsA = []string{replyB}
	}

}

func systemPrompt(role string) string {
	return fmt.Sprintf("From now on, we will have a debate. Your viewpoint is that %s. You must stick to your viewpoint and never agreeing with mine. do NOT say 'I understand on your point'", role)
}

func topicPrompt(topic string) string {
	return fmt.Sprintf(`I want to discuss with you on "%s".`, topic)
}

func usageExit() {
	fmt.Println("Usage:")
	fmt.Println("aitalk -topic {topic}")
	fmt.Println("aitalk -topic {topic} -role '{role description}' -role '{role description}'")
	os.Exit(0)
}

func printReply(text string) {
	for _, char := range text {
		fmt.Print(string(char))
		time.Sleep(time.Duration(rand.Intn(50)+50) * time.Millisecond)
	}
	fmt.Println()
}
