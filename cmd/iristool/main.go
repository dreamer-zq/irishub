package main

import (
	"os"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/cmd/iristool/cmd"
	debugcmd "github.com/irisnet/irishub/tools/debug"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/tmlibs/cli"
)

func init() {
	sdk.SetNetworkType("mainnet")
	rootCmd.AddCommand(
		debugcmd.RootCmd,
		cmd.NewCheckerCmd(),
	)

	//rootCmd.PersistentFlags().String("home", "", "data home")
	//viper.BindPFlags(rootCmd.PersistentFlags())
}

var rootCmd = &cobra.Command{
	Use:          "iristool",
	Short:        "Iris tool",
	SilenceUsage: true,
}

func main() {
	executor := cli.PrepareBaseCmd(rootCmd, "IRIS", app.DefaultNodeHome)
	if err := executor.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
