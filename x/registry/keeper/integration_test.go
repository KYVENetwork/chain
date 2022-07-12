package keeper_test

import (
	"fmt"
	"github.com/KYVENetwork/chain/x/registry"
	"github.com/KYVENetwork/chain/x/registry/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

const ALICE_ADDR = "cosmos1jq304cthpx0lwhpqzrdjrcza559ukyy347ju8f"
const BOB_ADDR = "cosmos1hvg7zsnrj6h29q9ss577mhrxa04rn94hfvl2ry"
const KYVE = uint64(1_000_000_000)

var DUMMY_ACCOUNTS []string

func runTxWithResult(msg sdk.Msg) (*sdk.Result, error) {
	cachedCtx, commit := s.ctx.CacheContext()
	resp, err := registry.NewHandler(s.app.RegistryKeeper)(cachedCtx, msg)
	if err == nil {
		commit()
		return resp, nil
	}
	return nil, err
}

func runTx(msg sdk.Msg) (success bool) {
	cachedCtx, commit := s.ctx.CacheContext()
	_, err := registry.NewHandler(s.app.RegistryKeeper)(cachedCtx, msg)
	if err == nil {
		commit()
		return true
	}
	return false
}

func runTxSuccess(t *testing.T, msg sdk.Msg) {
	success := runTx(msg)
	require.True(t, success)
}

func mint(address string, amount uint64) error {
	coins := sdk.NewCoins(sdk.NewInt64Coin("tkyve", int64(amount)))
	err := s.app.BankKeeper.MintCoins(s.ctx, types.ModuleName, coins)
	if err != nil {
		return err
	}

	s.Commit()

	sender, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, types.ModuleName, sender, coins)
	if err != nil {
		return err
	}
	return nil
}

func initDummyAccounts() {
	DUMMY_ACCOUNTS = make([]string, 50)
	rand.Seed(1)
	for i := 0; i < 50; i++ {
		byteAddr := make([]byte, 20)
		for k := 0; k < 20; k++ {
			byteAddr[k] = byte(rand.Int())
		}
		dummy, _ := sdk.Bech32ifyAddressBytes("cosmos", byteAddr)
		DUMMY_ACCOUNTS[i] = dummy
		mint(dummy, 1000*KYVE)
	}
}

func TestBlocks(t *testing.T) {
	s = new(KeeperTestSuite)
	s.SetupTest()

	initDummyAccounts()

	currentTime := s.ctx.BlockTime().Unix()
	s.CommitAfter(time.Second * 60)
	require.Equal(t, s.ctx.BlockTime().Unix(), currentTime+60)

	s.CommitAfter(time.Second * 60)
	require.Equal(t, s.ctx.BlockTime().Unix(), currentTime+2*60)

}

func TestFunding(t *testing.T) {
	pool := types.Pool{
		Creator:        govtypes.ModuleName,
		Name:           "Moontest",
		Runtime:        "@kyve/evm",
		Logo:           "9FJDam56yBbmvn8rlamEucATH5UcYqSBw468rlCXn8E",
		Config:         "{\"rpc\":\"https://rpc.api.moonbeam.network\",\"github\":\"https://github.com/KYVENetwork/evm\"}",
		UploadInterval: 60,
		OperatingCost:  100,
		BundleProposal: &types.BundleProposal{},
		MaxBundleSize:  100,
		Protocol: &types.Protocol{
			Version:     "1.3.0",
			LastUpgrade: uint64(s.ctx.BlockTime().Unix()),
			Binaries:    "{\"macos\":\"https://github.com/kyve-org/evm/releases/download/v1.0.5/kyve-evm-macos.zip\"}",
		},
		UpgradePlan: &types.UpgradePlan{},
		StartKey:    "0",
		Status:      types.POOL_STATUS_NOT_ENOUGH_VALIDATORS,
		MinStake:    0,
	}

	s.app.RegistryKeeper.AppendPool(s.ctx, pool)
	s.Commit()

	err := mint(ALICE_ADDR, 1000*KYVE)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	err = mint(BOB_ADDR, 1000*KYVE)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

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

func TestCommissionChange(t *testing.T) {

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

func TestStaking(t *testing.T) {
	UploadTimeout := s.app.RegistryKeeper.UploadTimeout(s.ctx)
	UnstakingTime := s.app.RegistryKeeper.UnbondingStakingTime(s.ctx)
	UndelegationTime := s.app.RegistryKeeper.UnbondingDelegationTime(s.ctx)
	_ = UndelegationTime

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

	alicaStaker, _ := s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	require.Equal(t, types.STAKER_STATUS_INACTIVE, alicaStaker.Status)

	bobStaker, _ := s.app.RegistryKeeper.GetStaker(s.ctx, BOB_ADDR, 0)
	require.Equal(t, types.STAKER_STATUS_ACTIVE, bobStaker.Status)

	// Reactivate Staker
	runTxSuccess(t, &types.MsgReactivateStaker{
		Creator: ALICE_ADDR,
		PoolId:  0,
	})

	pool, _ = s.app.RegistryKeeper.GetPool(s.ctx, 0)
	require.Len(t, pool.Stakers, 2)
	require.Len(t, pool.InactiveStakers, 0)

	alicaStaker, _ = s.app.RegistryKeeper.GetStaker(s.ctx, ALICE_ADDR, 0)
	require.Equal(t, types.STAKER_STATUS_ACTIVE, alicaStaker.Status)

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
