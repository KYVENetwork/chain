package keeper_test

import (
	"fmt"
	"testing"

	"github.com/KYVENetwork/chain/x/bundles/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBundlesKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Keeper Test Suite", types.ModuleName))
}
