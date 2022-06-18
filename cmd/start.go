/*
Copyright © 2022 Daniils Petrovs <thedanpetrov@gmail.com>

*/
package cmd

import (
	"fmt"
	"sync"

	"github.com/holodata/fakebee/bee"
	"github.com/holodata/fakebee/ytl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start producing events with FakeBee",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		origin := ytl.Origin{ChannelId: "UCyl1z3jo3XHR1riLFKG5UAg", VideoId: "HtGA1DfQr9o"}

		backend, _ := cmd.Flags().GetString("backend")
		broker, _ := cmd.InheritedFlags().GetString("broker")

		viper.Set("broker", broker)

		bee.NewEventWorker(bee.WithTopic("superchats"), bee.WithOrigin(origin), bee.WithBackend(backend)).Run(&wg)
		bee.NewEventWorker(bee.WithTopic("chats"), bee.WithOrigin(origin), bee.WithBackend(backend)).Run(&wg)
		wg.Add(2)

		// Wait for all workers to finish
		wg.Wait()
		fmt.Println("Finished all jobs!")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("backend", "b", "printer", "Producer backend (either printer or kafka)")
}
