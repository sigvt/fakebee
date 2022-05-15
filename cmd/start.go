/*
Copyright © 2022 Daniils Petrovs <thedanpetrov@gmail.com>

*/
package cmd

import (
	"fmt"
	"sync"

	"github.com/holodata/fakebee/bee"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start producing events with FakeBee",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		bee.NewEventWorker("UCyl1z3jo3XHR1riLFKG5UAg", "HtGA1DfQr9o", "chats").Run(&wg)
		bee.NewEventWorker("UCyl1z3jo3XHR1riLFKG5UAg", "HtGA1DfQr9o", "superchats").Run(&wg)
		wg.Add(2)

		// Wait for all workers to finish
		wg.Wait()
		fmt.Println("Finished all jobs!")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
