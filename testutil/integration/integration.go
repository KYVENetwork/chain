package integration

import (
	mrand "math/rand"
	"strconv"
	"time"

	"github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/math"
	"github.com/KYVENetwork/chain/util"

	abci "github.com/cometbft/cometbft/abci/types"

	globalTypes "github.com/KYVENetwork/chain/x/global/types"

	"github.com/KYVENetwork/chain/app"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	ALICE   = "kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd"
	BOB     = "kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq"
	CHARLIE = "kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s"
	DAVID   = "kyve1jxa7kp37jlm8hzwgc5qprquv9k7vawq79qhctt"

	STAKER_0         = "kyve1htgfatqevuvfzvl0sxp97ywteqhg5leha9emf4"
	POOL_ADDRESS_0_A = "kyve1qnf86dkvvtpdukx30r3vajav7rdq8snktm90hm"
	POOL_ADDRESS_0_B = "kyve10t8gnqjnem7tsu09erzswj3zm8599lsnex79rz"
	POOL_ADDRESS_0_C = "kyve13ztkcm2pket6mrmxj8rrmwc6supw7aqakg3uu3"

	STAKER_1         = "kyve1gnr35rwn8rmflnlzs6nn5hhkmzzkxg9ap8xepw"
	POOL_ADDRESS_1_A = "kyve1hpjgzljglmv00nstk3jvcw0zzq94nu0cuxv5ga"
	POOL_ADDRESS_1_B = "kyve14runw9qkltpz2mcx3gsfmlqyyvdzkt3rq3w6fm"
	POOL_ADDRESS_1_C = "kyve15w9m7zpq9ctsxsveqaqkp4uuvw98z5vct6s9g9"

	STAKER_2         = "kyve1xsemlxghgvusumhqzm2ztjw7dz9krvu3de54e2"
	POOL_ADDRESS_2_A = "kyve1u0870dkae6ql63hxvy9y7g65c0y8csfh8allzl"
	POOL_ADDRESS_2_B = "kyve16g3utghkvvlz53jk0fq96zwrhxmqfu36ue965q"
	POOL_ADDRESS_2_C = "kyve18gjtzsn6jme3qsczj9q7wefymlkfu7ngyq5f9c"

	STAKER_3         = "kyve1ca7rzyrxfpdm7j8jgccq4rduuf4sxpq0dhmwm4"
	POOL_ADDRESS_3_A = "kyve1d2clkfrw0r99ctgmkjvluzn6xm98yls06mnxv8"
	POOL_ADDRESS_3_B = "kyve1f36cvde6jnygcrz2yas4acp0akn9cw7vp5ze0w"
	POOL_ADDRESS_3_C = "kyve1gcnd8gya2ysfur6d6z4wpl9z54zadg7qzk8uyc"

	// To avoid giving burner permissions to a module for the tests
	BURNER = "kyve1ld23ktfwc9zstaq8aanwkkj8cf0ru6adtz59y5"
)

var (
	DUMMY    []string
	VALDUMMY []string
)

const (
	KYVE   = uint64(1_000_000)
	T_KYVE = int64(KYVE)
)

var (
	KYVE_DENOM = globalTypes.Denom
	A_DENOM    = "acoin"
	B_DENOM    = "bcoin"
	C_DENOM    = "ccoin"
)

func KYVECoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(KYVE_DENOM, amount)
}

func KYVECoins(amount int64) sdk.Coins {
	return sdk.NewCoins(KYVECoin(amount))
}

func ACoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(A_DENOM, amount)
}

func ACoins(amount int64) sdk.Coins {
	return sdk.NewCoins(ACoin(amount))
}

func BCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(B_DENOM, amount)
}

func BCoins(amount int64) sdk.Coins {
	return sdk.NewCoins(BCoin(amount))
}

func CCoin(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(C_DENOM, amount)
}

func CCoins(amount int64) sdk.Coins {
	return sdk.NewCoins(CCoin(amount))
}

func NewCleanChain() *KeeperTestSuite {
	return NewCleanChainAtTime(time.Now().Unix())
}

func NewCleanChainAtTime(startTime int64) *KeeperTestSuite {
	s := KeeperTestSuite{}
	s.SetupApp(startTime)
	s.initDummyAccounts()
	return &s
}

