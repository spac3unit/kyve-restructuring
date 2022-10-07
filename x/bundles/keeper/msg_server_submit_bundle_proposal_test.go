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

TEST CASES - msg_server_submit_bundle_proposal.go

// TODO: implement
* Call can propose if pool does not exist
* Call can propose if pool is currently upgrading
* Call can propose if pool is paused
* Call can propose if pool is out of funds
* Call can propose if pool has not reached the minimum stake
* Call can propose with a valaccount which does not exist
* Call can propose as a staker who is not the next uploader
* Call can propose before the upload interval passed
* Call can propose with an invalid from height

* Submit the first bundle proposal with empty storage id
* Submit the first bundle proposal with empty byte size
* Submit the first bundle proposal with bigger bundle size than allowed
* Submit the first bundle proposal with empty bundle
* Submit the first bundle proposal with empty value
* Submit the first bundle proposal with empty bundle hash
* Submit a bundle proposal with valid args

*/

var _ = Describe("msg_server_submit_bundle_proposal.go", Ordered, func() {
	s := i.NewCleanChain()

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
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.ALICE,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.ALICE,
			PoolId:     0,
			Valaddress: i.BOB,
		})

		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.BOB,
			Staker:  i.ALICE,
			PoolId:  0,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Submit the first bundle proposal with empty storage id", func() {
		// ARRANGE
		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.BOB,
			Staker:     i.ALICE,
			PoolId:     0,
			StorageId:  "",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   100,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "test_hash",
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
	})

	It("Submit the first bundle proposal with empty byte size", func() {
		// ARRANGE
		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.BOB,
			Staker:     i.ALICE,
			PoolId:     0,
			StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			ByteSize:   0,
			FromHeight: 0,
			ToHeight:   100,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "test_hash",
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
	})

	It("Submit the first bundle proposal with bigger bundle size than allowed", func() {
		// ARRANGE
		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.BOB,
			Staker:     i.ALICE,
			PoolId:     0,
			StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   101,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "test_hash",
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
	})

	It("Submit the first bundle proposal with empty bundle", func() {
		// ARRANGE
		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.BOB,
			Staker:     i.ALICE,
			PoolId:     0,
			StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   0,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "test_hash",
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
	})

	It("Submit the first bundle proposal with empty value", func() {
		// ARRANGE
		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.BOB,
			Staker:     i.ALICE,
			PoolId:     0,
			StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   100,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "",
			BundleHash: "test_hash",
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
	})

	It("Submit the first bundle proposal with empty bundle hash", func() {
		// ARRANGE
		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.BOB,
			Staker:     i.ALICE,
			PoolId:     0,
			StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   100,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "",
		})

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
	})

	It("Submit the first bundle proposal with valid args", func() {
		// ARRANGE
		s.CommitAfterSeconds(60)
		// TODO: test if waiting time is greater than 60 seconds

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.BOB,
			Staker:     i.ALICE,
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

		// ASSERT
		bundleProposal, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeTrue())

		Expect(bundleProposal.PoolId).To(Equal(uint64(0)))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))
		Expect(bundleProposal.Uploader).To(Equal(i.ALICE))
		Expect(bundleProposal.NextUploader).To(Equal(i.ALICE))
		Expect(bundleProposal.ByteSize).To(Equal(uint64(100)))
		Expect(bundleProposal.ToHeight).To(Equal(uint64(100)))
		Expect(bundleProposal.ToKey).To(Equal("99"))
		Expect(bundleProposal.ToValue).To(Equal("test_value"))
		Expect(bundleProposal.BundleHash).To(Equal("test_hash"))
		Expect(bundleProposal.CreatedAt).NotTo(BeZero())
		Expect(bundleProposal.VotersValid).To(ContainElement(i.ALICE))
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())
	})
})
