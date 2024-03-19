package main

import (
	"context"
	"fmt"

	"github.com/jin06/mercury/internal/broker"
	"github.com/spf13/cobra"
)

var cmd cobra.Command

func init() {
	cmd = cobra.Command{
		Use: "MQTT broker!",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("Hi~~ mecury is a mqtt server.")
			b := broker.NewBroker()
			b.Run(context.Background())
		},
	}
}

func main() {
	cmd.Execute()
}
