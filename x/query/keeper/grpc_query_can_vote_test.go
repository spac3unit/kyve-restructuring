package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	querytypes "github.com/KYVENetwork/chain/x/query/types"
	"github.com/KYVENetwork/chain/x/registry/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - grpc_query_can_vote.go

* Call can vote if pool does not exist
* Call can vote if pool is currently upgrading
* Call can vote if pool is paused
* Call can vote if pool is out of funds
* Call can vote if pool has not reached the minimum stake
* Call can vote with a valaccount which does not exist
* Call can vote if current bundle was dropped
* Call can vote with a different storage id than the current one
* Call can vote if voter has already voted valid
* Call can vote if voter has already voted invalid
* Call can vote if voter has already voted abstain
* Call can vote on an active pool with a data bundle with valid args

*/

var _ = Describe("grpc_query_can_vote.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		s = i.NewCleanChain()

		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "Moontest",
			MinStake:       200 * i.KYVE,
			UploadInterval: 60,
			MaxBundleSize:  100,
			Protocol:       &pooltypes.Protocol{},
			UpgradePlan:    &pooltypes.UpgradePlan{},
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
			Amount:     0,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgCreateStaker{
			Creator: i.STAKER_1,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakertypes.MsgJoinPool{
			Creator:    i.STAKER_1,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
			Amount:     0,
		})

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
			StorageId:  "test_storage_id",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   100,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "test_hash",
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Call can vote if pool does not exist", func() {
		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    1,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(sdkErrors.Wrapf(sdkErrors.ErrNotFound, types.ErrPoolNotFound.Error(), 1).Error()))

		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    1,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if pool is currently upgrading", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.UpgradePlan = &pooltypes.UpgradePlan{
			Version:     "1.0.0",
			Binaries:    "{}",
			ScheduledAt: 100,
			Duration:    3600,
		}

		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrPoolCurrentlyUpgrading.Error()))

		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if pool is paused", func() {
		// ARRANGE
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		pool.Paused = true

		s.App().PoolKeeper.SetPool(s.Ctx(), pool)

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrPoolPaused.Error()))

		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if pool is out of funds", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgDefundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrPoolOutOfFunds.Error()))

		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if pool has not reached the minimum stake", func() {
		// ARRANGE
		s.RunTxDelegatorSuccess(&delegationtypes.MsgUndelegate{
			Creator: i.STAKER_0,
			Staker:  i.STAKER_0,
			Amount:  50 * i.KYVE,
		})

		// wait for unbonding
		s.CommitAfterSeconds(s.App().DelegationKeeper.UnbondingDelegationTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrMinStakeNotReached.Error()))

		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote with a valaccount which does not exist", func() {
		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_0,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(stakertypes.ErrValaccountUnauthorized.Error()))

		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_0,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if previous bundle was dropped", func() {
		// ARRANGE
		// wait for timeout so bundle gets dropped
		s.CommitAfterSeconds(s.App().BundlesKeeper.UploadTimeout(s.Ctx()))
		s.CommitAfterSeconds(1)

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrBundleDropped.Error()))

		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote with a different storage id than the current one", func() {
		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "another_test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrInvalidStorageId.Error()))

		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "another_test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if voter has already voted valid", func() {
		// ARRANGE
		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).To(BeNil())

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrAlreadyVotedValid.Error()))

		_, txErr = s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if voter has already voted invalid", func() {
		// ARRANGE
		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_NO,
		})

		Expect(txErr).To(BeNil())

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeFalse())
		Expect(canVote.Reason).To(Equal(bundletypes.ErrAlreadyVotedInvalid.Error()))

		_, txErr = s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_NO,
		})

		Expect(txErr).NotTo(BeNil())
		Expect(txErr.Error()).To(Equal(canVote.Reason))
	})

	It("Call can vote if voter has already voted abstain", func() {
		// ARRANGE
		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_ABSTAIN,
		})

		Expect(txErr).To(BeNil())

		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeTrue())
		Expect(canVote.Reason).To(Equal("KYVE_VOTE_NO_ABSTAIN_ALLOWED"))

		_, txErr = s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).To(BeNil())
	})

	It("Call can vote on an active pool with a data bundle with valid args", func() {
		// ACT
		canVote, err := s.App().QueryKeeper.CanVote(sdk.WrapSDKContext(s.Ctx()), &querytypes.QueryCanVoteRequest{
			PoolId:    0,
			Staker:    i.STAKER_1,
			Voter:     i.VALADDRESS_1,
			StorageId: "test_storage_id",
		})

		// ASSERT
		Expect(err).To(BeNil())

		Expect(canVote.Possible).To(BeTrue())
		Expect(canVote.Reason).To(BeEmpty())

		_, txErr := s.RunTxBundles(&bundletypes.MsgVoteBundleProposal{
			Creator:   i.VALADDRESS_1,
			Staker:    i.STAKER_1,
			PoolId:    0,
			StorageId: "test_storage_id",
			Vote:      bundletypes.VOTE_TYPE_YES,
		})

		Expect(txErr).To(BeNil())
	})
})
