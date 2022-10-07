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

TEST CASES - msg_server_claim_uploader_role.go

* TODO: perform AssertPoolCanRun checks
* Try to claim uploader role without pool being funded
* Try to claim uploader role without being a staker
* Try to claim uploader role if the next uploader is not set yet
* Try to claim uploader role with non-existing valaccount
* Try to claim uploader role with wrong valaccount / TODO: clearify "wrong" valaccount
* TODO: claim uploader role if role is already claimed

*/

var _ = Describe("msg_server_claim_uploader_role.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create clean pool for every test case
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name: "Moontest",
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

	It("Try to claim uploader role without pool being funded", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ASSERT
		_, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeFalse())
	})

	It("Try to claim uploader role without being a staker", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_1,
			Staker:  i.STAKER_1,
			PoolId:  0,
		})

		// ASSERT
		_, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeFalse())
	})

	It("Try to claim uploader role if the next uploader is not set yet", func() {
		// ARRANGE
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

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ASSERT
		bundleProposal, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeTrue())

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
	})

	It("Try to claim uploader role with non existing valaccount", func() {
		// ARRANGE
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

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_1,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ASSERT
		_, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeFalse())
	})

	It("Try to claim uploader role with wrong valaccount", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Amount:  100 * i.KYVE,
		})

		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name: "Moontest2",
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_0,
		})

		// ACT
		s.RunTxBundlesError(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ASSERT
		_, found := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(found).To(BeFalse())
	})
})
