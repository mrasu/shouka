package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "shouka",
	Short: "Shouka is a tool to support your CI/CD",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() error {
	return rootCmd.Execute()
}
