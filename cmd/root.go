package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "terraseq",
	Short: "https://github.com/enelsr/terraseq",
	Long: `terraseq is a versatile tool designed for managing and transforming DNA
data from popular commercial genetic testing services
https://github.com/enelsr/terraseq`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}


