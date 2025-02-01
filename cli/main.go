package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"mycli/internal/telemetry"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func main() {
	telemetry.InitLogger()
	defer telemetry.CatchPanic()

	var name string
	var shout bool

	var greetCmd = &cobra.Command{
		Use:   "greet",
		Short: "Greet someone",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" {
				log.Error().Msg("--name flag is required")
				fmt.Fprintln(os.Stderr, "Error: --name flag is required")
				os.Exit(1)
			}

			start := time.Now()

			message := fmt.Sprintf("Hello and welcome %s!", name)
			flags := map[string]string{
				"name":  name,
				"shout": fmt.Sprintf("%t", shout),
			}

			// If --shout is provided, convert the message to uppercase
			// and simulate a potential error.
			if shout {
				message = strings.ToUpper(message)
				// Simulate a 10% chance for a "cough" error.
				if randFloat() < 0.1 {
					// Further, a 50% chance to simulate a 400 error.
					if randFloat() < 0.5 {
						errMsg := "400 Bad Request"
						log.Error().Str("name", name).Msg(errMsg)
						fmt.Fprintln(os.Stderr, errMsg)
						// Log the command execution details before exiting.
						telemetry.LogCommandExecution("greet", args, flags, "", fmt.Errorf(errMsg))
						os.Exit(1)
					} else {
						message = "HE--<BLERGH>"
						log.Warn().Str("name", name).Msg("Simulated cough output")
						// Log the command execution details with the cough output.
						telemetry.LogCommandExecution("greet", args, flags, message, nil)
						fmt.Println(message)
						return
					}
				}
			}

			// Log a successful command execution.
			telemetry.LogCommandExecution("greet", args, flags, message, nil)
			fmt.Println(message)

			// log the command's duration.
			duration := time.Since(start)
			log.Info().Dur("duration", duration).Msg("Command completed")
		},
	}

	greetCmd.Flags().StringVar(&name, "name", "", "Name to greet (required)")
	greetCmd.MarkFlagRequired("name")
	greetCmd.Flags().BoolVar(&shout, "shout", false, "Shout the greeting (all caps)")

	var rootCmd = &cobra.Command{Use: "mycli"}
	rootCmd.AddCommand(greetCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Msg("Command execution failed")
		fmt.Println(err)
		os.Exit(1)
	}
}

// randFloat returns a random float64 between 0.0 and 1.0.
func randFloat() float64 {
	return float64(time.Now().UnixNano()%1000) / 1000.0
}