func (suite *KeeperTestSuite) initDummyAccounts() {
	_ = suite.MintBaseCoins(ALICE, 1000*KYVE)
	_ = suite.MintBaseCoins(BOB, 1000*KYVE)
	_ = suite.MintBaseCoins(CHARLIE, 1000*KYVE)
	_ = suite.MintBaseCoins(DAVID, 1000*KYVE)

	_ = suite.MintBaseCoins(STAKER_0, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_0_A, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_0_B, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_0_C, 1000*KYVE)

	_ = suite.MintBaseCoins(STAKER_1, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_1_A, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_1_B, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_1_C, 1000*KYVE)

	_ = suite.MintBaseCoins(STAKER_2, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_2_A, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_2_B, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_2_C, 1000*KYVE)

	_ = suite.MintBaseCoins(STAKER_3, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_3_A, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_3_B, 1000*KYVE)
	_ = suite.MintBaseCoins(POOL_ADDRESS_3_C, 1000*KYVE)

	DUMMY = make([]string, 50)

	for i := 0; i < 50; i++ {
		byteAddr := make([]byte, 20)
		for k := 0; k < 20; k++ {
			randomSource := mrand.New(mrand.NewSource(int64(i + k)))
			byteAddr[k] = byte(randomSource.Int())
		}
		dummy, _ := sdk.Bech32ifyAddressBytes("kyve", byteAddr)
		DUMMY[i] = dummy
		_ = suite.MintBaseCoins(dummy, 1000*KYVE)
	}

	VALDUMMY = make([]string, 50)
	for i := 0; i < 50; i++ {
		byteAddr := make([]byte, 20)
		for k := 0; k < 20; k++ {
			randomSource := mrand.New(mrand.NewSource(int64(i + k + 100)))
			byteAddr[k] = byte(randomSource.Int())
		}
		dummy, _ := sdk.Bech32ifyAddressBytes("kyve", byteAddr)
		VALDUMMY[i] = dummy
		_ = suite.MintBaseCoins(dummy, 1000*KYVE)
	}
}

func (suite *KeeperTestSuite) MintBaseCoins(address string, amount uint64) error {
	return suite.MintCoins(address, sdk.NewCoins(
		// mint coins ukyve, A, B, C
		sdk.NewInt64Coin(KYVE_DENOM, int64(amount)),
		sdk.NewInt64Coin(A_DENOM, int64(amount)),
		sdk.NewInt64Coin(B_DENOM, int64(amount)),
		sdk.NewInt64Coin(C_DENOM, int64(amount)),
	))
}

func (suite *KeeperTestSuite) MintCoins(address string, coins sdk.Coins) error {
	// mint coins ukyve, A, B, C
	err := suite.app.BankKeeper.MintCoins(suite.ctx, mintTypes.ModuleName, coins)
	if err != nil {
		return err
	}

	suite.Commit()

	receiver, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return err
	}

	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, mintTypes.ModuleName, receiver, coins)
	if err != nil {
		return err
	}

	return nil
}

func (suite *KeeperTestSuite) MintDenomToModule(moduleAddress string, amount uint64, denom string) error {
	coins := sdk.NewCoins(sdk.NewInt64Coin(denom, int64(amount)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, mintTypes.ModuleName, coins)
	if err != nil {
		return err
	}

	suite.Commit()

	err = suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, mintTypes.ModuleName, moduleAddress, coins)
	if err != nil {
		return err
	}

	return nil
}

func (suite *KeeperTestSuite) App() *app.App {
	return suite.app
}

func (suite *KeeperTestSuite) Ctx() sdk.Context {
	return suite.ctx
}

func (suite *KeeperTestSuite) SetCtx(ctx sdk.Context) {
	suite.ctx = ctx
}

type TestValidatorAddress struct {
	Moniker string

	PrivateKey *ed25519.PrivKey

	Address        string
	ValAddress     string
	AccAddress     sdk.AccAddress
	ConsAccAddress sdk.ConsAddress
	ConsAddress    string

	PoolAccount [10]string
}

