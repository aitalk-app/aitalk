package cmd

import (
	"bufio"
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/shellfly/aoi/pkg/chatgpt"
)

type Intelligence interface {
	Name() string
	SetRole(string)
	Query(context.Context, []string) (string, error)
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

func (i *AI) Query(ctx context.Context, prompts []string) (reply string, err error) {
	i.spinner.Start()
	replyCh := make(chan string, 1)
	errCh := make(chan error, 1)
	go func() {
		defer i.spinner.Stop()
		reply, err := i.ai.Query(prompts)
		if err != nil {
			errCh <- err
			return
		}
		replyCh <- reply
	}()

	select {
	case reply = <-replyCh:
		reply = strings.TrimSpace(reply)
		for _, char := range reply {
			fmt.Print(string(char))
			os.Stdout.Sync()
			time.Sleep(time.Duration(rand.Intn(42)+10) * time.Millisecond)
		}
		fmt.Println()
		fmt.Println()
	case err = <-errCh:
		fmt.Println("failed to get reply: ", err)
	case <-ctx.Done():
		i.spinner.Stop()
	}

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

func (i *Human) Query(ctx context.Context, prompts []string) (reply string, err error) {
	replyCh := make(chan string, 1)
	errCh := make(chan error, 1)
	reader := bufio.NewReader(os.Stdin)
	go func() {
		for reply == "" {
			reply, err = reader.ReadString('\n')
			select {
			case <-ctx.Done():
				if err == nil {
					GlobalInputCh <- reply
				} else {
					GlobalInputErrCh <- err
				}
				return
			default:
			}
			if err != nil {
				errCh <- err
				return
			}

			reply = strings.TrimSpace(reply)
		}

		replyCh <- reply
	}()

	select {
	case reply = <-replyCh:
		fmt.Println()
	case err = <-errCh:
		fmt.Println("failed to get input: ", err)
	case <-ctx.Done():
	}

	return reply, err
}
