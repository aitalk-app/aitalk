package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// TODO: remove this ugly global channel
var (
	GlobalInputCh    = make(chan string, 1)
	GlobalInputErrCh = make(chan error, 1)
)

var uploadCmd = &cobra.Command{
	Use:   "upload <file path>",
	Short: "upload a talk in saved file",
	Args: func(cmd *cobra.Command, args []string) error {
		// Optionally run one of the validators provided by cobra
		return cobra.MinimumNArgs(1)(cmd, args)
	},
	Run: func(cmd *cobra.Command, args []string) {
		upload(cmd, args)
	},
}

func upload(cmd *cobra.Command, args []string) {
	for _, f := range args {
		uploadFile(f)
	}
}
func uploadFile(f string) {
	fmt.Printf("Uploading talk %s...\n", f)
	dat, err := os.ReadFile(f)
	if err != nil {
		fmt.Println("open file error: ", err)
		return
	}
	content := string(dat)
	parts := strings.SplitN(content, "...\n\n", 2)
	if len(parts) != 2 {
		fmt.Println("invalid file content")
		return
	}
	headers := parts[0]
	content = strings.TrimSpace(parts[1])
	var (
		model, lang, topic string
		roles              roles
	)
	for _, header := range strings.Split(headers, "\n") {
		header = strings.TrimSpace(header)
		if header == "..." {
			continue
		}
		headerParts := strings.Split(header, ":")
		if len(headerParts) != 2 {
			continue
		}
		name := headerParts[0]
		value := headerParts[1]
		switch name {
		case "model":
			model = value
		case "lang":
			lang = value
		case "topic":
			topic = value
		case "A":
			fallthrough
		case "B":
			roles = append(roles, value)
		}
	}

	uploadTalk(model, lang, roles, topic, content, false)
}

const (
	HOST       = "https://ai-talk.app"
	ConnectURL = HOST + "/connect?install_id=%s"
	TalkURL    = HOST + "/talks/%d"
	UploadAPI  = HOST + "/api/talks/upload"
)

func getInstallID() string {
	parent, err := os.UserHomeDir()
	if err != nil {
		parent = "."
	}
	dir := filepath.Join(parent, ".aitalk")
	makeDir(dir)
	var installID string
	filepath := filepath.Join(dir, "install_id")
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// create a new install id
		file, err := os.Create(filepath)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		defer file.Close()

		installID = uuid.NewString()
		_, err = file.WriteString(installID)
		if err != nil {
			fmt.Println(err)
			return ""
		}
		return installID
	} else if err != nil {
		fmt.Println("failed to check file exists: ", filepath)
		return ""
	}

	// read existed install id
	b, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
	}

	installID = string(b)
	return installID
}

type UploadTalkReq struct {
	Model     string `json:"model" validate:"required"`
	Language  string `json:"lang" validate:"required"`
	RoleA     string `json:"roleA" validate:"required"`
	RoleB     string `json:"roleB" validate:"required"`
	Topic     string `json:"topic" validate:"required"`
	Content   string `json:"content" validate:"required"`
	InstallId string `json:"install_id" validate:"required"`
}

type TalkCreateResp struct {
	ID int64 `json:"id"`
}

type CreateResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data TalkCreateResp `json:"data"`
}

func uploadTalk(model, lang string, roles roles, topic, content string, confirm bool) {
	if confirm {
		printInfo("Press <enter> to upload to https://ai-talk.app, <ctrl-d> to save locally")
		errCh := make(chan error, 1)
		go func() {
			reader := bufio.NewReader(os.Stdin)
			_, err := reader.ReadString('\n')
			errCh <- err
		}()

		select {
		case err := <-errCh:
			if err != nil {
				return
			}
		case err := <-GlobalInputErrCh:
			if err != nil {
				return
			}
		case <-GlobalInputCh:
		}

	}
	var roleA, roleB string
	if len(roles) > 0 {
		roleA, roleB = roles[0], roles[1]
	} else {
		roleA = "AI"
		roleB = "You"
	}
	data := UploadTalkReq{
		Model:     model,
		Topic:     topic,
		Language:  lang,
		RoleA:     roleA,
		RoleB:     roleB,
		Content:   content,
		InstallId: getInstallID(),
	}
	b, err := json.Marshal(data)
	if err != nil {
		printInfo("Marshal request error: %v", err)
		return
	}
	printInfo("Uploading...")
	req, err := http.NewRequest("POST", UploadAPI, bytes.NewBuffer(b))
	if err != nil {
		printInfo("Create request error: %v", err)
		return
	}
	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		printInfo("Send request error: %v", err)
		return
	}
	defer resp.Body.Close()
	var createResp CreateResp
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	if err != nil {
		printInfo("Decode response error: %v", err)
		return
	}
	if createResp.Code == 0 {
		printInfo("Uploaded success, view the talk at:")
		printInfo("    " + fmt.Sprintf(TalkURL, createResp.Data.ID))
	} else {
		printInfo("Uploaded failed: %s", createResp.Msg)
	}
}
