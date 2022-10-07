package keeper_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_leave_pool.go

* Leave a pool a staker has just joined as the first one
* // TODO: leave a pool multiple other stakers have joined previously
* Leave one of multiple pools a staker has previously joined
* Try to leave a pool again
* // TODO: try to leave a pool a staker has never joined

*/

var _ = Describe("msg_server_leave_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create pool
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name: "Moontest",
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		// create staker
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})

		// join pool
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Leave a pool a staker has just joined as the first one", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  0,
		})
		s.PerformValidityChecks()

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_0))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeTrue())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))

		s.PerformValidityChecks()

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.UnbondingStakingTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccountsOfStaker = s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(BeEmpty())

		_, found = s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeFalse())

		valaccountsOfPool = s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(BeEmpty())

		totalStakeOfPool = s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(BeZero())
	})

	It("Try to leave a pool again", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  0,
		})
		s.PerformValidityChecks()

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  0,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccountsOfStaker).To(HaveLen(1))

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.UnbondingStakingTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccountsOfStaker = s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccountsOfStaker).To(BeEmpty())
	})

	It("Try to leave one of multiple pools a staker has joined", func() {
		// ARRANGE
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name: "Moontest",
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_1,
		})
		s.PerformValidityChecks()

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgLeavePool{
			Creator: i.STAKER_0,
			PoolId:  1,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(2))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(Equal(uint64(1)))
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_1))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeTrue())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 1)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 1)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))

		// wait for leave pool
		s.CommitAfterSeconds(s.App().StakersKeeper.UnbondingStakingTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		valaccountsOfStaker = s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		_, found = s.App().StakersKeeper.GetValaccount(s.Ctx(), 1, i.STAKER_0)

		Expect(found).To(BeFalse())

		valaccountsOfPool = s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 1)

		Expect(valaccountsOfPool).To(BeEmpty())

		totalStakeOfPool = s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 1)
		Expect(totalStakeOfPool).To(BeZero())
	})
})
