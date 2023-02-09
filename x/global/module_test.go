package global_test

import (
	"fmt"
	"testing"

	"github.com/KYVENetwork/chain/x/global/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGlobalModule(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Test Suite", types.ModuleName))
}
