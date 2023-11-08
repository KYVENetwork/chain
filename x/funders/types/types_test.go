package types_test

import (
	"fmt"
	"testing"

	"github.com/KYVENetwork/chain/x/team/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTeamKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Types Test Suite", types.ModuleName))
}
