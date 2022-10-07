package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Keeper Suite")
}

func PayoutRewards(s *i.KeeperTestSuite, staker string, amount uint64) {
	err := s.App().PoolKeeper.ChargeFundersOfPool(s.Ctx(), 0, amount)
	Expect(err).To(BeNil())
	success := s.App().DelegationKeeper.PayoutRewards(s.Ctx(), staker, amount, pooltypes.ModuleName)
	Expect(success).To(BeTrue())
}

func CreateFundedPool(s *i.KeeperTestSuite) {
	s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
		Name: "Moontest",
		Protocol: &pooltypes.Protocol{
			Version:     "0.0.0",
			Binaries:    "{}",
			LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
		},
		UpgradePlan: &pooltypes.UpgradePlan{},
	})

	s.CommitAfterSeconds(7)

	s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
		Creator: i.ALICE,
		Id:      0,
		Amount:  50 * i.KYVE,
	})

	s.CommitAfterSeconds(7)

	pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
	Expect(poolFound).To(BeTrue())
	Expect(pool.TotalFunds).To(Equal(50 * i.KYVE))
}

func CheckAndContinueChainForOneMonth(s *i.KeeperTestSuite) {
	s.PerformValidityChecks()
	for d := 0; d < 31; d++ {
		s.CommitAfterSeconds(60 * 60 * 24)
		s.PerformValidityChecks()
	}
}
