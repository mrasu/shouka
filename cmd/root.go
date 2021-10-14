package cmd

import (
	"github.com/mrasu/shouka/libs/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shouka",
	Short: "Shouka is a tool to support your CI/CD",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var verbose bool

func Execute() error {
	cobra.OnInitialize(initialize)

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	return rootCmd.Execute()
}

func initialize() {
	if verbose {
		log.SetLevel(log.Debug)
	}
}
