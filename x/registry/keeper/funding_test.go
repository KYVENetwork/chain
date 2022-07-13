package keeper_test

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFunding(t *testing.T) {
	createGenesis(t)
	testFunding(t)
}

func testFunding(t *testing.T) {

	fundPool0 := runTx(&types.MsgFundPool{
		Creator: ALICE_ADDR,
		Id:      0,
		Amount:  99 * KYVE,
	})

	_, foundFunder0 := s.app.RegistryKeeper.GetFunder(s.ctx, ALICE_ADDR, 0)
	require.True(t, fundPool0)
	require.True(t, foundFunder0)

	fundPool1 := runTx(&types.MsgFundPool{
		Creator: ALICE_ADDR,
		Id:      1,
		Amount:  0,
	})

	_, foundFunder1 := s.app.RegistryKeeper.GetFunder(s.ctx, ALICE_ADDR, 1)
	require.False(t, fundPool1)
	require.False(t, foundFunder1)
}
