package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/darrenparkinson/wtr/cmd"

	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// version is a linker flag set by goreleaser
var version = "0.0.0"

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".wtr-cli")
	viper.SetConfigType("json")

	_ = viper.SafeWriteConfig()
	_ = viper.ReadInConfig()

	viper.SetEnvPrefix("webex")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	cmdRetrieve := cmd.RetrieveCmd()
	cmdRefresh := cmd.RefreshCmd()

	var rootCmd = &cobra.Command{
		Use:     "wtr",
		Short:   "Webex Token Retriever/Refresher",
		Version: version,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(cmdRetrieve)
	rootCmd.AddCommand(cmdRefresh)

	rootCmd.PersistentFlags().BoolP("debug", "d", false, "verbose debug output")

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	wtrLogo := aec.GreenF.Apply(wtrFigletStr)
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\nv{{.Version}}\nhttps://github.com/darrenparkinson/wtr\n", wtrLogo))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const wtrFigletStr = `          _        
__      _| |_ _ __ 2
\ \ /\ / / __| '__|
 \ V  V /| |_| |   
  \_/\_/  \__|_|                                                                                                                     
`
