package main

import (
	"fmt"
	"os"

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

	cmdRetrieve := cmd.RetrieveCmd()

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

	rootCmd.AddCommand(cmdRetrieve)

	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	wtrLogo := aec.RedF.Apply(wtrFigletStr)
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\nv{{.Version}}\n", wtrLogo))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

const wtrFigletStr = `          _        
__      _| |_ _ __ 
\ \ /\ / / __| '__|
 \ V  V /| |_| |   
  \_/\_/  \__|_|                                                                                                                     
`
