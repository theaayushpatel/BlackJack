package cli_client

import (
	"fmt"
	"github.com/theaayushpatel/BlackJack/common_web"
	"os"

	"github.com/spf13/cobra"
)

var (
	StrRemoteHost string
	StrRemotePort string

	BoolJson bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "BlackJack_cli",
	Short: "BlackJack (remote) CLI configuration",
	Long: `The CLI client tool could be used to configure BlackJack A.L.O.A.
from the command line. The tool relies on RPC so it could be used 
remotely.

` + "Version: " + common_web.VERSION,
}

func GenBashComplete() {
	target := "/etc/bash_completion.d/BlackJack.sh"
	if _, err := os.Stat(target); os.IsNotExist(err) {
		rootCmd.GenBashCompletionFile(target)
	}
}

func Execute() {
	// ToDo: this should be changed to a dedicated command which is sourced in in .bashrc to assure updates on start of bash
	GenBashComplete()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&StrRemoteHost, "host", "localhost", "The host with the listening BlackJack RPC server")
	rootCmd.PersistentFlags().StringVar(&StrRemotePort, "port", "50051", "The port on which the BlackJack RPC server is listening")

	/*
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	*/
}
