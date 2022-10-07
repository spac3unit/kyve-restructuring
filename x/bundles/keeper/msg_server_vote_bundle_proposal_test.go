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

TEST CASES - msg_server_vote_bundle_proposal.go

* Try to vote valid as the only voter on bundle proposal
* Try to vote invalid as the only voter on bundle proposal
* Try to vote abstain as the only voter on bundle proposal
* Try to vote abstain on proposal again
* Try to vote valid on proposal after abstain vote
* Try to vote invalid on proposal after abstain vote
* Try to vote unspecified on proposal
* TODO: Try to vote as not the first voter on bundle proposal

*/

var _ = Describe("msg_server_vote_bundle_proposal.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create clean pool for every test case
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:          "Moontest",
			MaxBundleSize: 100,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
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

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Try to vote valid as the only voter on bundle proposal", func() {
		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).NotTo(ContainElement(i.STAKER_1))
	})

	It("Try to vote invalid as the only voter on bundle proposal", func() {
		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_NO,
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).NotTo(ContainElement(i.STAKER_1))
	})

	It("Try to vote abstain as the only voter on bundle proposal", func() {
		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).To(ContainElement(i.STAKER_1))
	})

	It("Try to vote abstain on proposal again", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).To(ContainElement(i.STAKER_1))
	})

	It("Try to vote valid on proposal after abstain vote", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).NotTo(ContainElement(i.STAKER_1))
	})

	It("Try to vote invalid on proposal after abstain vote", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_NO,
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).To(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).NotTo(ContainElement(i.STAKER_1))
	})

	It("Try to vote unspecified on proposal", func() {
		// ACT
		s.RunTxBundlesError(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			Vote:      bundletypes.VOTE_TYPE_UNSPECIFIED,
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.VotersValid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersInvalid).NotTo(ContainElement(i.STAKER_1))
		Expect(bundleProposal.VotersAbstain).NotTo(ContainElement(i.STAKER_1))
	})
})
