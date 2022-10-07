package keeper_test

import (
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_join_pool.go

* Test if a newly created staker is participating in no pools yet
* Join the first pool as the first staker to a newly created pool
* // TODO: join a pool where other stakers have already joined
* Self-Delegate more KYVE after joining a pool
* Try to join the same pool with the same valaddress again
* Try to join the same pool with a different valaddress
* Try to join another pool with the same valaddress again
* Try to join another pool with a different valaddress
* Join a pool with a valaddress which does not exist on chain yet
* Join a pool with a valaddress which does not exist on chain yet and send 0 funds
* Join a pool with an invalid valaddress
* Join a pool and fund the valaddress with more KYVE than available in balance
* Kick out lowest staker by joining a full pool
* // TODO: fail to kick out lowest staker because not enough stake
* Kick out lowest staker with respect to stake + delegation
* // TODO: fail to kick out lowest staker because not enough delegation

*/

var _ = Describe("msg_server_join_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalanceStaker0 := uint64(0)
	initialBalanceValaddress0 := uint64(0)

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

		initialBalanceStaker0 = s.GetBalanceFromAddress(i.STAKER_0)
		initialBalanceValaddress0 = s.GetBalanceFromAddress(i.VALADDRESS_0)
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Test if a newly created staker is participating in no pools yet", func() {
		// ASSERT
		valaccounts := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccounts).To(HaveLen(0))
	})

	It("Join the first pool as the first staker to a newly created pool", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterValaddress0 := s.GetBalanceFromAddress(i.VALADDRESS_0)

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(100 * i.KYVE))
		Expect(balanceAfterValaddress0 - initialBalanceValaddress0).To(Equal(100 * i.KYVE))

		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_0))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Self-Delegate more KYVE after joining a pool", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     100 * i.KYVE,
		})

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		// ACT
		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.STAKER_0,
			Staker:  i.STAKER_0,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal(i.VALADDRESS_0))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool = s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)

		Expect(totalStakeOfPool).To(Equal(150 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Try to join the same pool with the same valaddress again", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     100 * i.KYVE,
		})

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))
	})

	It("Try to join the same pool with a different valaddress", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     100 * i.KYVE,
		})

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_1,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))
	})

	It("Try to join another pool with the same valaddress again", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     100 * i.KYVE,
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

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_0,
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccountsOfStaker).To(HaveLen(1))
	})

	It("Try to join another pool with a different valaddress", func() {
		// ARRANGE
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     100 * i.KYVE,
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

		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     1,
			Valaddress: i.VALADDRESS_1,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)
		Expect(valaccountsOfStaker).To(HaveLen(2))
	})

	It("Join a pool with a valaddress which does not exist on chain yet", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: "kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x",
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterUnknown := s.GetBalanceFromAddress("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x")

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(Equal(100 * i.KYVE))
		Expect(balanceAfterUnknown).To(Equal(100 * i.KYVE))

		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x"))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Join a pool with a valaddress which does not exist on chain yet and send 0 funds", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: "kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x",
			Amount:     0 * i.KYVE,
		})

		// ASSERT
		balanceAfterStaker0 := s.GetBalanceFromAddress(i.STAKER_0)
		balanceAfterUnknown := s.GetBalanceFromAddress("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x")

		Expect(initialBalanceStaker0 - balanceAfterStaker0).To(BeZero())
		Expect(balanceAfterUnknown).To(BeZero())

		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(HaveLen(1))

		valaccount, found := s.App().StakersKeeper.GetValaccount(s.Ctx(), 0, i.STAKER_0)

		Expect(found).To(BeTrue())

		Expect(valaccount.Staker).To(Equal(i.STAKER_0))
		Expect(valaccount.PoolId).To(BeZero())
		Expect(valaccount.Valaddress).To(Equal("kyve1dx0nvx7y9d44jvr2dr6r2p636jea3f9827rn0x"))
		Expect(valaccount.Points).To(BeZero())
		Expect(valaccount.IsLeaving).To(BeFalse())

		valaccountsOfPool := s.App().StakersKeeper.GetAllValaccountsOfPool(s.Ctx(), 0)

		Expect(valaccountsOfPool).To(HaveLen(1))

		totalStakeOfPool := s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)
		Expect(totalStakeOfPool).To(Equal(100 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.STAKER_0)).To(Equal(totalStakeOfPool))
		Expect(s.App().DelegationKeeper.GetDelegationAmountOfDelegator(s.Ctx(), i.STAKER_0, i.STAKER_0)).To(Equal(totalStakeOfPool))
	})

	It("Join a pool with an invalid valaddress", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: "invalid_valaddress",
			Amount:     100 * i.KYVE,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.STAKER_0)

		Expect(valaccountsOfStaker).To(BeEmpty())
	})

	It("Join a pool and fund the valaddress with more KYVE than available in balance", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: "invalid_valaddress",
			Amount:     initialBalanceStaker0 + 1,
		})

		// ASSERT
		valaccountsOfStaker := s.App().StakersKeeper.GetValaccountsFromStaker(s.Ctx(), i.ALICE)

		Expect(valaccountsOfStaker).To(BeEmpty())
	})

	It("Kick out lowest staker by joining a full pool", func() {

		// Arrange
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     1,
		})

		for k := 0; k < 49; k++ {
			s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
				Creator: i.DUMMY[k],
				Amount:  150 * i.KYVE,
			})
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:    i.DUMMY[k],
				PoolId:     0,
				Valaddress: i.VALDUMMY[k],
				Amount:     1,
			})
		}

		// STAKER_0 is lowest staker and all stakers are full now.
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.BOB,
			Amount:  150 * i.KYVE,
		})

		// Act
		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.BOB,
			PoolId:     0,
			Valaddress: i.BOB,
			Amount:     1,
		})

		// Assert
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 150) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).ToNot(ContainElement(i.STAKER_0))
	})

	It("Kick out lowest staker with respect to stake + delegation", func() {
		// ARRANGE
		Expect(stakerstypes.MaxStakers).To(Equal(50))

		s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
			Creator:    i.STAKER_0,
			PoolId:     0,
			Valaddress: i.VALADDRESS_0,
			Amount:     1 * i.KYVE,
		})

		for k := 0; k < 49; k++ {
			s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
				Creator: i.DUMMY[k],
				Amount:  150 * i.KYVE,
			})
			s.RunTxStakersSuccess(&stakerstypes.MsgJoinPool{
				Creator:    i.DUMMY[k],
				PoolId:     0,
				Valaddress: i.VALDUMMY[k],
				Amount:     1 * i.KYVE,
			})
		}

		// Alice is lowest staker and all stakers are full now.
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 100) * i.KYVE))

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.BOB,
			Amount:  150 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.ALICE,
			Staker:  i.STAKER_0,
			Amount:  150 * i.KYVE,
		}) // Staker0 has now 250 delegation

		// ACT
		s.RunTxStakersError(&stakerstypes.MsgJoinPool{
			Creator:    i.BOB,
			PoolId:     0,
			Valaddress: i.BOB,
			Amount:     1,
		})

		// ASSERT
		Expect(s.App().DelegationKeeper.GetDelegationOfPool(s.Ctx(), 0)).To(Equal((150*49 + 250) * i.KYVE))
		Expect(s.App().StakersKeeper.GetAllStakerAddressesOfPool(s.Ctx(), 0)).To(ContainElement(i.STAKER_0))
	})
})
