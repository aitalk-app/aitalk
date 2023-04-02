package main

import (
	"bufio"
	"os"
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

func NewAI(name, apiKey, model string) (*AI, error) {
	bot, err := chatgpt.NewAI("", apiKey, model)
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
	defer i.spinner.Stop()

	return i.ai.Query(prompts)
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

	return reader.ReadString('\n')
}
