package keeper

import (
	"context"

	"github.com/KYVENetwork/chain/x/stakers/types"
)

// UpdateCommission creates a queue entry to update the staker commission.
// After the `CommissionChangeTime` is over the new commission will be applied.
// If an update is currently in the queue it will get removed from the queue
// and the user needs to wait again for the full time to pass.
func (k msgServer) UpdateCommission(goCtx context.Context, msg *types.MsgUpdateCommission) (*types.MsgUpdateCommissionResponse, error) {
	// TODO no-op, will be replaced by a per pool commission

	return &types.MsgUpdateCommissionResponse{}, nil
}
