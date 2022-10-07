package keeper_test

import (
	delegationtypes "github.com/KYVENetwork/chain/x/delegation/types"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
)

/*

TEST CASES - logic_gov_bonding_test.go

* Simple Staking
* Multiple Staking
* Multiple Staking + Delegation
* Multiple Staking + Multiple Delegation

*/

var _ = Describe("Delegation Gov Logic", Ordered, func() {
	s := i.NewCleanChain()
	aliceAcc, _ := sdk.AccAddressFromBech32(i.ALICE)
	bobAcc, _ := sdk.AccAddressFromBech32(i.BOB)
	charlieAcc, _ := sdk.AccAddressFromBech32(i.CHARLIE)

	BeforeEach(func() {
		s = i.NewCleanChain()

		Expect(s.App().DelegationKeeper.TotalProtocolBonding(s.Ctx())).To(Equal(sdk.NewInt(0)))
	})

	AfterEach(func() {
		CheckAndContinueChainForOneMonth(&s)
	})

	It("Simple Staking", func() {

		// Arrange
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.ALICE,
			Amount:  100 * i.KYVE,
		})

		// Assert
		Expect(s.App().DelegationKeeper.TotalProtocolBonding(s.Ctx())).To(Equal(sdk.NewInt(int64(100 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), aliceAcc)).To(Equal(sdk.NewInt(int64(100 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), bobAcc)).To(Equal(sdk.NewInt(int64(0 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), charlieAcc)).To(Equal(sdk.NewInt(int64(0 * i.KYVE))))
	})

	It("Multiple Staking", func() {
		Expect(s.App().DelegationKeeper.TotalProtocolBonding(s.Ctx())).To(Equal(sdk.NewInt(0)))

		// Arrange
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.ALICE,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.BOB,
			Amount:  50 * i.KYVE,
		})

		// Assert
		Expect(s.App().DelegationKeeper.TotalProtocolBonding(s.Ctx())).To(Equal(sdk.NewInt(int64(150 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), aliceAcc)).To(Equal(sdk.NewInt(int64(100 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), bobAcc)).To(Equal(sdk.NewInt(int64(50 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), charlieAcc)).To(Equal(sdk.NewInt(int64(0 * i.KYVE))))
	})

	It("Multiple Staking + Delegation", func() {
		Expect(s.App().DelegationKeeper.TotalProtocolBonding(s.Ctx())).To(Equal(sdk.NewInt(0)))

		// Arrange
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.ALICE,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.BOB,
			Amount:  50 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.CHARLIE,
			Staker:  i.ALICE,
			Amount:  80 * i.KYVE,
		})

		// Assert
		Expect(s.App().DelegationKeeper.TotalProtocolBonding(s.Ctx())).To(Equal(sdk.NewInt(int64((100 + 50 + 80) * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), aliceAcc)).To(Equal(sdk.NewInt(int64(100 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), bobAcc)).To(Equal(sdk.NewInt(int64(50 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), charlieAcc)).To(Equal(sdk.NewInt(int64(80 * i.KYVE))))
	})

	It("Multiple Staking + Multiple Delegation", func() {
		Expect(s.App().DelegationKeeper.TotalProtocolBonding(s.Ctx())).To(Equal(sdk.NewInt(0)))

		// Arrange
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.ALICE,
			Amount:  100 * i.KYVE,
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.BOB,
			Amount:  50 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.CHARLIE,
			Staker:  i.ALICE,
			Amount:  80 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.DUMMY[0],
			Staker:  i.ALICE,
			Amount:  60 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.DUMMY[1],
			Staker:  i.BOB,
			Amount:  40 * i.KYVE,
		})

		s.RunTxDelegatorSuccess(&delegationtypes.MsgDelegate{
			Creator: i.DUMMY[2],
			Staker:  i.BOB,
			Amount:  75 * i.KYVE,
		})

		dummy0Acc, _ := sdk.AccAddressFromBech32(i.DUMMY[0])
		dummy1Acc, _ := sdk.AccAddressFromBech32(i.DUMMY[1])
		dummy2Acc, _ := sdk.AccAddressFromBech32(i.DUMMY[2])

		// Assert
		Expect(s.App().DelegationKeeper.TotalProtocolBonding(s.Ctx())).To(Equal(sdk.NewInt(int64((100 + 50 + 80 + 60 + 40 + 75) * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), aliceAcc)).To(Equal(sdk.NewInt(int64(100 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), bobAcc)).To(Equal(sdk.NewInt(int64(50 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), charlieAcc)).To(Equal(sdk.NewInt(int64(80 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), dummy0Acc)).To(Equal(sdk.NewInt(int64(60 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), dummy1Acc)).To(Equal(sdk.NewInt(int64(40 * i.KYVE))))
		Expect(s.App().DelegationKeeper.GetBondingOfAddress(s.Ctx(), dummy2Acc)).To(Equal(sdk.NewInt(int64(75 * i.KYVE))))
	})

})
