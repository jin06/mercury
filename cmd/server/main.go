package main

import (
	"context"
	"fmt"

	"github.com/jin06/mercury/internal/broker"
	"github.com/jin06/mercury/internal/config"
	"github.com/spf13/cobra"
)

var cmd cobra.Command

func init() {
	cmd = cobra.Command{
		Use: "MQTT broker!",
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Println("Hi~~ mecury is a mqtt server.")
			if path, err := c.Flags().GetString("config"); err != nil {
				return err
			} else if err := config.Init(path); err != nil {
				return err
			}
			b := broker.NewBroker()
			b.Run(context.Background())
			return nil
		},
	}
	cmd.PersistentFlags().String("config", "mercury.yaml", "Specify config file path")
}

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