func (suite *KeeperTestSuite) CreateValidatorFromFullAddress(address TestValidatorAddress, kyveStake int64) {
	valAddress := util.MustValaddressFromOperatorAddress(address.Address)

	msg, _ := stakingtypes.NewMsgCreateValidator(
		valAddress,
		address.PrivateKey.PubKey(),
		sdk.NewInt64Coin(globalTypes.Denom, kyveStake),
		stakingtypes.Description{Moniker: address.Moniker},
		stakingtypes.NewCommissionRates(math.LegacyMustNewDecFromStr("0.1"), math.LegacyMustNewDecFromStr("1"), math.LegacyMustNewDecFromStr("1")),
		math.NewInt(1),
	)

	_, err := suite.RunTx(msg)
	if err != nil {
		panic(err)
	}

	suite.Commit()
}

func (suite *KeeperTestSuite) CreateValidatorWithoutCommit(address, moniker string, kyveStake int64) {
	valAddress := util.MustValaddressFromOperatorAddress(address)

	msg, _ := stakingtypes.NewMsgCreateValidator(
		valAddress,
		ed25519.GenPrivKeyFromSecret([]byte(valAddress)).PubKey(),
		sdk.NewInt64Coin(globalTypes.Denom, kyveStake),
		stakingtypes.Description{Moniker: moniker},
		stakingtypes.NewCommissionRates(math.LegacyMustNewDecFromStr("0.1"), math.LegacyMustNewDecFromStr("1"), math.LegacyMustNewDecFromStr("1")),
		math.NewInt(1),
	)

	_, err := suite.RunTx(msg)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) SelfDelegateValidator(address string, amount uint64) {
	valAddress := util.MustValaddressFromOperatorAddress(address)

	msg := stakingtypes.NewMsgDelegate(
		address,
		valAddress,
		sdk.NewInt64Coin(globalTypes.Denom, int64(amount)),
	)

	_, err := suite.RunTx(msg)
	if err != nil {
		panic(err)
	}

	suite.Commit()
}

func (suite *KeeperTestSuite) SelfUndelegateValidator(address string, amount uint64) {
	valAddress := util.MustValaddressFromOperatorAddress(address)

	msg := stakingtypes.NewMsgUndelegate(
		address,
		valAddress,
		sdk.NewInt64Coin(globalTypes.Denom, int64(amount)),
	)

	_, err := suite.RunTx(msg)
	if err != nil {
		panic(err)
	}

	suite.Commit()
}

func (suite *KeeperTestSuite) CreateNewValidator(moniker string, kyveStake uint64) TestValidatorAddress {
	a := GenerateTestValidatorAddress(moniker)
	_ = suite.MintBaseCoins(a.Address, 10*kyveStake)
	msg, _ := stakingtypes.NewMsgCreateValidator(
		util.MustValaddressFromOperatorAddress(a.Address),
		a.PrivateKey.PubKey(),
		sdk.NewInt64Coin(globalTypes.Denom, int64(kyveStake)),
		stakingtypes.Description{Moniker: moniker},
		stakingtypes.NewCommissionRates(math.LegacyMustNewDecFromStr("0.1"), math.LegacyMustNewDecFromStr("1"), math.LegacyMustNewDecFromStr("1")),
		math.NewInt(1),
	)
	_, err := suite.RunTx(msg)
	if err != nil {
		panic(err)
	}
	suite.Commit()
	return a
}

func (suite *KeeperTestSuite) CreateValidator(address, moniker string, kyveStake int64) {
	suite.CreateValidatorWithoutCommit(address, moniker, kyveStake)
	suite.Commit()
}

func (suite *KeeperTestSuite) CreateZeroDelegationValidator(address, name string) {
	// create zero delegation validator by overwriting the min-self-delegation with an invalid value
	// it is fine for the test
	suite.CreateValidator(address, name, int64(100*KYVE))
	val, _ := sdk.ValAddressFromBech32(util.MustValaddressFromOperatorAddress(address))
	validator, _ := suite.App().StakingKeeper.GetValidator(suite.Ctx(), val)
	validator.MinSelfDelegation = math.ZeroInt()
	_ = suite.App().StakingKeeper.SetValidator(suite.Ctx(), validator)
	suite.RunTxSuccess(stakingtypes.NewMsgUndelegate(
		address,
		util.MustValaddressFromOperatorAddress(address),
		sdk.NewInt64Coin("tkyve", int64(100*KYVE)),
	))
}

