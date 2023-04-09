package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

const (
	TalkURL   = "https://ai-talk.app/talks/%d"
	UploadAPI = "https://ai-talk.app/api/talks/upload"
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

	b, err := os.ReadFile(filepath) // just pass the file name
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
	Data TalkCreateResp `json:"data"`
}

func uploadTalk(model, lang string, roles roles, topic, content string) {
	fmt.Println("Press <enter> to upload to https://ai-talk.app, <ctrl-d> to save locally")
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	var roleA, roleB string
	if len(roles) > 0 {
		roleA, roleB = roles[0], roles[1]
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
		fmt.Println("Marshal request error: ", err)
		return
	}
	fmt.Println("Uploading...")
	req, err := http.NewRequest("POST", UploadAPI, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Create request error: ", err)
		return
	}
	req.Header.Set("User-Agent", userAgent())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Send request error: ", err)
		return
	}
	defer resp.Body.Close()
	var createResp CreateResp
	err = json.NewDecoder(resp.Body).Decode(&createResp)
	if err != nil {
		fmt.Println("Decode response error:", err)
		return
	}
	fmt.Println("Uploaded success, view the talk at:")
	fmt.Println("    " + fmt.Sprintf(TalkURL, createResp.Data.ID))
}
