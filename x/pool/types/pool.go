package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (m *Pool) GetPoolAccount() sdk.AccAddress {
	name := fmt.Sprintf("%s/%d", ModuleName, m.Id)

	return authTypes.NewModuleAddress(name)
}
