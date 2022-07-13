package keeper_test

import (
	"github.com/KYVENetwork/chain/x/registry/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRedelegation(t *testing.T) {
	createGenesis(t)
	testRedelegation(t)
}

func testRedelegation(t *testing.T) {

	runTxSuccess(t, &types.MsgStakePool{
		Creator: BOB_ADDR,
		Id:      0,
		Amount:  50,
	})

	runTxSuccess(t, &types.MsgStakePool{
		Creator: ALICE_ADDR,
		Id:      0,
		Amount:  50,
	})

	pool, _ := s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.Len(t, pool.Stakers, 2)

	// Delegate
	runTxSuccess(t, &types.MsgDelegatePool{
		Creator: DUMMY_ACCOUNTS[0],
		Id:      0,
		Staker:  BOB_ADDR,
		Amount:  100 * KYVE,
	})
	s.Commit()

	delegator, found := s.app.RegistryKeeper.GetDelegator(s.ctx, 0, BOB_ADDR, DUMMY_ACCOUNTS[0])

	require.True(t, found)
	require.Equal(t, 100*KYVE, delegator.DelegationAmount)

	// Redelegate to staker that does not exist; even though it does not make sense it is allowed
	redelegateRes := runTx(&types.MsgRedelegatePool{
		Creator:    DUMMY_ACCOUNTS[0],
		FromPoolId: 0,
		FromStaker: BOB_ADDR,
		ToPoolId:   0,
		ToStaker:   DUMMY_ACCOUNTS[2],
		Amount:     50 * KYVE,
	})
	require.True(t, redelegateRes)

	// DO invalid redelegation from account without delegation
	res := runTx(&types.MsgRedelegatePool{
		Creator:    DUMMY_ACCOUNTS[0],
		FromPoolId: 0,
		FromStaker: DUMMY_ACCOUNTS[0],
		ToPoolId:   0,
		ToStaker:   ALICE_ADDR,
		Amount:     50 * KYVE,
	})
	require.False(t, res)

	s.Commit()
	// cannot commit within same blocktime
	res = runTx(&types.MsgRedelegatePool{
		Creator:    DUMMY_ACCOUNTS[0],
		FromPoolId: 0,
		FromStaker: BOB_ADDR,
		ToPoolId:   0,
		ToStaker:   ALICE_ADDR,
		Amount:     5 * KYVE,
	})
	require.False(t, res)

	s.CommitAfterSeconds(60*60*24 - 1)

	// Fill up queue
	for i := 0; i < 4; i++ {
		runTxSuccess(t, &types.MsgRedelegatePool{
			Creator:    DUMMY_ACCOUNTS[0],
			FromPoolId: 0,
			FromStaker: BOB_ADDR,
			ToPoolId:   0,
			ToStaker:   ALICE_ADDR,
			Amount:     5 * KYVE,
		})
		s.CommitAfterSeconds(60*60*24 - 1)
	}

	s.CommitAfterSeconds(1)

	// All delegation spells taken
	res = runTx(&types.MsgRedelegatePool{
		Creator:    DUMMY_ACCOUNTS[0],
		FromPoolId: 0,
		FromStaker: BOB_ADDR,
		ToPoolId:   0,
		ToStaker:   ALICE_ADDR,
		Amount:     5 * KYVE,
	})
	require.False(t, res)

	// Wait for first entry to expire

	s.CommitAfterSeconds(10)
	runTxSuccess(t, &types.MsgRedelegatePool{
		Creator:    DUMMY_ACCOUNTS[0],
		FromPoolId: 0,
		FromStaker: BOB_ADDR,
		ToPoolId:   0,
		ToStaker:   ALICE_ADDR,
		Amount:     5 * KYVE,
	})

	s.CommitAfterSeconds(1)
	res = runTx(&types.MsgRedelegatePool{
		Creator:    DUMMY_ACCOUNTS[0],
		FromPoolId: 0,
		FromStaker: BOB_ADDR,
		ToPoolId:   0,
		ToStaker:   ALICE_ADDR,
		Amount:     5 * KYVE,
	})
	require.False(t, res)

}
