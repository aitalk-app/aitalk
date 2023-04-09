/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "authenticate and connect talks",
	Long:  `connect talks generated on this machine to an ai-talk.app account.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auth called")
	},
}
