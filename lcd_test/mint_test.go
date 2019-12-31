package lcdtest

import (
	"net/http"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/irisnet/irishub/modules/mint"
)

func TestMintingQueries(t *testing.T) {
	kb, err := keys.NewKeyringFromDir(InitClientHome(""), nil)
	require.NoError(t, err)
	addr, _, err := CreateAddr(name1, kb)
	require.NoError(t, err)
	cleanup, _, _, port, err := InitializeLCD(1, []sdk.AccAddress{addr}, true)
	require.NoError(t, err)
	defer cleanup()

	res, body := Request(t, port, "GET", "/mint/parameters", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	var params mint.Params
	require.NoError(t, cdc.UnmarshalJSON(extractResultFromResponse(t, []byte(body)), &params))
}
