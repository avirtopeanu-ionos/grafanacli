package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

func main() {
	// Seed the random number generator.
	rand.Seed(time.Now().UnixNano())

	var name string
	var shout bool

	// Define the "greet" command.
	var greetCmd = &cobra.Command{
		Use:   "greet",
		Short: "Greet someone",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Since --name is marked as required, this check is extra safety.
			if name == "" {
				fmt.Fprintln(os.Stderr, "Error: --name flag is required")
				os.Exit(1)
			}

			message := fmt.Sprintf("Hello and welcome, %s!", name)

			// If --shout is provided, convert the message to uppercase.
			if shout {
				message = strings.ToUpper(message)

				if rand.Float64() < 0.1 {
					if rand.Float64() < 0.5 {
						// Simulated 400 Bad Request (exit with error).
						return fmt.Errorf("400 Bad Request")
					} else {
						// Wrong output due to cough.
						fmt.Println("HE--<BLERGH>")
						return nil
					}
				}
			}

			fmt.Println(message)
			return nil
		},
	}

	// Define flags for the greet command.
	greetCmd.Flags().StringVar(&name, "name", "", "Name to greet (required)")
	greetCmd.MarkFlagRequired("name")
	greetCmd.Flags().BoolVar(&shout, "shout", false, "Shout the greeting (all caps)")

	// Create the root command and add the greet command.
	var rootCmd = &cobra.Command{Use: "mycli"}
	rootCmd.AddCommand(greetCmd)

	// Execute the root command.
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
