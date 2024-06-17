package v1_5

type GovProposalData struct {
	title    string
	summary  string
	proposer string
}

var govProposalsData = map[uint64]GovProposalData{
	1: {
		title:    "Parameter Change: Enabling Validator Creation, Increasing Voting Period",
		summary:  "Currently, it is not possible to create a validator, and the voting period is set to 1 hour. If passed, this parameter change proposal would enable validator creation on the KYVE network and increase the voting period to 1 week for optimal governance participation.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	2: {
		title:    "Parameter Change: Enabling Inflation",
		summary:  "Currently, inflation has been disabled on the KYVE network while we rolled out our Foundation Delegation Program. This proposal, if passed, will enable inflation to be set to our default parameters, aiming to reach an APY of 20%.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	3: {
		title:    "v1.1.0 Software Upgrade",
		summary:  "Recently, while undergoing an independent code review, we found a few non-critical issues. If passed, this software upgrade proposal would implement fixes and improvements to resolve these issues on the KYVE network. Additionally, if passed, token transfers via IBC would be enabled.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	4: {
		title:    "v1.2.0 Software Upgrade",
		summary:  "This software upgrade proposal, if passed, would enable full Ledger support for **all** of KYVE's custom transactions.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	5: {
		title:    "Pool Creation: Cosmos Hub",
		summary:  "After ample testing and observation on Kaon, our official testnet, the Cosmos Hub data pool works smoothly and has successfully activated the testnet protocol layer. That being said, this pool creation proposal, if passed, would enact the creation of the Cosmos Hub data pool on the KYVE network, officially kickstarting the mainnet protocol layer, commonly known as our decentralized data lake and core product.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	6: {
		title:    "Parameter Change: Protocol Updates",
		summary:  "To provide a smooth protocol layer launch and user experience, various network parameters must be updated to the most recent and refined calculations. This parameter change proposal, if passed, would adopt newly calculated data pool parameters on the KYVE network, increasing the overall security and functionality of the protocol layer.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	7: {
		title:    "Efficient Frontier – KYVE Market Making Proposal",
		summary:  "Efficient Frontier, a cryptocurrency algorithmic trading and market-making firm, proposes a market-making agreement to the KYVE community. The proposal involves KYVE Foundation loaning $1,225,000 worth of KYVE and $125,000 worth of USDC to Efficient Frontier for market-making purposes. The total KYVE transferred will be determined by the prices at the time of token delivery, and Efficient Frontier commits to making markets on CeFi exchanges and the Osmosis DEX. The loaned amount will either be returned at the end of the market-making period or purchased at pre-determined strike prices. The intention is to provide ample liquidity for the KYVE ecosystem, facilitating efficient trading on both centralized and decentralized exchanges.",
		proposer: "kyve13jcsyasn6g56e2advm792ukp6cwaxn4pe9rtdl",
	},
	8: {
		title:    "v1.3.0 Software Upgrade",
		summary:  "If approved, this proposed software upgrade would introduce a weighted round-robin method for selecting uploaders and also incorporate a split of inflation between the chain and protocol. Additionally, this upgrade would correctly track the delegations from investor accounts affected in the `v1.1.0` upgrade.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	9: {
		title:    "Pool Operating Cost (Follow-up Proposal)",
		summary:  "This proposal is a follow-up proposal to the v1.3.0 upgrade which introduces inflation splitting. After the upgrade is passed (which will be approx. at 1st Aug. 12 pm UTC) the operating cost for the pool needs to be adjusted to maintain the same pool economics.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	10: {
		title:    "KYVE Validator Delegation Program Update",
		summary:  "Since the launch of the protocol layer, the [KYVE Validator Delegation Program](https://blog.kyve.network/kyve-validator-delegation-program-update-mainnet-delegation-and-curated-airdrop-2c615b8ee76f) (rewarding the top Incentivized Testnet participants, KVDP) is the layer’s largest source of delegation. We have become aware that certain validators within this program are taking advantage of the situation and setting their commission rate to 100% to get the most economic benefit. In doing so, those who delegate to them do not receive a part of the rewards, discouraging delegation. In order to give new protocol validators a fair chance, increase network competition, and further support delegation, we propose limiting the commission rate to 50% for those within the KYVE Validator Delegation Program. These validators would have 96 hours (4 days) to comply with the new limit, along with a second chance option for those who missed the implementation of this new rule.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	11: {
		title:    "Data Pool Creation: Osmosis",
		summary:  "As a next step in expanding our decentralized data lake, we propose the creation of a new data pool on Mainnet. This data pool creation proposal, if passed, would enact the validation and archiving of Osmosis blocks and the corresponding block results.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	12: {
		title:    "Data Pool Creation: Archway",
		summary:  "As a subsequent action in the progression of our decentralized data lake’s growth, we suggest establishing a fresh data pool on the Mainnet. The proposal for setting up this data pool, upon approval, will lead to the verification and preservation of Archway blocks along with their associated block results.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	13: {
		title:    "Data Pool Creation: Axelar",
		summary:  "As a subsequent action in the progression of our decentralized data lake’s growth, we suggest establishing a new data pool on mainnet. The proposal for setting up this data pool, upon approval, will lead to the validation and archiving of Axelar blocks along with their associated block results.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	14: {
		title:    "Parameter Change: Min Delegation",
		summary:  "Since the launch of the first pool Cosmos Hub back in June 2023, the $KYVE delegated into the protocol significantly increased. For security reasons the min delegation has to be increased in order to guarantee that no node can have more than 50% voting power inside a pool.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	15: {
		title:    "Parameter Change: Storage Cost",
		summary:  "The storage_cost parameter is responsible for estimating a fair payout for the uploader to compensate them for uploading data to storage providers like Arweave or Bundlr. Since this parameter has not been updated since genesis, and real market conditions can now be observed, we propose an updated storage cost value of 0.103996 ukyve per byte.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	16: {
		title:    "Software Upgrade: Pool binaries",
		summary:  "We propose the software upgrade of all `kyvejs/tendermint` and `kyvejs/tendermint-bsync` binaries, affecting all pools. This software upgrade proposal, if passed, would enact the upgrade of the binaries for each pool to the specified versions.",
		proposer: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	},
	17: {
		title:    "Pool Creation: Archway // State-Sync",
		summary:  "As a next step in expanding our decentralized data lake, we propose the creation of a new data pool, initially launching on Kaon. This data pool creation proposal, if passed, would enact the validation and archiving of Archway state-sync snapshots. With this, it will be possible to use archived and validated state-sync snapshots from KYVE via KSYNC to rapidly join the Archway network.",
		proposer: "kyve1fst07guqk7u4elhfj4z79fzgkr5wtfqqm7w20n",
	},
	18: {
		title:    "Upgrade Pool Runtimes",
		summary:  "To participate as a Protocol Validator in the KYVE Network, it is essential to use the kyvejs/tendermint and kyvejs/tendermint-bsync runtimes as they implement the validation logic. We have identified an issue with the core implementation and an improper handling of non-deterministic data on Osmosis block results. Therefore, we propose to update the minimum required binary version of all pools to ensure they incorporate the latest release with fixes for those problems to avoid future incorrect vote slashes.",
		proposer: "kyve1fst07guqk7u4elhfj4z79fzgkr5wtfqqm7w20n",
	},
	19: {
		title:    "Parameter Change: Min Delegation",
		summary:  "Since the launch of the first pool Cosmos Hub back in June 2023, the $KYVE delegated into the protocol significantly increased. For security reasons the min delegation has to be increased in order to guarantee that no node can have more than 50% voting power inside a pool.",
		proposer: "kyve1fst07guqk7u4elhfj4z79fzgkr5wtfqqm7w20n",
	},
	20: {
		title:    "Data Pool Creation: Cronos & Cronos // State-Sync",
		summary:  "As a subsequent action in the progression of our decentralized data lake’s growth, we suggest establishing a new pair of data pools on mainnet. The proposal for setting up those data pools, upon approval, will lead to the validation and archiving of Cronos blocks, block results, and Cronos state-sync snapshots.",
		proposer: "kyve1fst07guqk7u4elhfj4z79fzgkr5wtfqqm7w20n",
	},
	21: {
		title:    "v1.4.0 Software Upgrade",
		summary:  "This software upgrade proposal, if passed, would upgrade the KYVE Mainnet network to version 1.4 featuring CosmosSDK v0.47 and the implementation of the new funders concept, as discussed [here](https://commonwealth.im/kyve/discussion/13420-enhancing-kyves-funders-concept). After ample testing on the Kaon testnet, this upgrade is now proposed for mainnet.",
		proposer: "kyve1fst07guqk7u4elhfj4z79fzgkr5wtfqqm7w20n",
	},
	22: {
		title:    "Parameter Change: Storage Cost",
		summary:  "The storage_cost parameter is responsible for estimating a fair payout for the uploader to compensate them for uploading data to storage providers like Arweave or Bundlr. Since this parameter has not been updated since genesis, and real market conditions can now be observed, we propose an updated storage cost value of 0.2772 ukyve per byte.",
		proposer: "kyve1fst07guqk7u4elhfj4z79fzgkr5wtfqqm7w20n",
	},
}
