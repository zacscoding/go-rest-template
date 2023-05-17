package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zacscoding/go-rest-template/pkg/version"
)

func init() {
	rootCmd.AddCommand(versionCommand)
}

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Print version info",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("API Server %s", version.ShortVersionInfo())
		fmt.Println()
	},
}
