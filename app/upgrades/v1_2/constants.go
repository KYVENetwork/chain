package v1_2

// UpgradeName is the name of this specific software upgrade used on-chain.
const UpgradeName = "v1.2.0"

// TestnetChainID is the Chain ID of the KYVE testnet (Kaon).
const TestnetChainID = "kaon-1"

// MainnetChainID is the Chain ID of the KYVE mainnet.
const MainnetChainID = "kyve-1"

// TestnetProposers is a mapping between Proposal ID and Proposer Address for
// the KYVE testnet (Kaon).
var TestnetProposers = map[uint64]string{
	1: "kyve1qfnu3k5vwgfyzsaqe0w4ssqp99delgtg2qz0jc",
	2: "kyve1qfnu3k5vwgfyzsaqe0w4ssqp99delgtg2qz0jc",
	3: "kyve1qfnu3k5vwgfyzsaqe0w4ssqp99delgtg2qz0jc",
	4: "kyve1qfnu3k5vwgfyzsaqe0w4ssqp99delgtg2qz0jc",
}

// MainnetProposers is a mapping between Proposal ID and Proposer Address for
// the KYVE mainnet.
var MainnetProposers = map[uint64]string{
	1: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
	2: "kyve1afu42029ujjcja4yry3rx6x43k33k88ep5wvjz",
}
