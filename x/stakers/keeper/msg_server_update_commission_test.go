package keeper_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	stakerstypes "github.com/KYVENetwork/chain/x/stakers/types"
)

/*

TEST CASES - msg_server_update_commission.go

* Get the default commission from a newly created staker
* Update commission to 50% from previously default commission
* Update commission to 0% from previously default commission
* Update commission to 100% from previously default commission
* Update commission with an invalid number from previously default commission
* Update commission with a negative number from previously default commission
* Update commission with a too high number from previously default commission
* Update commission multiple times during the commission change time
* Update commission multiple times during the commission change time with the same value
* // TODO: commission should reset if staker unstakes everything and stakes again
* // TODO: update commission with multiple stakers

*/

var _ = Describe("msg_server_update_commission.go", Ordered, func() {
	s := i.NewCleanChain()

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create staker
		s.RunTxStakersSuccess(&stakerstypes.MsgCreateStaker{
			Creator: i.STAKER_0,
			Amount:  100 * i.KYVE,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Get the default commission from a newly created staker", func() {
		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))
	})

	It("Update commission to 50% from previously default commission", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "0.5",
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.CommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal("0.5"))
	})

	It("Update commission to 0% from previously default commission", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "0",
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.CommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal("0"))
	})

	It("Update commission to 100% from previously default commission", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "1",
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.CommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal("1"))
	})

	It("Update commission with an invalid number from previously default commission", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "teset",
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))
	})

	It("Update commission with a negative number from previously default commission", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "-0.5",
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))
	})

	It("Update commission with a too high number from previously default commission", func() {
		// ACT
		s.RunTxStakersError(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "2",
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))
	})

	It("Update commission multiple times during the commission change time", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "0.5",
		})
		s.PerformValidityChecks()

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "0.2",
		})
		s.PerformValidityChecks()

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "0.3",
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.CommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal("0.3"))
	})

	It("Update commission multiple times during the commission change time with the same value", func() {
		// ACT
		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "0.5",
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: "0.2",
		})

		s.RunTxStakersSuccess(&stakerstypes.MsgUpdateCommission{
			Creator:    i.STAKER_0,
			Commission: stakerstypes.DefaultCommission,
		})
		s.PerformValidityChecks()

		// ASSERT
		staker, _ := s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))

		// wait for update
		s.CommitAfterSeconds(s.App().StakersKeeper.CommissionChangeTime(s.Ctx()))
		s.CommitAfterSeconds(1)

		staker, _ = s.App().StakersKeeper.GetStaker(s.Ctx(), i.STAKER_0)
		Expect(staker.Commission).To(Equal(stakerstypes.DefaultCommission))
	})
})
