package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/darrenparkinson/wtr/internal/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RetrieveCmd is the entry point for the Refresh command
func RetrieveCmd() *cobra.Command {
	var command = &cobra.Command{
		Use:          "retrieve",
		Short:        "retrieve an initial token using parameters in config file",
		Example:      `  wtr retrieve`,
		SilenceUsage: false,
	}

	command.Flags().BoolP("output", "o", false, "output token details to console")
	command.Flags().BoolP("json", "j", false, "output token details as json to console")
	command.Flags().IntP("timeout", "t", 60, "timeout in seconds to wait for response")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		// Flags
		verbose, _ := cmd.Flags().GetBool("debug")
		output, _ := cmd.Flags().GetBool("output")
		jsonf, _ := cmd.Flags().GetBool("json")
		timeout, _ := cmd.Flags().GetInt("timeout")
		// Config
		clientID := viper.GetString("clientid")
		secret := viper.GetString("secret")
		scopes := viper.GetString("scopes")
		redirectPort := viper.GetString("redirectPort")

		if clientID == "" || secret == "" {
			log.Fatal(errors.New("clientid and secret required"))
		}

		log.SetOutput(ioutil.Discard)
		if verbose {
			log.SetOutput(os.Stderr)
		}

		webex := auth.Webex{
			ClientID:     clientID,
			ClientSecret: secret,
			Scopes:       strings.Split(scopes, " "),
			RedirectPort: redirectPort,
		}
		log.Printf("retrieving access token: timeout %ds", timeout)
		t, err := webex.GetAccessToken(timeout)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("token retrieved: expires %s", t.Expires())
		if output {
			fmt.Println(t)
		}
		if jsonf {
			res, _ := json.MarshalIndent(t, "", "  ")
			fmt.Println(string(res))
		}
		err = t.Save()
		return err
	}
	return command
}
