package keeper_test

import (
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/delegation/types"
)

/*

TEST CASES - msg_server_withdraw_rewards_test.go

* Split rewards; test rounding
* Withdraw from not delegator
* Test invalid Payout

TODO test withdraw with multiple slashes

*/

var _ = Describe("Delegation - Withdraw Rewards", Ordered, func() {
	s := i.NewCleanChain()

	aliceSelfDelegation := 0 * i.KYVE
	bobSelfDelegation := 0 * i.KYVE

	BeforeEach(func() {
		s = i.NewCleanChain()

		CreateFundedPool(&s)

		// Stake
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{Creator: i.ALICE, Amount: aliceSelfDelegation})
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{Creator: i.BOB, Amount: bobSelfDelegation})

		_, stakerFound := s.App().StakersKeeper.GetStaker(s.Ctx(), i.ALICE)
		Expect(stakerFound).To(BeTrue())

		_, stakerFound = s.App().StakersKeeper.GetStaker(s.Ctx(), i.BOB)
		Expect(stakerFound).To(BeTrue())
	})

	AfterEach(func() {
		CheckAndContinueChainForOneMonth(&s)
	})

	It("Split rewards; test rounding", func() {

		// Arrange
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(990 * i.KYVE))
		Expect(s.GetBalanceFromAddress(i.DUMMY[2])).To(Equal(990 * i.KYVE))

		Expect(s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)).To(Equal(aliceSelfDelegation + 30*i.KYVE))

		delegationModuleBalanceBefore := s.GetBalanceFromModule(types.ModuleName)
		poolModuleBalanceBefore := s.GetBalanceFromModule(pooltypes.ModuleName)
		s.PerformValidityChecks()

		// Act
		// Alice: 100
		// Dummy0: 10
		// Dummy1: 0
		PayoutRewards(&s, i.ALICE, 20*i.KYVE)

		delegationModuleBalanceAfter := s.GetBalanceFromModule(types.ModuleName)
		poolModuleBalanceAfter := s.GetBalanceFromModule(pooltypes.ModuleName)

		Expect(delegationModuleBalanceAfter).To(Equal(delegationModuleBalanceBefore + 20*i.KYVE))
		Expect(poolModuleBalanceAfter).To(Equal(poolModuleBalanceBefore - 20*i.KYVE))

		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[0])).To(Equal(uint64(6666666666)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[1])).To(Equal(uint64(6666666666)))
		Expect(s.App().DelegationKeeper.GetOutstandingRewards(s.Ctx(), i.ALICE, i.DUMMY[2])).To(Equal(uint64(6666666666)))

		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
		})
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[1],
			Staker:  i.ALICE,
		})
		s.RunTxDelegatorSuccess(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[2],
			Staker:  i.ALICE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(uint64(996666666666)))
		Expect(s.GetBalanceFromAddress(i.DUMMY[1])).To(Equal(uint64(996666666666)))
		Expect(s.GetBalanceFromAddress(i.DUMMY[2])).To(Equal(uint64(996666666666)))

		Expect(s.GetBalanceFromModule(types.ModuleName)).To(Equal(uint64(30000000002)))
	})

	It("Withdraw from not delegator", func() {

		balanceDummy1Before := s.GetBalanceFromAddress(i.DUMMY[0])
		balanceCharlieBefore := s.GetBalanceFromAddress(i.CHARLIE)
		balanceAliceBefore := s.GetBalanceFromAddress(i.ALICE)
		delegationBalance := s.GetBalanceFromModule(types.ModuleName)
		s.PerformValidityChecks()

		s.RunTxDelegatorError(&types.MsgWithdrawRewards{
			Creator: i.CHARLIE,
			Staker:  i.ALICE,
		})

		s.RunTxDelegatorError(&types.MsgWithdrawRewards{
			Creator: i.DUMMY[0],
			Staker:  i.CHARLIE,
		})

		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(balanceDummy1Before))
		Expect(s.GetBalanceFromAddress(i.CHARLIE)).To(Equal(balanceCharlieBefore))
		Expect(s.GetBalanceFromAddress(i.ALICE)).To(Equal(balanceAliceBefore))
		Expect(s.GetBalanceFromModule(types.ModuleName)).To(Equal(delegationBalance))

	})

	It("Test invalid Payout", func() {
		forkedCtx, _ := s.Ctx().CacheContext()
		success := s.App().DelegationKeeper.PayoutRewards(forkedCtx, i.ALICE, 20000*i.KYVE, pooltypes.ModuleName)
		Expect(success).To(BeFalse())

		success = s.App().DelegationKeeper.PayoutRewards(s.Ctx(), i.DUMMY[20], 0*i.KYVE, pooltypes.ModuleName)
		Expect(success).To(BeFalse())
	})

})
