package keeper_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - dropped bundles

* Produce a dropped bundle because not enough validators voted
* TODO: Produce a dropped bundle because the only funder can not pay for the bundle reward
* TODO: Produce a dropped bundle because multiple funders with same amount can not pay for the bundle reward
* TODO: Produce a dropped bundle because multiple funders with different amount can not pay for the bundle reward

*/

var _ = Describe("dropped bundles", Ordered, func() {
	s := i.NewCleanChain()

	initialBalanceStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
	initialBalanceValaddress0 := s.GetBalanceFromAddress(i.VALADDRESS_0)

	initialBalanceStaker1 := s.GetBalanceFromAddress(i.STAKER_1)
	initialBalanceValaddress1 := s.GetBalanceFromAddress(i.VALADDRESS_1)

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create clean pool for every test case
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "Moontest",
			MaxBundleSize:  100,
			StartKey:       "0",
			UploadInterval: 60,
			OperatingCost:  10_000,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceValaddress0 = s.GetBalanceFromAddress(i.VALADDRESS_0)

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)
		initialBalanceValaddress1 = s.GetBalanceFromAddress(i.VALADDRESS_1)

		s.CommitAfterSeconds(60)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Produce a dropped bundle because not enough validators voted", func() {
		// ARRANGE
		// stake a bit more than first node so >50% is reached
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  200 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
		})

		initialBalanceStaker1 = s.GetBalanceFromAddress(i.STAKER_1)

		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.VALADDRESS_0,
			Staker:     i.STAKER_0,
			PoolId:     0,
			StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   100,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "test_hash",
		})

		// ACT
		// do not vote so bundle gets dropped
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		// check if bundle got not finalized on pool
		pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		Expect(poolFound).To(BeTrue())

		Expect(pool.CurrentKey).To(Equal(""))
		Expect(pool.CurrentValue).To(BeEmpty())
		Expect(pool.CurrentHeight).To(BeZero())
		Expect(pool.TotalBundles).To(BeZero())

		// check if finalized bundle exists
		_, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
		Expect(finalizedBundleFound).To(BeFalse())

		// check if bundle proposal got dropped
		bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposalFound).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.ByteSize).To(BeZero())
		Expect(bundleProposal.ToHeight).To(BeZero())
		Expect(bundleProposal.ToKey).To(BeEmpty())
		Expect(bundleProposal.ToValue).To(BeEmpty())
		Expect(bundleProposal.BundleHash).To(BeEmpty())
		Expect(bundleProposal.CreatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(BeEmpty())
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		// check uploader status
		valaccountUploader, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)
		Expect(valaccountUploader.Points).To(BeZero())

		balanceValaddress := s.GetBalanceFromAddress(valaccountUploader.Valaddress)
		Expect(balanceValaddress).To(Equal(initialBalanceValaddress0))

		balanceUploader := s.GetBalanceFromAddress(valaccountUploader.Staker)

		Expect(balanceUploader).To(Equal(initialBalanceStaker0))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(BeZero())

		// check voter status
		valaccountVoter, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_1)
		Expect(valaccountVoter.Points).To(Equal(uint64(1)))

		balanceVoterValaddress := s.GetBalanceFromAddress(valaccountVoter.Valaddress)
		Expect(balanceVoterValaddress).To(Equal(initialBalanceValaddress1))

		balanceVoter := s.GetBalanceFromAddress(valaccountVoter.Staker)
		Expect(balanceVoter).To(Equal(initialBalanceStaker1))

		Expect(balanceVoter).To(Equal(initialBalanceStaker1))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.STAKER_1, i.STAKER_1)).To(BeZero())

		// check pool funds
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		funder, _ := pool.GetFunder(i.ALICE)

		Expect(pool.Funders).To(HaveLen(1))
		Expect(funder.Amount).To(Equal(100 * i.KYVE))
	})

	//PIt("Produce dropped bundle because pool has not enough funds", func() {
	//	// ARRANGE
	//	s.RunTxPoolSuccess(&pooltypes.MsgDefundPool{
	//		Creator: i.ALICE,
	//		Amount:  100 * i.KYVE,
	//	})
	//
	//	// fund amount which definetely not cover bundle reward
	//	s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
	//		Creator: i.ALICE,
	//		Amount:  1,
	//	})
	//
	//	initialBalanceAlice := s.GetBalanceFromAddress(i.ALICE)
	//
	//	// stake a bit more than first node so >50% is reached
	//	s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
	//		Creator: i.BOB,
	//		Amount:  200 * i.KYVE,
	//	})
	//
	//	s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
	//		Creator:    i.BOB,
	//		PoolId:     0,
	//		Valaddress: i.ALICE,
	//	})
	//
	//	initialBalanceValaddress = s.GetBalanceFromAddress(i.BOB)
	//
	//	s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
	//		Creator:    i.BOB,
	//		Staker:     i.ALICE,
	//		PoolId:     0,
	//		StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
	//		ByteSize:   100,
	//		FromHeight: 0,
	//		ToHeight:   100,
	//		FromKey:    "0",
	//		ToKey:      "99",
	//		ToValue:    "test_value",
	//		BundleHash: "test_hash",
	//	})
	//
	//	// ACT
	//	s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
	//		Creator:   i.ALICE,
	//		Staker:    i.BOB,
	//		PoolId:    0,
	//		StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
	//		Vote:      bundletypes.VOTE_TYPE_YES,
	//	})
	//
	//	s.CommitAfterSeconds(60)
	//
	//	s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
	//		Creator:    i.BOB,
	//		Staker:     i.ALICE,
	//		PoolId:     0,
	//		StorageId:  "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
	//		ByteSize:   100,
	//		FromHeight: 100,
	//		ToHeight:   200,
	//		FromKey:    "99",
	//		ToKey:      "199",
	//		ToValue:    "test_value2",
	//		BundleHash: "test_hash2",
	//	})
	//
	//	// ASSERT
	//	// check if bundle got not finalized on pool
	//	pool, poolFound := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
	//	Expect(poolFound).To(BeTrue())
	//
	//	Expect(pool.CurrentKey).To(Equal(""))
	//	Expect(pool.CurrentValue).To(BeEmpty())
	//	Expect(pool.CurrentHeight).To(BeZero())
	//	Expect(pool.TotalBundles).To(BeZero())
	//
	//	// check if finalized bundle exists
	//	_, finalizedBundleFound := s.App().BundlesKeeper.GetFinalizedBundle(s.Ctx(), 0, 0)
	//	Expect(finalizedBundleFound).To(BeFalse())
	//
	//	// check if bundle proposal got dropped
	//	bundleProposal, bundleProposalFound := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
	//	Expect(bundleProposalFound).To(BeTrue())
	//
	//	Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
	//	Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
	//	Expect(bundleProposal.Uploader).To(Equal(i.ALICE))
	//	Expect(bundleProposal.NextUploader).To(Equal(i.ALICE))
	//	Expect(bundleProposal.ByteSize).To(Equal(uint64(100)))
	//	Expect(bundleProposal.ToHeight).To(Equal(uint64(100)))
	//	Expect(bundleProposal.ToKey).To(Equal("99"))
	//	Expect(bundleProposal.ToValue).To(Equal("test_value"))
	//	Expect(bundleProposal.BundleHash).To(Equal("test_hash"))
	//	Expect(bundleProposal.CreatedAt).NotTo(BeZero())
	//	Expect(bundleProposal.VotersValid).To(ContainElement(i.ALICE))
	//	Expect(bundleProposal.VotersInvalid).To(BeEmpty())
	//	Expect(bundleProposal.VotersAbstain).To(BeEmpty())
	//
	//	// check uploader status
	//	valaccountUploader, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.ALICE)
	//	Expect(valaccountUploader.Points).To(BeZero())
	//
	//	balanceValaddress := s.GetBalanceFromAddress(valaccountUploader.Valaddress)
	//	Expect(balanceValaddress).To(Equal(initialBalanceValaddress))
	//
	//	balanceStaker := s.GetBalanceFromAddress(valaccountUploader.Staker)
	//
	//	Expect(balanceStaker).To(Equal(initialBalanceAlice))
	//
	//	//staker, stakerFound := s.App().StakersKeeper.GetStaker(s.Ctx(), valaccountUploader.Staker)
	//	//Expect(stakerFound).To(BeTrue())
	//
	//	//Expect(staker.Amount).To(Equal(100 * i.KYVE)) TODO
	//	//Expect(s.App().StakersKeeper.GetTotalStake(s.Ctx(), 0)).To(Equal(300 * i.KYVE))
	//
	//	// check voter status
	//	valaccountVoter, _ := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.BOB)
	//	Expect(valaccountVoter.Points).To(BeZero())
	//
	//	// check pool funds
	//	pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
	//
	//	Expect(pool.Funders).To(BeEmpty())
	//})
})
