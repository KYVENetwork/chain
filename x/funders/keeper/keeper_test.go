package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KYVENetwork/chain/x/funders/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFundersKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Keeper Test Suite", types.ModuleName))
}
