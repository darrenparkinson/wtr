package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/darrenparkinson/wtr/internal/auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RefreshCmd is the entry point for the Refresh command
func RefreshCmd() *cobra.Command {
	var command = &cobra.Command{
		Use:          "refresh",
		Short:        "refresh an existing token using details saved to the config file",
		Example:      `  wtr refresh`,
		SilenceUsage: false,
	}

	command.Flags().BoolP("output", "o", false, "output token details to console")
	command.Flags().BoolP("json", "j", false, "output token details as json to console")

	command.RunE = func(cmd *cobra.Command, args []string) error {
		// Flags
		verbose, _ := cmd.Flags().GetBool("debug")
		output, _ := cmd.Flags().GetBool("output")
		jsonf, _ := cmd.Flags().GetBool("json")
		// Config
		clientID := viper.GetString("clientid")
		secret := viper.GetString("secret")
		expiry := viper.GetInt64("expiration")
		refreshToken := viper.GetString("refresh_token")

		if refreshToken == "" || expiry == 0 {
			log.Fatal(errors.New("refreshToken and expiry required in configuration file: run retrieve first"))
		}

		if clientID == "" || secret == "" || refreshToken == "" {
			log.Fatal(errors.New("clientid, secret and refreshToken required in configuration file: run retrieve first"))
		}

		log.SetOutput(ioutil.Discard)
		if verbose {
			log.SetOutput(os.Stderr)
		}

		log.Println("refreshing access token")
		t, err := auth.RefreshToken(clientID, secret, refreshToken)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("token refreshed: expires %s", t.Expires())
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
