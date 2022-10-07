package keeper_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	i "github.com/KYVENetwork/chain/testutil/integration"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
)

/*

TEST CASES - msg_server_defund_pool.go

* Defund 50 KYVE from a funder who has previously funded 100 KYVE
* Try to defund more than actually funded
* Defund full funding amount from a funder who has previously funder 100 KYVE
* Defund as highest funder 75 KYVE in order to be the lowest funder afterwards

*/

var _ = Describe("msg_server_defund_pool.go", Ordered, func() {
	s := i.NewCleanChain()

	initialBalance := s.GetBalanceFromAddress(i.ALICE)

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

		// fund pool
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})
	})

	AfterEach(func() {
		s.PerformValidityChecks()
	})

	It("Defund 50 KYVE from a funder who has previously funded 100 KYVE", func() {
		// ACT
		s.RunTxPoolSuccess(&pooltypes.MsgDefundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  50 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(initialBalance - balanceAfter).To(Equal(50 * i.KYVE))

		Expect(pool.Funders).To(HaveLen(1))
		Expect(pool.TotalFunds).To(Equal(50 * i.KYVE))

		funder, funderFound := pool.GetFunder(i.ALICE)

		Expect(funderFound).To(BeTrue())
		Expect(funder).To(Equal(pooltypes.Funder{
			Address: i.ALICE,
			Amount:  50 * i.KYVE,
		}))
		Expect(pool.GetLowestFunder()).To(Equal(funder))
	})

	It("Try to defund more than actually funded", func() {
		// ACT
		s.RunTxPoolError(&pooltypes.MsgDefundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  101 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(initialBalance - balanceAfter).To(Equal(100 * i.KYVE))

		Expect(pool.Funders).To(HaveLen(1))
		Expect(pool.TotalFunds).To(Equal(100 * i.KYVE))

		funder, funderFound := pool.GetFunder(i.ALICE)

		Expect(funderFound).To(BeTrue())
		Expect(funder).To(Equal(pooltypes.Funder{
			Address: i.ALICE,
			Amount:  100 * i.KYVE,
		}))
		Expect(pool.GetLowestFunder()).To(Equal(funder))
	})

	It("Defund full funding amount from a funder who has previously funder 100 KYVE", func() {
		// ACT
		s.RunTxPoolSuccess(&pooltypes.MsgDefundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		// ASSERT
		balanceAfter := s.GetBalanceFromAddress(i.ALICE)

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		Expect(initialBalance - balanceAfter).To(BeZero())

		Expect(pool.Funders).To(BeEmpty())
		Expect(pool.TotalFunds).To(BeZero())

		funder, funderFound := pool.GetFunder(i.ALICE)

		Expect(funderFound).To(BeFalse())
		Expect(funder).To(Equal(pooltypes.Funder{}))
		Expect(pool.GetLowestFunder()).To(Equal(funder))
	})

	It("Defund as highest funder 75 KYVE in order to be the lowest funder afterwards", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  50 * i.KYVE,
		})

		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		funderBob, _ := pool.GetFunder(i.BOB)
		Expect(pool.GetLowestFunder()).To(Equal(funderBob))

		// ACT
		s.RunTxPoolSuccess(&pooltypes.MsgDefundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  75 * i.KYVE,
		})

		// ASSERT
		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)
		funderAlice, _ := pool.GetFunder(i.ALICE)
		Expect(pool.GetLowestFunder()).To(Equal(funderAlice))
	})
})
