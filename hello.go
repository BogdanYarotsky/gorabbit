package main

import (
	"fmt"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gorabbit",
	Short: "RabbitMQ CLI",
}

var declareCmd = &cobra.Command{
	Use:   "declare [resource]",
	Short: "Declare a resource like exchange, queue, etc.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resource := args[0]
		switch resource {
		case "exchange":
			name, _ := cmd.Flags().GetString("name")
			exchangeType, _ := cmd.Flags().GetString("type")
			fmt.Printf("Declaring exchange: %s of type %s\n", name, exchangeType)
			conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
			if err != nil {
				log.Fatalf("Failed to connect to RabbitMQ: %v", err)
			}
			defer conn.Close()

			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("Failed to open a channel: %v", err)
			}
			defer ch.Close()
			err = ch.ExchangeDeclare(
				name,         // name
				exchangeType, // type
				true,         // durable
				false,        // auto-deleted
				false,        // internal
				false,        // no-wait
				nil,          // arguments
			)
			if err != nil {
				log.Fatalf("Failed to declare an exchange: %v", err)
			}
		}
	},
}

func init() {
	declareCmd.Flags().StringP("name", "n", "", "Name of the exchange")
	declareCmd.Flags().StringP("type", "t", "", "Type of the exchange")
	declareCmd.MarkFlagRequired("name")
	declareCmd.MarkFlagRequired("type")
	rootCmd.AddCommand(declareCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
