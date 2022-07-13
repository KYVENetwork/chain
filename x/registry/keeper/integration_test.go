package keeper_test

import (
	"fmt"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	testInitIntegration(t)
	testIntegrationFunding(t)
	testCommissionChange(t)
	testStaking(t)
}

func testInitIntegration(t *testing.T) {
	createGenesis(t)
}

func testIntegrationFunding(t *testing.T) {

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

func testCommissionChange(t *testing.T) {

	QueueTime := time.Duration(s.app.RegistryKeeper.CommissionChangeTime(s.ctx)) * time.Second
	fmt.Printf("Queue Time: %d\n", QueueTime)

	// Stake into pool
	stakePool := runTx(&types.MsgStakePool{
		Creator: ALICE_ADDR,
		Id:      0,
		Amount:  10 * KYVE,
	})
	require.True(t, stakePool)
	s.Commit()
	staker, found := s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	require.True(t, found)
	require.Equal(t, types.DefaultCommission, staker.Commission)

	// Submit commission change
	s.Commit()
	changeCommission := runTx(&types.MsgUpdateCommission{
		Creator:    ALICE_ADDR,
		Id:         0,
		Commission: "0.1",
	})
	require.True(t, changeCommission)

	fmt.Printf("state: %v\nentires: %v\n", s.app.RegistryKeeper.GetCommissionChangeQueueState(s.ctx), s.app.RegistryKeeper.GetAllCommissionChangeQueueEntries(s.ctx))

	// Commission should not have changed
	staker, found = s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	require.True(t, found)
	require.Equal(t, types.DefaultCommission, staker.Commission)

	s.CommitAfter(QueueTime) // Wait timeout
	s.Commit()               // Execute endblock after time was upgraded

	// Commission should have changed
	staker, found = s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	require.True(t, found)
	require.Equal(t, "0.1", staker.Commission)

	// =================
	// Add second staker
	// =================

	stakePool2 := runTx(&types.MsgStakePool{
		Creator: BOB_ADDR,
		Id:      0,
		Amount:  10 * KYVE,
	})
	require.True(t, stakePool2)
	s.Commit()
	staker2, found2 := s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	require.True(t, found2)
	// Commission should be default
	require.Equal(t, types.DefaultCommission, staker2.Commission)

	// Create two queue entries
	runTxSuccess(t, &types.MsgUpdateCommission{
		Creator:    ALICE_ADDR,
		Id:         0,
		Commission: "0.2",
	})
	s.CommitAfter(QueueTime / 2)
	runTxSuccess(t, &types.MsgUpdateCommission{
		Creator:    BOB_ADDR,
		Id:         0,
		Commission: "0.5",
	})
	s.CommitAfter(QueueTime / 3)

	fmt.Printf("state: %v\nentires: %v\n", s.app.RegistryKeeper.GetCommissionChangeQueueState(s.ctx), s.app.RegistryKeeper.GetAllCommissionChangeQueueEntries(s.ctx))

	// No change should have been applied yet.
	aliceStaker, found := s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	require.Equal(t, "0.1", aliceStaker.Commission)

	// Overwrite Bob which will reset the timer
	runTxSuccess(t, &types.MsgUpdateCommission{
		Creator:    BOB_ADDR,
		Id:         0,
		Commission: "0.8",
	})

	// Alice should have become active
	s.CommitAfter(QueueTime / 2)
	s.Commit()

	// Alice should now be due, bob should be unchanged

	fmt.Printf("state: %v\nentires: %v\n", s.app.RegistryKeeper.GetCommissionChangeQueueState(s.ctx), s.app.RegistryKeeper.GetAllCommissionChangeQueueEntries(s.ctx))

	aliceStaker, _ = s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	require.Equal(t, "0.2", aliceStaker.Commission)
	bobStaker, _ := s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	require.Equal(t, types.DefaultCommission, bobStaker.Commission)

	s.CommitAfter(QueueTime / 2)
	s.Commit()

	// bob should now have become active
	bobStaker, _ = s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	require.Equal(t, "0.8", bobStaker.Commission)

	fmt.Printf("state: %v\nentires: %v\n", s.app.RegistryKeeper.GetCommissionChangeQueueState(s.ctx), s.app.RegistryKeeper.GetAllCommissionChangeQueueEntries(s.ctx))
}

func testStaking(t *testing.T) {
	UploadTimeout := s.app.RegistryKeeper.UploadTimeout(s.ctx)
	UnstakingTime := s.app.RegistryKeeper.UnbondingStakingTime(s.ctx)
	UndelegationTime := s.app.RegistryKeeper.UnbondingDelegationTime(s.ctx)
	_ = UndelegationTime

	// Unstake everything.
	aliceStaker, _ := s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	runTxSuccess(t, &types.MsgUnstakePool{
		Creator: ALICE_ADDR,
		Id:      0,
		Amount:  aliceStaker.Amount,
	})
	// Unstake everything.
	bobStaker, _ := s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	runTxSuccess(t, &types.MsgUnstakePool{
		Creator: BOB_ADDR,
		Id:      0,
		Amount:  bobStaker.Amount,
	})

	s.CommitAfterSeconds(UnstakingTime + 1)
	s.CommitAfterSeconds(1)
	pool, _ := s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.Len(t, pool.Stakers, 0)
	require.Len(t, pool.InactiveStakers, 0)
	require.Equal(t, uint64(0), pool.TotalStake)
	require.Equal(t, uint64(0), pool.TotalInactiveStake)

	s.CommitAfterSeconds(UploadTimeout + 10)
	pool, _ = s.app.RegistryKeeper.GetPool(s.ctx, 0)
	// Stake into pool
	stakePool := runTx(&types.MsgStakePool{
		Creator: ALICE_ADDR,
		Id:      0,
		Amount:  10 * KYVE,
	})
	require.True(t, stakePool)

	stakePool2 := runTx(&types.MsgStakePool{
		Creator: BOB_ADDR,
		Id:      0,
		Amount:  10 * KYVE,
	})
	require.True(t, stakePool2)

	pool, _ = s.app.RegistryKeeper.GetPool(s.ctx, 0)

	s.CommitAfterSeconds(2)

	runTxSuccess(t, &types.MsgClaimUploaderRole{
		Creator: ALICE_ADDR,
		Id:      0,
	})

	pool, found := s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.True(t, found)

	require.Equal(t, ALICE_ADDR, pool.BundleProposal.NextUploader)

	s.CommitAfterSeconds(UploadTimeout + pool.UploadInterval + 1)
	s.Commit()

	pool, found = s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.True(t, found)
	require.Len(t, pool.Stakers, 1)
	require.Len(t, pool.InactiveStakers, 1)

	aliceStaker, _ = s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	require.Equal(t, types.STAKER_STATUS_INACTIVE, aliceStaker.Status)

	bobStaker, _ = s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	require.Equal(t, types.STAKER_STATUS_ACTIVE, bobStaker.Status)

	// Reactivate Staker
	runTxSuccess(t, &types.MsgReactivateStaker{
		Creator: ALICE_ADDR,
		PoolId:  0,
	})

	pool, _ = s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.Len(t, pool.Stakers, 2)
	require.Len(t, pool.InactiveStakers, 0)

	aliceStaker, _ = s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	require.Equal(t, types.STAKER_STATUS_ACTIVE, aliceStaker.Status)

	bobStaker, _ = s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	require.Equal(t, types.STAKER_STATUS_ACTIVE, bobStaker.Status)

	s.Commit()

	for i := 0; i < 48; i++ {
		dummyStake := runTx(&types.MsgStakePool{
			Creator: DUMMY_ACCOUNTS[i],
			Id:      0,
			Amount:  5 * KYVE,
		})
		require.True(t, dummyStake)
	}

	failedStake := runTx(&types.MsgStakePool{
		Creator: DUMMY_ACCOUNTS[48],
		Id:      0,
		Amount:  2 * KYVE,
	})
	require.False(t, failedStake)
	pool, _ = s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.Len(t, pool.Stakers, 50)
	require.Len(t, pool.InactiveStakers, 0)

	//Make Bob to Unstake a few KYVE

	bobStaker, _ = s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	runTxSuccess(t, &types.MsgUnstakePool{
		Creator: BOB_ADDR,
		Id:      0,
		Amount:  9_800_000_000,
	})

	s.CommitAfterSeconds(10)

	// bobs stake should still be 10
	bobStaker, _ = s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	require.Equal(t, 10*KYVE, bobStaker.Amount)

	// User is currently NextUploader, after this long timeout he will get slashed
	// Simultaneously the user is unstaking all of his tokens.
	s.CommitAfterSeconds(UnstakingTime + 1)
	s.Commit()

	pool, _ = s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.NotContains(t, pool.Stakers, BOB_ADDR)
	require.Len(t, pool.Stakers, 49)
	require.NotContains(t, pool.InactiveStakers, BOB_ADDR)
	require.Len(t, pool.InactiveStakers, 0)

	_, found = s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	require.False(t, found)

	// Check bobs balance.
	acc, _ := sdk.AccAddressFromBech32(BOB_ADDR)
	balance := s.app.BankKeeper.GetBalance(s.ctx, acc, "tkyve")

	// Bobs balacne should be 1000KYVE - 0.02 * 10 KYVE
	require.Equal(t, 1000*KYVE-200_000_000, balance.Amount.Uint64())

	// Overstake lowest staker
	runTxSuccess(t, &types.MsgStakePool{
		Creator: DUMMY_ACCOUNTS[48],
		Id:      0,
		Amount:  1 * KYVE,
	})

	s.Commit()
	// Bob should now be inactive
	pool, _ = s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.Len(t, pool.InactiveStakers, 0)
	require.Len(t, pool.Stakers, 50)
	require.NotContains(t, pool.InactiveStakers, BOB_ADDR)
	require.NotContains(t, pool.Stakers, BOB_ADDR)

	lowestStaker, _ := s.app.RegistryKeeper.GetStaker(s.ctx, pool.LowestStaker, 0)
	require.Equal(t, DUMMY_ACCOUNTS[48], lowestStaker.Account)

	// Kick out lowest staker and check if staker becomes inactive
	runTxSuccess(t, &types.MsgStakePool{
		Creator: DUMMY_ACCOUNTS[49],
		Id:      0,
		Amount:  100 * KYVE,
	})

	s.CommitAfterSeconds(1)
	pool, _ = s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.Len(t, pool.InactiveStakers, 1)
	require.Len(t, pool.Stakers, 50)
	require.Contains(t, pool.InactiveStakers, DUMMY_ACCOUNTS[48])
	require.NotContains(t, pool.Stakers, DUMMY_ACCOUNTS[48])

	d48, _ := s.app.RegistryKeeper.GetStaker(s.ctx, DUMMY_ACCOUNTS[48], 0)
	fmt.Printf("Dummy48: %v\n", d48)
	fmt.Printf("Stakers: %v\n", pool.Stakers)
	fmt.Printf("InactiveStakers: %v\n", pool.InactiveStakers)

	reactiveShouldFail := runTx(&types.MsgReactivateStaker{
		Creator: DUMMY_ACCOUNTS[48],
		PoolId:  0,
	})
	require.False(t, reactiveShouldFail)

	s.Commit()

	// Increase stake
	runTxSuccess(t, &types.MsgStakePool{
		Creator: DUMMY_ACCOUNTS[48],
		Id:      0,
		Amount:  100 * KYVE,
	})

	// Reactivate Staker now
	runTxSuccess(t, &types.MsgReactivateStaker{
		Creator: DUMMY_ACCOUNTS[48],
		PoolId:  0,
	})

	s.CommitAfterSeconds(1)
	pool, _ = s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.Len(t, pool.InactiveStakers, 1)
	require.Len(t, pool.Stakers, 50)
	require.Contains(t, pool.Stakers, DUMMY_ACCOUNTS[48])
	require.NotContains(t, pool.InactiveStakers, DUMMY_ACCOUNTS[48])

}
