package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/irisnet/irishub/modules/asset/01-token/internal/types"
)

// QuerierToken return the definition of token by symbol
func QuerierToken(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokenParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
	}
	return queryToken(ctx, keeper, params.Symbol)
}

// QuerierTokens return the token list by symbol or owner
func QuerierTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokensParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
	}

	if len(params.Symbol) > 0 {
		return queryToken(ctx, keeper, params.Symbol)
	}

	if len(params.Owner) > 0 {
		return queryTokensByOwner(ctx, keeper, params.Owner)
	}

	return queryAllTokens(ctx, keeper)
}

// QuerierFees return the fee of issue or mint by symbol
func QuerierFees(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryTokenFeesParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
	}

	issueFee := GetTokenIssueFee(ctx, keeper, params.Symbol)
	mintFee := GetTokenMintFee(ctx, keeper, params.Symbol)

	fees := types.TokenFeesOutput{
		Exist:    keeper.HasTokenSymbol(ctx, params.Symbol),
		IssueFee: issueFee,
		MintFee:  mintFee,
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, fees)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// QuerierParameters return the system param of the token module
func QuerierParameters(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	params := keeper.GetParamSet(ctx)

	res, err := codec.MarshalJSONIndent(keeper.cdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return res, nil
}

func queryToken(ctx sdk.Context, keeper Keeper, symbol string) ([]byte, sdk.Error) {
	token, found := keeper.GetToken(ctx, symbol)
	if !found {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("token %s does not exist", symbol))
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, types.Tokens{token})
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func queryAllTokens(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	tokens := keeper.GetAllTokens(ctx)
	bz, err := codec.MarshalJSONIndent(keeper.cdc, tokens)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func queryTokensByOwner(ctx sdk.Context, keeper Keeper, ownerStr string) ([]byte, sdk.Error) {
	owner, err := sdk.AccAddressFromBech32(ownerStr)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("failed to parse params: %s", err))
	}
	// Query from db
	iter := keeper.GetTokens(ctx, owner, "")
	defer iter.Close()

	var tokens types.Tokens
	for ; iter.Valid(); iter.Next() {
		var symbol string
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &symbol)
		token, found := keeper.GetToken(ctx, symbol)
		if !found {
			continue
		}

		tokens = append(tokens, token)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, tokens)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
