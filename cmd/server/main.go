package main

import (
	"os"

	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Application server",
	Run:   runApplication,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "f", "", "config file")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
