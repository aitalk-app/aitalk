package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/shellfly/aoi/pkg/chatgpt"
)

type Intelligent interface {
	Name() string
	SetRole(string)
	Query([]string) (string, error)
	IsAI() bool
}

type AI struct {
	name    string
	ai      *chatgpt.AI
	spinner *spinner.Spinner
}

func NewAI(name, host, apiKey, model string) (*AI, error) {
	bot, err := chatgpt.NewAI(host, apiKey, model)
	if err != nil {
		return nil, err
	}
	//bot.ToggleDebug()
	return &AI{
		name:    name,
		ai:      bot,
		spinner: spinner.New(spinner.CharSets[14], 200*time.Millisecond, spinner.WithSuffix(" thinking...")),
	}, nil
}

func (i *AI) Name() string {
	return i.name
}

func (i *AI) IsAI() bool {
	return true
}

func (i *AI) SetRole(role string) {
	i.ai.SetSystem(role)
}

func (i *AI) Query(prompts []string) (string, error) {
	i.spinner.Start()
	reply, err := i.ai.Query(prompts)
	i.spinner.Stop()
	if err != nil {
		fmt.Println("failed to get reply: ", err)
		return "", err
	}
	reply = strings.TrimSpace(reply)
	for _, char := range reply {
		fmt.Print(string(char))
		time.Sleep(time.Duration(rand.Intn(42)+10) * time.Millisecond)
	}
	fmt.Println()
	fmt.Println()
	return reply, err
}

type Human struct {
	name string
}

func (i *Human) Name() string {
	return i.name
}

func (i *Human) IsAI() bool {
	return false
}

func (i *Human) SetRole(string) {}

func (i *Human) Query(prompts []string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	reply, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("failed to get input: ", err)
		return "", err
	}
	fmt.Println()
	return reply, err
}
