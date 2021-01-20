package cmd

import (
	"encoding/csv"
	"os"
	"path/filepath"

	"github.com/irisnet/irishub/app"
	"github.com/irisnet/irishub/server"
	sdk "github.com/irisnet/irishub/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/abci/types"
	dbm "github.com/tendermint/tm-db"
)

func NewCheckerCmd() *cobra.Command {
	//Add state commands
	checkerCmd := &cobra.Command{
		Use:   "check",
		Short: "check some module state",
	}
	checkerCmd.AddCommand(
		CheckAccountCmd(),
	)
	return checkerCmd
}

func NewApp() (sdk.Context, *app.IrisApp) {
	serverCtx := server.NewDefaultContext()
	home := viper.GetString("home")
	db, err := openDB(home)
	if err != nil {
		panic(err)
	}
	app := app.NewIrisApp(serverCtx.Logger, db, serverCtx.Config.Instrumentation, nil)
	ctx := app.BaseApp.NewContext(true, types.Header{})
	return ctx, app
}

func LoadCSV(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if reader == nil {
		panic("NewReader return nil")
	}

	records, err := reader.ReadAll()
	if err == nil {
		panic(err)
	}
	return records
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := dbm.NewGoLevelDB("application", dataDir)
	return db, err
}
