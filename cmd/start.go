/*
Copyright © 2022 Daniils Petrovs <thedanpetrov@gmail.com>

*/
package cmd

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/holodata/fakebee/bee"
	"github.com/holodata/fakebee/ytl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Origin struct {
	VideoId   string   `json:"videoId"`
	ChannelId string   `json:"channelId`
	Topics    []string `json:"topics"`
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start producing events with FakeBee",
	Run: func(cmd *cobra.Command, args []string) {
		var wg sync.WaitGroup

		origins := []Origin{}
		viper.UnmarshalKey("origins", &origins)

		log.Printf("Origins configued: %+v\n", origins)

		backend, _ := cmd.Flags().GetString("backend")
		broker, _ := cmd.InheritedFlags().GetString("broker")

		viper.Set("broker", broker)

		for _, origin := range origins {
			for _, topic := range origin.Topics {
				worker := bee.NewEventWorker(
					bee.WithTopic(topic),
					bee.WithOrigin(ytl.Origin{ChannelId: origin.ChannelId, VideoId: origin.VideoId}),
					bee.WithBackend(backend),
				)

				// Spread out the timing of the tickers
				time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
				worker.Run(&wg)
				wg.Add(1)
			}
		}

		// Wait for all workers to finish
		wg.Wait()
		log.Println("Finished all jobs!")
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP("backend", "b", "printer", "Producer backend (either printer or kafka)")
}
