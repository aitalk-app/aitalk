package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/shellfly/aoi/pkg/color"
	"github.com/shellfly/aoi/pkg/command"
	"github.com/spf13/cobra"
)

const NAME = "aitalk"

var (
	VERSION = "0.2.0"

	Version                     bool
	Topic, Lang, Model          string
	OpenaiAPIHost, OpenaiAPIKey string
	Roles                       roles
	Duration                    time.Duration
)

var rootCmd = &cobra.Command{
	Use:   "aitalk",
	Short: "aitalk is CLI tool to generate and share ai talks",
	Long:  `Complete documentation is available at https://ai-talk.app`,
	Run: func(cmd *cobra.Command, args []string) {
		main(cmd, args)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&Version, "version", "v", false, "print version")
	rootCmd.Flags().StringVarP(&Topic, "topic", "", "", "topic (required)")
	rootCmd.Flags().StringVarP(&Model, "model", "", "gpt-3.5-turbo", "model to use")
	rootCmd.Flags().StringVarP(&Lang, "lang", "", "en", "language")
	rootCmd.Flags().StringVarP(&OpenaiAPIHost, "openai_api_host", "", os.Getenv("OPENAI_API_HOST"), "OpenAI API host")
	rootCmd.Flags().StringVarP(&OpenaiAPIKey, "openai_api_key", "", os.Getenv("OPENAI_API_KEY"), "OpenAI API key")
	rootCmd.Flags().VarP(&Roles, "role", "", "list of roles")
	rootCmd.Flags().DurationVarP(&Duration, "timeout", "", 5*time.Minute, "timeout for the talk")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(langCmd)
	rootCmd.AddCommand(uploadCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func userAgent() string {
	return fmt.Sprintf("%s/%s", NAME, VERSION)
}

type roles []string

func (rs *roles) Type() string {
	return "Role"
}

func (rs *roles) String() string {
	return strings.Join(*rs, ",")
}

func (rs *roles) Set(value string) error {
	*rs = append(*rs, value)
	return nil
}

func main(cmd *cobra.Command, args []string) {
	if Version {
		fmt.Println(userAgent())
		return
	}

	if Topic == "" {
		usageExit()
	}

	var (
		A, B Intelligence
		err  error
	)
	if len(Roles) == 0 {
		A, err = NewAI("AI", OpenaiAPIHost, OpenaiAPIKey, Model)
		if err != nil {
			fmt.Println(err)
			return
		}
		B = &Human{name: "You"}
	} else if len(Roles) == 2 {
		A, err = NewAI("A", OpenaiAPIHost, OpenaiAPIKey, Model)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		A.SetRole(systemPrompt(Roles[0], Lang))
		B, err = NewAI("B", OpenaiAPIHost, OpenaiAPIKey, Model)
		if err != nil {
			fmt.Println(err)
			return
		}
		B.SetRole(systemPrompt(Roles[1], Lang))
	} else {
		usageExit()
	}

	ctx, cancel := context.WithCancel(context.Background())
	if B.IsAI() {
		ctx, cancel = context.WithTimeout(ctx, Duration)
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
	promptsA = []string{Topic}
	discussions := []string{}
loop:
	for {
		fmt.Println(color.Yellow(A.Name() + ":"))
		replyA, err = A.Query(ctx, promptsA)
		if err != nil {
			return
		}
		if replyA != "" {
			// when ctrl-c replyA could be empty
			discussions = append(discussions, A.Name()+": "+strings.TrimSpace(replyA))
		}
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				fmt.Println("[timed out]")
			}
			break loop
		default:
		}

		fmt.Println(color.Green(B.Name() + ":"))
		replyB, err = B.Query(ctx, promptsB)
		if err != nil {
			return
		}
		if replyB != "" {
			// when ctrl-c replyB could be empty
			discussions = append(discussions, B.Name()+": "+strings.TrimSpace(replyB))
		}
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
	saveTalk(Model, Lang, Roles, Topic, content)
	uploadTalk(Model, Lang, Roles, Topic, content, true)
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
	fmt.Println(`Usage
  1. Generate a talk with AI interactively:
     aitalk --topic {topic}

  2. Generate a talk automatically:
     aitalk --topic {topic} --lang {lang} --role '{role description}' --role '{role description}'`)
	os.Exit(0)
}
