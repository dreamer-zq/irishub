package cmd

import (
	"github.com/irisnet/irishub/app/v1/auth"
	v2 "github.com/irisnet/irishub/app/v2"
	"github.com/spf13/cobra"
)

// CheckAccountCmd implements the default command for a tx query.
func CheckAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "account [account.csv]",
		Short:   "Matches this txhash over all committed blocks",
		Example: "iriscli tendermint tx <transaction hash>",
		//Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, app := NewApp()

			protocol := app.Engine.GetCurrentProtocol()
			p := protocol.(*v2.ProtocolV2)

			var srcAccMap = make(AccMap)
			p.AccountMapper.IterateAccounts(ctx, func(acc auth.Account) (stop bool) {
				srcAccMap[acc.GetAddress().String()] = AccValue{
					Balance: acc.GetCoins().String(),
				}
				return false
			})
			return nil
		},
	}
	return cmd
}

type AccValue struct {
	Balance string
	Type    string
	Module  string
}

type AccMap map[string]AccValue
