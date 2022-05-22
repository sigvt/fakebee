/*
Copyright © 2022 Daniils Petrovs <thedanpetrov@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/holodata/fakebee/ytl"
	"github.com/pioz/faker"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fakebee",
	Short: "A fake Youtube Live event producer.",
	Long: `Produces fake Youtube Live events to Kafka to different topics 
	and with different message payloads.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.fakebee.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	faker.SetSeed(12345)
	ytl.RegisterBuilders()

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
