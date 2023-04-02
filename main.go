package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
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
		version             bool
	)
	flag.DurationVar(&duration, "timeout", 5*time.Minute, "timeout for the talk")
	flag.StringVar(&model, "model", "gpt-3.5-turbo", "model to use")
	flag.StringVar(&openaiAPIKey, "openai_api_key", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	flag.StringVar(&lang, "lang", "en", "language")
	flag.StringVar(&topic, "topic", "", "topic")
	flag.BoolVar(&version, "version", false, "topic")
	flag.Var(&roles, "role", "list of roles")
	flag.Parse()
	if version {
		fmt.Println(userAgent())
		return
	}
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

	ctx, cancel := context.WithCancel(context.Background())
	if B.IsAI() {
		ctx, cancel = context.WithTimeout(ctx, duration)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	defer signal.Reset()
	go func() {
		<-sigChan
		cancel()
	}()

	var (
		promptsA, promptsB []string
		replyA, replyB     string
	)
	promptsA = []string{topic}
	discussions := []string{}
loop:
	for {
		fmt.Println(color.Yellow(A.Name() + ":"))
		replyA, err = A.Query(promptsA)
		if err != nil {
			return
		}
		discussions = append(discussions, A.Name()+": "+strings.TrimSpace(replyA))
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("[timed out]")
			}
			break loop
		default:
		}

		fmt.Println(color.Green(B.Name() + ":"))
		replyB, err = B.Query(promptsB)
		if err != nil {
			return
		}
		discussions = append(discussions, B.Name()+": "+strings.TrimSpace(replyB))
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("[timed out]")
			}
			break loop
		default:
		}

		promptsA = []string{replyB}
	}

	content := strings.Join(discussions, "\n\n")
	if content == "" {
		return
	}
	saveTalk(topic, content)
	uploadTalk(model, topic, lang, content, roles)
}

func systemPrompt(role, lang string) string {
	system := fmt.Sprintf("From now on, we will have a debate. Your viewpoint is that %s. You must stick to your viewpoint and never agreeing with mine. do NOT say 'I understand on your point', limit up to 100 words for every reply", role)
	language := command.Languages[lang]
	if language != "" {
		system = system + " Please reply in " + language
	}
	return system
}

func usageExit() {
	fmt.Println("Usage:")
	fmt.Println("aitalk -topic {topic}")
	fmt.Println("aitalk -topic {topic} -role '{role description}' -role '{role description}'")
	os.Exit(0)
}
