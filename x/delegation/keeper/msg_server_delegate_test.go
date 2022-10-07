package keeper_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/delegation/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_delegate_test.go

* Delegate 10 KYVE to ALICE
* Try delegating to non-existent staker
* Delegate more than available
* Payout delegators
* Don't pay out rewards twice

*/

var _ = Describe("Delegation - Delegation", Ordered, func() {
	s := i.NewCleanChain()
	const aliceSelfDelegation uint64 = 100 * i.KYVE
	const bobSelfDelegation uint64 = 200 * i.KYVE

	BeforeEach(func() {
		s = i.NewCleanChain()

		CreateFundedPool(&s)

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.ALICE,
			Amount:  aliceSelfDelegation,
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.BOB,
			Amount:  bobSelfDelegation,
		})

		_, stakerFound := s.App().StakersKeeper.GetStaker(s.Ctx(), i.ALICE)
		Expect(stakerFound).To(BeTrue())

		s.CommitAfterSeconds(7)
	})

	AfterEach(func() {
		CheckAndContinueChainForOneMonth(&s)
	})

	It("Delegate 10 KYVE to ALICE", func() {

		// Arrange
		bobBalance := s.GetBalanceFromAddress(i.BOB)

		// Act
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		// Assert
		CheckAndContinueChainForOneMonth(&s)
		bobBalanceAfter := s.GetBalanceFromAddress(i.BOB)
		Expect(bobBalanceAfter).To(Equal(bobBalance - 10*i.KYVE))

		aliceDelegation := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegation).To(Equal(10*i.KYVE + aliceSelfDelegation))
	})

	It("Try delegating to non-existent staker", func() {

		// Arrange
		bobBalance := s.GetBalanceFromAddress(i.BOB)
		s.PerformValidityChecks()

		// Act
		s.RunTxDelegatorError(&types.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.CHARLIE,
			Amount:  10 * i.KYVE,
		})

		// Assert
		Expect(s.GetBalanceFromAddress(i.BOB)).To(Equal(bobBalance))

		aliceDelegation := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegation).To(Equal(aliceSelfDelegation))
	})

	It("Delegate more than available", func() {

		// Arrange
		bobBalance := s.GetBalanceFromAddress(i.BOB)
		aliceDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		s.PerformValidityChecks()

		// Act
		_, delegateErr := s.RunTxDelegator(&types.MsgDelegate{
			Creator: i.BOB,
			Staker:  i.ALICE,
			Amount:  2000 * i.KYVE,
		})

		// Assert
		Expect(delegateErr).ToNot(BeNil())

		aliceDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegationBefore).To(Equal(aliceDelegationAfter))

		bobBalanceAfter := s.GetBalanceFromAddress(i.BOB)
		Expect(bobBalanceAfter).To(Equal(bobBalance))
	})

	It("Payout delegators", func() {

		// Arrange
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  100 * i.KYVE,
		})
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  209 * i.KYVE,
		})
		poolModuleBalance := s.GetBalanceFromModule(pooltypes.ModuleName)
		Expect(poolModuleBalance).To(Equal(50 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(BeZero())
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(BeZero())
		s.PerformValidityChecks()

		// Act
		PayoutRewards(&s, i.ALICE, 10*i.KYVE)

		// Name    amount   shares
		// Alice:   100		100/(409) * 10 * 1e9 = 2.444.987.775
		// Dummy0:  100		100/(409) * 10 * 1e9 = 2.444.987.775
		// Dummy1:  209		209/(409) * 10 * 1e9 = 5.110.024.449
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.ALICE)).To(Equal(uint64(2_444_987_775)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(uint64(2_444_987_775)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(uint64(5_110_024_449)))

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})

		// Assert
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.ALICE)).To(Equal(uint64(2_444_987_775)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(uint64(0)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(uint64(5_110_024_449)))

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(uint64(900*i.KYVE + 2_444_987_775)))
		Expect(s.GetBalanceFromModule(pooltypes.ModuleName)).To(Equal(40 * i.KYVE))
		Expect(s.GetBalanceFromModule(types.ModuleName)).To(Equal((200+409)*i.KYVE + uint64(2_444_987_775+5_110_024_449+1)))
	})

	It("Don't pay out rewards twice", func() {

		// Arrange
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  100 * i.KYVE,
		})
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  200 * i.KYVE,
		})
		poolModuleBalance := s.GetBalanceFromModule(pooltypes.ModuleName)
		Expect(poolModuleBalance).To(Equal(50 * i.KYVE))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(BeZero())
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(BeZero())

		PayoutRewards(&s, i.ALICE, 10*i.KYVE)

		// Alice: 100
		// Dummy0: 100
		// Dummy1: 200
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(uint64(2_500_000_000)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(uint64(5_000_000_000)))
		s.PerformValidityChecks()

		// Act
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})

		// Assert
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(uint64(0)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(uint64(5_000_000_000)))

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(uint64(900*i.KYVE + 2_500_000_000)))
		Expect(s.GetBalanceFromModule(pooltypes.ModuleName)).To(Equal(40 * i.KYVE))
		Expect(s.GetBalanceFromModule(types.ModuleName)).To(Equal(600*i.KYVE + 7_500_000_000))
	})

})
