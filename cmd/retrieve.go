package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func RetrieveCmd() *cobra.Command {
	var command = &cobra.Command{
		Use:          "retrieve",
		Short:        "Retrieve an initial token using parameters in config file",
		Example:      `  wtr retrieve`,
		SilenceUsage: false,
	}
	command.RunE = func(cmd *cobra.Command, args []string) error {
		log.Println("retrieving token")
		return nil

	}
	return command
}
