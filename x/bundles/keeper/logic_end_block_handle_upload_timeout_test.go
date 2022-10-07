package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - logic_end_block_handle_upload_timeout.go

* First staker who joins gets automatically chosen as next uploader
* Next uploader gets removed due to pool upgrading
* Next uploader gets removed due to pool being paused
* Next uploader gets removed due to pool having no funds
* Next uploader gets removed due to pool not reaching min stake
* Staker is next uploader of genesis bundle and upload interval and timeout does not pass
* Staker is next uploader of genesis bundle and upload timeout does not pass but upload interval passes
* Staker is next uploader of genesis bundle and upload timeout does pass together with upload interval
* Staker is next uploader of bundle proposal and upload interval does not pass
* Staker is next uploader of bundle proposal and upload timeout does not pass
* Staker is next uploader of bundle proposal and upload timeout passes
* A bundle proposal with no quorum does not reach the upload interval
* A bundle proposal with no quorum does reach the upload interval
* Staker who just left the pool is next uploader of bundle proposal and upload timeout passes
* TODO: test with multiple pools

*/

var _ = Describe("logic_end_block_handle_upload_timeout.go", Ordered, func() {
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
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("First staker who joins gets automatically chosen as next uploader", func() {
		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Next uploader gets removed due to pool upgrading", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		pool.UpgradePlan = &pooltypes.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{}",
			ScheduledAt: uint64(s.Ctx().BlockTime().Unix()),
			Duration:    3600,
		}

		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Next uploader gets removed due to pool being paused", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		pool.Paused = true

		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Next uploader gets removed due to pool having no funds", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.RunTxPoolSuccess(&pooltypes.MsgDefundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Next uploader gets removed due to pool not reaching min stake", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgUndelegate{
			Creator: i.STAKER_0,
			Staker:  i.STAKER_0,
			Amount:  50 * i.KYVE,
		})

		s.CommitAfterSeconds(s.App().DelegationKeeper.UnbondingDelegationTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(50 * i.KYVE))
	})

	It("Staker is next uploader of genesis bundle and upload interval and timeout does not pass", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Staker is next uploader of genesis bundle and upload timeout does not pass but upload interval passes", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ACT
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Staker is next uploader of genesis bundle and upload timeout does pass together with upload interval", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.UploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(BeEmpty())

		// check if next uploader got removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(BeEmpty())

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())

		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(BeZero())

		// check if next uploader got slashed
		slashAmountRatio, _ := sdk.NewDecFromStr(s.App().StakersKeeper.TimeoutSlash(s.Ctx()))
		expectedBalance := 100*i.KYVE - uint64(sdk.NewDec(int64(100*i.KYVE)).Mul(slashAmountRatio).RoundInt64())

		Expect(expectedBalance).To(Equal(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))
	})

	It("Staker is next uploader of bundle proposal and upload interval does not pass", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

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
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Staker is next uploader of bundle proposal and upload timeout does not pass", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

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
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Staker is next uploader of bundle proposal and upload timeout passes", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

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
		s.CommitAfterSeconds(s.App().BundlesKeeper.UploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(BeEmpty())
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		// check if next uploader got removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(BeEmpty())

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())

		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(BeZero())

		// check if next uploader got slashed
		slashAmountRatio, _ := sdk.NewDecFromStr(s.App().StakersKeeper.TimeoutSlash(s.Ctx()))
		expectedBalance := 100*i.KYVE - uint64(sdk.NewDec(int64(100*i.KYVE)).Mul(slashAmountRatio).RoundInt64())

		Expect(expectedBalance).To(Equal(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))
	})

	It("A bundle proposal with no quorum does not reach the upload interval", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

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

		// ACT
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("A bundle proposal with no quorum does reach the upload interval", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

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

		// ACT
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)

		Expect(bundleProposal.StorageId).To(BeEmpty())
		Expect(bundleProposal.Uploader).To(BeEmpty())
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_0))
		Expect(bundleProposal.ByteSize).To(BeZero())
		Expect(bundleProposal.ToHeight).To(BeZero())
		Expect(bundleProposal.ToKey).To(BeEmpty())
		Expect(bundleProposal.ToValue).To(BeEmpty())
		Expect(bundleProposal.BundleHash).To(BeEmpty())
		Expect(bundleProposal.VotersValid).To(BeEmpty())
		Expect(bundleProposal.VotersInvalid).To(BeEmpty())
		Expect(bundleProposal.VotersAbstain).To(BeEmpty())

		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(2))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(100 * i.KYVE))
	})

	It("Staker who just left the pool is next uploader of bundle proposal and upload timeout passes", func() {
		// ARRANGE
		s.RunTxBundlesSuccess(&bundletypes.MsgClaimUploaderRole{
			Creator: i.VALADDRESS_0,
			Staker:  i.STAKER_0,
			PoolId:  0,
		})

		s.CommitAfterSeconds(60)

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

		// remove valaccount directory from pool
		s.App().StakersKeeper.RemoveValaccountFromPool(s.Ctx(), 0, i.STAKER_0)

		// ACT
		s.CommitAfterSeconds(s.App().BundlesKeeper.UploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(60)
		s.CommitAfterSeconds(1)

		// ASSERT
		bundleProposal, _ := s.App().BundlesKeeper.GetBundleProposal(s.Ctx(), 0)
		Expect(bundleProposal.NextUploader).To(Equal(i.STAKER_1))
		Expect(bundleProposal.StorageId).To(Equal("y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI"))

		// check if next uploader got removed from pool
		poolStakers := s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)
		Expect(poolStakers).To(HaveLen(1))

		_, found := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(found).To(BeTrue())

		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal(100 * i.KYVE))

		// check if next uploader got slashed
		expectedBalance := 100 * i.KYVE

		Expect(expectedBalance).To(Equal(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)))
	})
})
