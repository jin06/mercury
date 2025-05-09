package main

import (
	"context"

	"github.com/common-nighthawk/go-figure"
	"github.com/jin06/mercury/internal/broker"
	"github.com/jin06/mercury/internal/config"
	"github.com/spf13/cobra"
)

var cmd cobra.Command

func init() {
	cmd = cobra.Command{
		Use: "MQTT broker!",
		RunE: func(c *cobra.Command, args []string) error {
			// fmt.Println("Hi~~ mercury is a mqtt server.")
			myFigure := figure.NewFigure("Mercury MQTT Broker", "", true)
			myFigure.Print()
			if path, err := c.Flags().GetString("config"); err != nil {
				return err
			} else if err := config.Init(path); err != nil {
				return err
			}
			b := broker.NewBroker()
			ctx := context.TODO()
			if err := b.Run(ctx); err != nil {
				panic(err)
			}
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
