package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmd cobra.Command

func init() {
	cmd = cobra.Command{
		Use: "MQTT broker!",
		Run: func(c *cobra.Command, args []string) {
			fmt.Println("Hi~~ mecury is a mqtt server.")

		},
	}
}

func main() {
	cmd.Execute()
}
