package app

import (
	"cosmossdk.io/log"
	bundlestypes "github.com/KYVENetwork/chain/x/bundles/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
	"slices"
)

type PriorityProposalHandler struct {
	logger    log.Logger
	txDecoder sdk.TxDecoder
}

func NewPriorityProposalHandler(logger log.Logger, decoder sdk.TxDecoder) *PriorityProposalHandler {
	return &PriorityProposalHandler{
		logger:    logger,
		txDecoder: decoder,
	}
}

var priorityTypes = []string{
	reflect.TypeOf(bundlestypes.MsgSubmitBundleProposal{}).Name(),
	reflect.TypeOf(bundlestypes.MsgVoteBundleProposal{}).Name(),
	reflect.TypeOf(bundlestypes.MsgClaimUploaderRole{}).Name(),
	reflect.TypeOf(bundlestypes.MsgSkipUploaderRole{}).Name(),
}

// PrepareProposal returns a PrepareProposalHandler that separates transactions into different queues
// This function is only called by the block proposer and therefore does NOT need to be deterministic
func (h *PriorityProposalHandler) PrepareProposal() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		// Separate the transactions into different queues
		// priorityQueue: transactions that should be executed before the default transactions
		// defaultQueue: transactions that should be executed last
		// The order is priorityQueue -> defaultQueue

		var priorityQueue [][]byte
		var defaultQueue [][]byte

		// Iterate through the transactions and separate them into different queues
		for _, rawTx := range req.Txs {
			tx, err := h.txDecoder(rawTx)
			if err != nil {
				h.logger.Error("failed to decode transaction", "error", err)
				continue
			}
			msgs, err := tx.GetMsgsV2()
			if err != nil {
				h.logger.Error("failed to get messages from transaction", "error", err)
				continue
			}

			// We only care about transactions with a single message
			if len(msgs) == 1 {
				msg := msgs[0]
				msgType := string(msg.ProtoReflect().Type().Descriptor().Name())

				if slices.Contains(priorityTypes, msgType) {
					priorityQueue = append(priorityQueue, rawTx)
					continue
				}
			}

			// Otherwise, add the tx to the default queue
			defaultQueue = append(defaultQueue, rawTx)
		}

		// Append the transactions in the correct order
		// If we have reached the max tx bytes, we stop adding more transactions
		var orderedTxs [][]byte
		var countTxBytes int64
		for _, rawTx := range append(priorityQueue, defaultQueue...) {
			countTxBytes += int64(len(rawTx))
			if countTxBytes > req.MaxTxBytes {
				h.logger.Info("max tx bytes reached", "max_tx_bytes", req.MaxTxBytes)
				break
			}
			orderedTxs = append(orderedTxs, rawTx)
		}

		return &abci.ResponsePrepareProposal{
			Txs: orderedTxs,
		}, nil
	}
}
