package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_end_block_handle_pool_upgrades.go

* Schedule a pool upgrade in the past
* Schedule a pool upgrade in the future
* Schedule pool upgrade now and with no upgrade duration
* // TODO: schedule pool upgrade with same version but different binaries
* // TODO: Try to schedule pool upgrade with same binaries but different version
* // TODO: Try to schedule pool upgrade with same version and same binaries

*/

var _ = Describe("logic_end_block_handle_pool_upgrades.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create clean pool for every test case
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "Moontest",
			MaxBundleSize:  100,
			StartKey:       "0",
			MinStake:       100 * i.KYVE,
			UploadInterval: 60,
			OperatingCost:  10_000,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Schedule a pool upgrade in the past", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		upgradePlan := pooltypes.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{\"linux\":\"test\"}",
			ScheduledAt: uint64(s.Ctx().BlockTime().Unix()) - 120,
			Duration:    3600,
		}

		pool.UpgradePlan = &upgradePlan

		// ACT
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)
		s.CommitAfterSeconds(1)

		// ASSERT
		// check if pool is currently upgrading
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Protocol.Version).To(Equal(upgradePlan.Version))
		Expect(pool.Protocol.Binaries).To(Equal(upgradePlan.Binaries))

		Expect(pool.UpgradePlan).To(Equal(&upgradePlan))

		// check if upgrade is done after upgrade duration
		s.CommitAfterSeconds(3600)
		s.CommitAfterSeconds(1)

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Protocol.Version).To(Equal(upgradePlan.Version))
		Expect(pool.Protocol.Binaries).To(Equal(upgradePlan.Binaries))

		Expect(pool.UpgradePlan).To(Equal(&pooltypes.UpgradePlan{}))
	})

	It("Schedule a pool upgrade in the future", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		upgradePlan := pooltypes.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{\"linux\":\"test\"}",
			ScheduledAt: uint64(s.Ctx().BlockTime().Unix()) + 120,
			Duration:    3600,
		}

		pool.UpgradePlan = &upgradePlan

		// ACT
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)
		s.CommitAfterSeconds(1)

		// ASSERT
		// check if pool upgrade is still only scheduled
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Protocol.Version).To(Equal("0.0.0"))
		Expect(pool.Protocol.Binaries).To(Equal("{}"))

		Expect(pool.UpgradePlan).To(Equal(&upgradePlan))

		s.CommitAfterSeconds(120)
		s.CommitAfterSeconds(1)

		// check if pool is currently upgrading
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Protocol.Version).To(Equal(upgradePlan.Version))
		Expect(pool.Protocol.Binaries).To(Equal(upgradePlan.Binaries))

		Expect(pool.UpgradePlan).To(Equal(&upgradePlan))

		s.CommitAfterSeconds(3600)
		s.CommitAfterSeconds(1)

		// check if upgrade is done after upgrade duration
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Protocol.Version).To(Equal(upgradePlan.Version))
		Expect(pool.Protocol.Binaries).To(Equal(upgradePlan.Binaries))

		Expect(pool.UpgradePlan).To(Equal(&pooltypes.UpgradePlan{}))
	})

	It("Schedule pool upgrade now and with no upgrade duration", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		upgradePlan := pooltypes.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{\"linux\":\"test\"}",
			ScheduledAt: uint64(s.Ctx().BlockTime().Unix()),
			Duration:    0,
		}

		pool.UpgradePlan = &upgradePlan

		// ACT
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)
		s.CommitAfterSeconds(1)

		// ASSERT
		// check if upgrade is done after upgrade duration
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(pool.Protocol.Version).To(Equal(upgradePlan.Version))
		Expect(pool.Protocol.Binaries).To(Equal(upgradePlan.Binaries))

		Expect(pool.UpgradePlan).To(Equal(&pooltypes.UpgradePlan{}))
	})
})
