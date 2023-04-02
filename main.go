package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

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
	flag.DurationVar(&duration, "timeout", 5*time.Minute, "timeout for the talk")
	flag.StringVar(&model, "model", "gpt-3.5-turbo", "model to use")
	flag.StringVar(&openaiAPIKey, "openai_api_key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	flag.StringVar(&lang, "lang", "", "language")
	flag.StringVar(&topic, "topic", "", "topic")
	flag.Var(&roles, "role", "list of roles")
	flag.Parse()
	if topic == "" {
		usageExit()
	}

	var (
		A, B Intelligent
		err  error
	)
	if len(roles) == 0 {
		A, err = NewAI("AI", openaiAPIKey, model)
		if err != nil {
			fmt.Println(err)
			return
		}
		B = &Human{name: "You"}
	} else if len(roles) == 2 {
		A, err = NewAI("A", openaiAPIKey, model)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		A.SetRole(systemPrompt(roles[0], lang))
		B, err = NewAI("B", openaiAPIKey, model)
		if err != nil {
			fmt.Println(err)
			return
		}
		B.SetRole(systemPrompt(roles[1], lang))
	} else {
		usageExit()
	}

	outputDir := makeDir("talks")
	discussions := []string{}
	filename := topic + ".txt"
	path := filepath.Join(outputDir, filename)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("failed to open ", filename, err)
		return
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Reset(os.Interrupt)
	go func() {
		if B.IsAI() {
			select {
			case <-sigChan:
			case <-time.After(duration):
				fmt.Println("[timeout]")
			}
		} else {
			<-sigChan
		}
		if len(discussions) == 0 {
			return
		}

		_, _ = f.WriteString(strings.Join(discussions, "\n\n"))
		f.Close()
		fmt.Println("The output is saved in ", path)
		os.Exit(0)
	}()

	fmt.Printf("topic: %s\n", topic)
	if A.IsAI() && B.IsAI() {
		fmt.Printf("A: %s\n", roles[0])
		fmt.Printf("B: %s\n", roles[1])
	}

	var (
		promptsA, promptsB []string
		replyA, replyB     string
	)
	promptsA = []string{topic}
	for {
		fmt.Println(color.Yellow(A.Name() + ":"))
		replyA, err = A.Query(promptsA)
		if err != nil {
			fmt.Println("failed to get reply from A: ", err)
			return
		}
		fmt.Println(strings.TrimSpace(replyA))
		fmt.Println()
		discussions = append(discussions, A.Name()+": "+strings.TrimSpace(replyA))

		fmt.Println(color.Green(B.Name() + ":"))
		replyB, err = B.Query(promptsB)
		if err != nil {
			fmt.Println("failed to get reply from B: ", err)
			return
		}
		if B.IsAI() {
			fmt.Println(strings.TrimSpace(replyB))
		}
		fmt.Println()
		discussions = append(discussions, B.Name()+": "+strings.TrimSpace(replyB))

		promptsA = []string{replyB}
	}
}

func systemPrompt(role, lang string) string {
	system := fmt.Sprintf("From now on, we will have a debate. Your viewpoint is that %s. You must stick to your viewpoint and never agreeing with mine. do NOT say 'I understand on your point', limit up to 100 words for every reply", role)
	language := command.Languages[lang]
	if language != "" {
		system = system + " Please reply in " + language
	}
	return system
}

// func topicPrompt(topic string) string {
// 	return fmt.Sprintf(`I want to discuss with you on "%s".`, topic)
// }

func usageExit() {
	fmt.Println("Usage:")
	fmt.Println("aitalk -topic {topic}")
	fmt.Println("aitalk -topic {topic} -role '{role description}' -role '{role description}'")
	os.Exit(0)
}
