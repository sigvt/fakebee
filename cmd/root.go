/*
Copyright © 2022 Daniils Petrovs <thedanpetrov@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/holodata/fakebee/ytl"
	"github.com/pioz/faker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fakebee",
	Short: "A fake Youtube Live event producer.",
	Long: `Produces fake Youtube Live events to Kafka to different topics 
	and with different message payloads.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	faker.SetSeed(12345)
	ytl.RegisterBuilders()

	rootCmd.PersistentFlags().String("broker", "localhost:9092", "Kafka broker to connect to")
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
}
