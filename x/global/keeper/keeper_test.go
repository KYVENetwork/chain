package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KYVENetwork/chain/x/global/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGlobalKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Keeper Test Suite", types.ModuleName))
}
