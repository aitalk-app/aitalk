/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"sort"

	"github.com/shellfly/aoi/pkg/command"
	"github.com/spf13/cobra"
)

var langCmd = &cobra.Command{
	Use:   "lang",
	Short: "show supported languages",
	Run: func(cmd *cobra.Command, args []string) {
		languages := [][]string{}
		for lang, language := range command.Languages {
			languages = append(languages, []string{lang, language})
		}
		sort.Slice(languages, func(i, j int) bool {
			return languages[i][0] < languages[j][0]
		})
		for _, l := range languages {
			fmt.Println(l[0], l[1])
		}
	},
}
