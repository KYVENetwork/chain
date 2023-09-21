package keeper_test

import (
	"testing"

	testkeeper "github.com/KYVENetwork/chain/testutil/keeper"
	"github.com/KYVENetwork/chain/x/funders/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.FundersKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
