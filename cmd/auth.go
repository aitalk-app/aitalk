/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const AuthHint = `Open the URL in a web browser to link your install ID with your account:

https://ai-talk.app/connect?install_id=%s


This will associate all talks uploaded from this machine (past and future ones) to your account, and allow you to manage them at ai-talk.app.
`

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "authenticate and connect talks",
	Long:  `connect talks generated on this machine to an ai-talk.app account.`,
	Run: func(cmd *cobra.Command, args []string) {
		hint := fmt.Sprintf(AuthHint, getInstallID())
		fmt.Println(hint)
	},
}
