/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/kalru/git-worktree/pkg/switchMenu"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editor string

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch",
	Short: "Switch worktrees",
	Long: `Select current worktrees from a menu to switch to them.

Use / to search and filter by name.`,
	Run: func(cmd *cobra.Command, args []string) {
		switchMenu.Run()
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// switchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// switchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&editor, "editor", "e", "", "Code editor to use when opening worktrees")
	viper.BindPFlag("editor", rootCmd.PersistentFlags().Lookup("editor"))
}
