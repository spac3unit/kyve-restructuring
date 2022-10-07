package keeper_test

import (
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	"github.com/KYVENetwork/chain/x/delegation/types"
)

/*

TEST CASES - msg_server_redelegate.go

* Redelegate 1 KYVE to Bob
* Redelegate more than delegated
* Redelegate without delegation
* Redelegate to non-existent staker
* Exhaust all redelegation spells
* Expire redelegation spells

*/

var _ = Describe("Delegation - Redelegation", Ordered, func() {
	s := i.NewCleanChain()

	aliceSelfDelegation := 100 * i.KYVE
	bobSelfDelegation := 100 * i.KYVE

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

		_, stakerFound = s.App().StakersKeeper.GetStaker(s.Ctx(), i.BOB)
		Expect(stakerFound).To(BeTrue())

		s.CommitAfterSeconds(7)
	})

	AfterEach(func() {
		CheckAndContinueChainForOneMonth(&s)
	})

	It("Redelegate 1 KYVE to Bob", func() {

		// Arrange
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		aliceDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		bobDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(aliceDelegationBefore).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(bobDelegationBefore).To(Equal(bobSelfDelegation))

		// Act
		s.RunTxDelegatorSuccess(&types.MsgRedelegate{
			Creator:    i.DUMMY[0],
			FromStaker: i.ALICE,
			ToStaker:   i.BOB,
			Amount:     1 * i.KYVE,
		})
		s.CommitAfterSeconds(10)

		// Assert
		CheckAndContinueChainForOneMonth(&s)
		aliceDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegationBefore).To(Equal(aliceDelegationAfter + 1*i.KYVE))

		bobDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(bobDelegationBefore).To(Equal(bobDelegationAfter - 1*i.KYVE))
	})

	It("Redelegate more than delegated", func() {

		// Arrange
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		aliceDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		bobDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(aliceDelegationBefore).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(bobDelegationBefore).To(Equal(bobSelfDelegation))
		s.PerformValidityChecks()

		// Act
		s.RunTxDelegatorError(&types.MsgRedelegate{
			Creator:    i.DUMMY[0],
			FromStaker: i.ALICE,
			ToStaker:   i.BOB,
			Amount:     11 * i.KYVE,
		})
		s.CommitAfterSeconds(10)

		// Assert
		aliceDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegationBefore).To(Equal(aliceDelegationAfter))

		bobDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(bobDelegationBefore).To(Equal(bobDelegationAfter))
	})

	It("Redelegate without delegation", func() {

		// Arrange
		aliceDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		bobDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(aliceDelegationBefore).To(Equal(aliceSelfDelegation))
		Expect(bobDelegationBefore).To(Equal(bobSelfDelegation))
		s.PerformValidityChecks()

		// Act
		s.RunTxDelegatorError(&types.MsgRedelegate{
			Creator:    i.DUMMY[0],
			FromStaker: i.ALICE,
			ToStaker:   i.CHARLIE,
			Amount:     1 * i.KYVE,
		})

		// Assert
		aliceDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegationBefore).To(Equal(aliceDelegationAfter))

		bobDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(bobDelegationBefore).To(Equal(bobDelegationAfter))
	})

	It("Redelegate to non-existent staker", func() {

		// Arrange
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})
		aliceDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		bobDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(aliceDelegationBefore).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(bobDelegationBefore).To(Equal(bobSelfDelegation))
		s.PerformValidityChecks()

		// Act
		s.RunTxDelegatorError(&types.MsgRedelegate{
			Creator:    i.DUMMY[0],
			FromStaker: i.ALICE,
			ToStaker:   i.CHARLIE,
			Amount:     1 * i.KYVE,
		})

		// Assert
		aliceDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegationBefore).To(Equal(aliceDelegationAfter))

		bobDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(bobDelegationBefore).To(Equal(bobDelegationAfter))
	})

	It("Exhaust all redelegation spells", func() {

		// Arrange
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		aliceDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		bobDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(aliceDelegationBefore).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(bobDelegationBefore).To(Equal(bobSelfDelegation))
		s.PerformValidityChecks()

		// Act
		redelegationMessage := types.MsgRedelegate{
			Creator:    i.DUMMY[0],
			FromStaker: i.ALICE,
			ToStaker:   i.BOB,
			Amount:     1 * i.KYVE,
		}

		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)

		// Assert
		aliceDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegationAfter).To(Equal(aliceSelfDelegation + 5*i.KYVE))

		bobDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(bobDelegationAfter).To(Equal(bobSelfDelegation + 5*i.KYVE))

		// Expect to fail.
		// Now all redelegation spells are exhausted
		s.RunTxDelegatorError(&redelegationMessage)
	})

	It("Expire redelegation spells", func() {

		// Arrange
		s.RunTxDelegatorSuccess(&types.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  10 * i.KYVE,
		})
		Expect(s.GetBalanceFromAddress(i.DUMMY[0])).To(Equal(990 * i.KYVE))
		aliceDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		bobDelegationBefore := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(aliceDelegationBefore).To(Equal(aliceSelfDelegation + 10*i.KYVE))
		Expect(bobDelegationBefore).To(Equal(bobSelfDelegation))

		redelegationMessage := types.MsgRedelegate{
			Creator:    i.DUMMY[0],
			FromStaker: i.ALICE,
			ToStaker:   i.BOB,
			Amount:     1 * i.KYVE,
		}

		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(10)
		s.PerformValidityChecks()

		// Act
		s.CommitAfterSeconds(s.App().DelegationKeeper.RedelegationCooldown(s.Ctx()) - 50)
		s.CommitAfterSeconds(1)

		// One redelegation available
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(1)

		// Redelegations are now all used again
		s.RunTxDelegatorError(&redelegationMessage)
		s.PerformValidityChecks()

		// Act 2

		// Expire next two spells
		s.CommitAfterSeconds(25)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		// No two delegation within same block
		s.RunTxDelegatorError(&redelegationMessage)
		s.CommitAfterSeconds(1)
		s.RunTxDelegatorSuccess(&redelegationMessage)
		s.CommitAfterSeconds(1)

		// Assert
		aliceDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.ALICE)
		Expect(aliceDelegationAfter).To(Equal(aliceSelfDelegation + 2*i.KYVE))

		bobDelegationAfter := s.App().DelegationKeeper.GetDelegationAmount(s.Ctx(), i.BOB)
		Expect(bobDelegationAfter).To(Equal(bobSelfDelegation + 8*i.KYVE))

		// Expect to fail.
		// Now all redelegation spells are exhausted
		s.RunTxDelegatorError(&redelegationMessage)
	})

})