func (suite *KeeperTestSuite) SetDelegationToZero(address string) {
	val, _ := sdk.ValAddressFromBech32(util.MustValaddressFromOperatorAddress(address))
	validator, _ := suite.App().StakingKeeper.GetValidator(suite.Ctx(), val)
	validator.MinSelfDelegation = math.ZeroInt()
	_ = suite.App().StakingKeeper.SetValidator(suite.Ctx(), validator)
	suite.RunTxSuccess(stakingtypes.NewMsgUndelegate(
		address,
		util.MustValaddressFromOperatorAddress(address),
		sdk.NewInt64Coin("tkyve", validator.BondedTokens().Int64()),
	))
}

func (suite *KeeperTestSuite) Commit() {
	suite.CommitAfter(time.Second * 0)
}

func (suite *KeeperTestSuite) CommitAfterSeconds(seconds uint64) {
	suite.CommitAfter(time.Second * time.Duration(seconds))
}

func (suite *KeeperTestSuite) CommitAfter(t time.Duration) {
	header := suite.ctx.BlockHeader()
	header.Time = header.Time.Add(t)

	_, err := suite.app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height: header.Height,
		Time:   header.Time,
		DecidedLastCommit: abci.CommitInfo{
			Round: 0,
			Votes: suite.VoteInfos,
		},
	})
	if err != nil {
		panic(err)
	}
	_, err = suite.app.Commit()
	if err != nil {
		panic(err)
	}

	header.Height += 1

	suite.ctx = suite.app.BaseApp.NewUncachedContext(false, header)
}

func (suite *KeeperTestSuite) WaitSeconds(seconds uint64) {
	suite.Wait(time.Second * time.Duration(seconds))
}

func (suite *KeeperTestSuite) Wait(t time.Duration) {
	header := suite.ctx.BlockHeader()
	header.Time = header.Time.Add(t)

	suite.ctx = suite.ctx.WithBlockTime(suite.ctx.BlockTime().Add(t))
}

func GenerateTestValidatorAddress(moniker string) TestValidatorAddress {
	a := TestValidatorAddress{}
	a.Moniker = moniker
	a.PrivateKey = ed25519.GenPrivKeyFromSecret([]byte(moniker))

	a.AccAddress = sdk.AccAddress(a.PrivateKey.PubKey().Address())
	bech32Address, _ := sdk.Bech32ifyAddressBytes("kyve", a.AccAddress)
	a.Address = bech32Address
	a.ValAddress = util.MustValaddressFromOperatorAddress(a.Address)

	a.ConsAccAddress = sdk.ConsAddress(a.PrivateKey.PubKey().Address())
	bech32ConsAddress, _ := sdk.Bech32ifyAddressBytes("kyvevalcons", a.AccAddress)
	a.ConsAddress = bech32ConsAddress

	for i := 0; i < 10; i++ {
		poolAddress := ed25519.GenPrivKeyFromSecret([]byte("pool_address" + moniker + strconv.Itoa(i))).PubKey().Address()
		address, _ := sdk.Bech32ifyAddressBytes("kyve", sdk.AccAddress(poolAddress))
		a.PoolAccount[i] = address
	}

	return a
}

func (suite *KeeperTestSuite) ResetAbciVotes() {
	suite.VoteInfos = nil
}

func (suite *KeeperTestSuite) AddAbciCommitVotes(addresses ...sdk.ConsAddress) {
	suite.addAbciVotes(2, addresses...)
}

func (suite *KeeperTestSuite) AddAbciAbsentVote(addresses ...sdk.ConsAddress) {
	suite.addAbciVotes(1, addresses...)
}

func (suite *KeeperTestSuite) addAbciVotes(blogFlagId int32, addresses ...sdk.ConsAddress) {
	suite.VoteInfos = make([]abci.VoteInfo, 0)
	for _, address := range addresses {
		suite.VoteInfos = []abci.VoteInfo{{
			Validator: abci.Validator{
				Address: address,
				Power:   1,
			},
			BlockIdFlag: types.BlockIDFlag(blogFlagId),
		}}
	}
}
