package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_gov_test.go

* Create Pool
* Update Pool
* Pause Pool
* Pause Pool when already paused
* Unpause Pool
* Unpause Pool when already unpaused
* Create Upgrade Pool proposal

*/

var _ = Describe("Gov Pool", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		s = i.NewCleanChain()

		err := s.App().PoolKeeper.CreatePool(s.Ctx(), &pooltypes.CreatePoolProposal{
			Title:          i.GOV,
			Description:    "desc",
			Name:           "Moonbeam",
			Runtime:        "@kyve/evm",
			Logo:           "https://arweave.net/9FJDam56yBbmvn8rlamEucATH5UcYqSBw468rlCXn8E",
			Config:         "{\"config\": \"test\"}",
			StartKey:       "0",
			UploadInterval: 60,
			OperatingCost:  2_500_000_000,
			MinStake:       100_000_000_000,
			MaxBundleSize:  100,
			Version:        "1",
			Binaries:       "{\"b1\": \"string\"}",
		})
		Expect(err).To(BeNil())

	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Create Pool", func() {

		// Overwrite BeforeEach
		s = i.NewCleanChain()

		// Arrange
		_, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(found).To(BeFalse())

		// Act
		err := s.App().PoolKeeper.CreatePool(s.Ctx(), &pooltypes.CreatePoolProposal{
			Title:          i.GOV,
			Description:    "desc",
			Name:           "Moonbeam",
			Runtime:        "@kyve/evm",
			Logo:           "https://arweave.net/9FJDam56yBbmvn8rlamEucATH5UcYqSBw468rlCXn8E",
			Config:         "{\"config\": \"test\"}",
			StartKey:       "0",
			UploadInterval: 60,
			OperatingCost:  2_500_000_000,
			MinStake:       100_000_000_000,
			MaxBundleSize:  100,
			Version:        "1",
			Binaries:       "{\"b1\": \"string\"}",
		})
		Expect(err).To(BeNil())

		// Assert
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(found).To(BeTrue())
		Expect(pool.Name).To(Equal("Moonbeam"))
		Expect(pool.Runtime).To(Equal("@kyve/evm"))
		Expect(pool.Logo).To(Equal("https://arweave.net/9FJDam56yBbmvn8rlamEucATH5UcYqSBw468rlCXn8E"))
		Expect(pool.Config).To(Equal("{\"config\": \"test\"}"))
		Expect(pool.StartKey).To(Equal("0"))
		Expect(pool.UploadInterval).To(Equal(uint64(60)))
		Expect(pool.OperatingCost).To(Equal(uint64(2_500_000_000)))
		Expect(pool.MinStake).To(Equal(uint64(100_000_000_000)))
		Expect(pool.MaxBundleSize).To(Equal(uint64(100)))
		Expect(pool.Protocol.Version).To(Equal("1"))
		Expect(pool.Protocol.Binaries).To(Equal("{\"b1\": \"string\"}"))

	})

	It("Update Pool", func() {
		// Act
		err := s.App().PoolKeeper.UpdatePool(s.Ctx(), &pooltypes.UpdatePoolProposal{
			Title:       "gov",
			Description: "desc",
			Id:          0,
			Payload:     "{\"name\": \"Bitcoin\"}",
		})

		// Assert
		Expect(err).To(BeNil())

		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(found).To(BeTrue())
		Expect(pool.Name).To(Equal("Bitcoin"))
		Expect(pool.Runtime).To(Equal("@kyve/evm"))
		Expect(pool.Logo).To(Equal("https://arweave.net/9FJDam56yBbmvn8rlamEucATH5UcYqSBw468rlCXn8E"))
		Expect(pool.Config).To(Equal("{\"config\": \"test\"}"))
		Expect(pool.StartKey).To(Equal("0"))
		Expect(pool.UploadInterval).To(Equal(uint64(60)))
		Expect(pool.OperatingCost).To(Equal(uint64(2_500_000_000)))
		Expect(pool.MinStake).To(Equal(uint64(100_000_000_000)))
		Expect(pool.MaxBundleSize).To(Equal(uint64(100)))
		Expect(pool.Protocol.Version).To(Equal("1"))
		Expect(pool.Protocol.Binaries).To(Equal("{\"b1\": \"string\"}"))
	})

	It("Pause Pool", func() {
		// Arrange
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(found).To(BeTrue())
		Expect(pool.Paused).To(BeFalse())

		// Act
		err := s.App().PoolKeeper.PausePool(s.Ctx(), &pooltypes.PausePoolProposal{
			Title:       i.GOV,
			Description: "desc",
			Id:          0,
		})
		Expect(err).To(BeNil())

		// Assert
		poolAfter, foundAfter := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(foundAfter).To(BeTrue())
		Expect(poolAfter.Paused).To(BeTrue())

		poolAfter.Paused = false
		Expect(pool).To(Equal(poolAfter))
	})

	It("Pause Pool when already paused", func() {
		// Arrange
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.Paused = true
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)
		Expect(found).To(BeTrue())
		Expect(pool.Paused).To(BeTrue())

		// Act
		err := s.App().PoolKeeper.PausePool(s.Ctx(), &pooltypes.PausePoolProposal{
			Title:       i.GOV,
			Description: "desc",
			Id:          0,
		})
		Expect(err).ToNot(BeNil())

		// Assert
		poolAfter, foundAfter := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(foundAfter).To(BeTrue())
		Expect(pool).To(Equal(poolAfter))
	})

	It("Unpause Pool", func() {
		// Arrange
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.Paused = true
		s.App().PoolKeeper.SetPool(s.Ctx(), pool)
		Expect(found).To(BeTrue())
		Expect(pool.Paused).To(BeTrue())

		// Act
		err := s.App().PoolKeeper.UnpausePool(s.Ctx(), &pooltypes.UnpausePoolProposal{
			Title:       i.GOV,
			Description: "desc",
			Id:          0,
		})
		Expect(err).To(BeNil())

		// Assert
		poolAfter, foundAfter := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(foundAfter).To(BeTrue())
		Expect(poolAfter.Paused).To(BeFalse())

		poolAfter.Paused = true
		Expect(pool).To(Equal(poolAfter))
	})

	It("Unpause Pool when already unpaused", func() {
		// Arrange
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(found).To(BeTrue())
		Expect(pool.Paused).To(BeFalse())

		// Act
		err := s.App().PoolKeeper.UnpausePool(s.Ctx(), &pooltypes.UnpausePoolProposal{
			Title:       i.GOV,
			Description: "desc",
			Id:          0,
		})
		Expect(err).ToNot(BeNil())

		// Assert
		poolAfter, foundAfter := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(foundAfter).To(BeTrue())
		Expect(pool).To(Equal(poolAfter))
	})

	It("Create Upgrade Pool proposal", func() {
		// Arrange
		pool, found := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(found).To(BeTrue())
		Expect(pool.Paused).To(BeFalse())

		// Act
		s.App().PoolKeeper.UpgradePool(s.Ctx(), &pooltypes.SchedulePoolUpgradeProposal{
			Title:       i.GOV,
			Description: "desc",
			Runtime:     "@kyve/evm",
			Version:     "new version",
			ScheduledAt: uint64(s.Ctx().BlockTime().Unix() + 1000),
			Duration:    60,
			Binaries:    "{\"test\": \"link.com\"}",
		})

		// Assert
		poolAfter, foundAfter := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(foundAfter).To(BeTrue())

		Expect(poolAfter.UpgradePlan.Version).To(Equal("new version"))
		Expect(poolAfter.UpgradePlan.ScheduledAt).To(Equal(uint64(s.Ctx().BlockTime().Unix() + 1000)))
		Expect(poolAfter.UpgradePlan.Duration).To(Equal(uint64(60)))
		Expect(poolAfter.UpgradePlan.Binaries).To(Equal("{\"test\": \"link.com\"}"))

		// Fast-forward
		s.CommitAfterSeconds(2000)
		s.CommitAfterSeconds(1)

		//TODO check
		//Expect(poolAfter.UpgradePlan).To(Equal(pooltypes.UpgradePlan{}))
	})
})
