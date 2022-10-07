package keeper_test

import (
	i "github.com/KYVENetwork/chain/testutil/integration"
	bundletypes "github.com/KYVENetwork/chain/x/bundles/types"
	pooltypes "github.com/KYVENetwork/chain/x/pool/types"
	stakertypes "github.com/KYVENetwork/chain/x/stakers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*

TEST CASES - funding bundles

* Produce a valid bundle with only one funder
* Produce a valid bundle with multiple funders and same funding amounts
* Produce a valid bundle with multiple funders and different funding amounts

*/

var _ = Describe("funding bundles", Ordered, func() {
	s := i.NewCleanChain()

	initialBalanceAlice := s.GetBalanceFromAddress(i.ALICE)
	initialBalanceBob := s.GetBalanceFromAddress(i.BOB)

	BeforeEach(func() {
		// init new clean chain
		s = i.NewCleanChain()

		// create clean pool for every test case
		s.App().PoolKeeper.AppendPool(s.Ctx(), pooltypes.Pool{
			Name:           "Moontest",
			MaxBundleSize:  100,
			StartKey:       "0",
			UploadInterval: 60,
			OperatingCost:  10_000,
			Protocol: &pooltypes.Protocol{
				Version:     "0.0.0",
				Binaries:    "{}",
				LastUpgrade: uint64(s.Ctx().BlockTime().Unix()),
			},
			UpgradePlan: &pooltypes.UpgradePlan{},
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

	It("Produce a valid bundle with only one funder", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
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
			StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   100,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "test_hash",
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.VALADDRESS_0,
			Staker:     i.STAKER_0,
			PoolId:     0,
			StorageId:  "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			ByteSize:   100,
			FromHeight: 100,
			ToHeight:   200,
			FromKey:    "99",
			ToKey:      "199",
			ToValue:    "test_value2",
			BundleHash: "test_hash2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		totalReward := 100*s.App().BundlesKeeper.StorageCost(s.Ctx()) + pool.OperatingCost

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert total pool funds
		Expect(pool.TotalFunds).To(Equal(100*i.KYVE - totalReward))
		Expect(pool.Funders).To(HaveLen(1))

		// assert individual funds
		funder, _ := pool.GetFunder(i.ALICE)
		Expect(funder.Amount).To(Equal(100*i.KYVE - totalReward))

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)

		Expect(balanceAlice).To(Equal(initialBalanceAlice - 100*i.KYVE))
	})

	It("Produce a valid bundle with multiple funders and same funding amounts", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  100 * i.KYVE,
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  100 * i.KYVE,
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
			StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   100,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "test_hash",
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.VALADDRESS_0,
			Staker:     i.STAKER_0,
			PoolId:     0,
			StorageId:  "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			ByteSize:   100,
			FromHeight: 100,
			ToHeight:   200,
			FromKey:    "99",
			ToKey:      "199",
			ToValue:    "test_value2",
			BundleHash: "test_hash2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		totalReward := 100*s.App().BundlesKeeper.StorageCost(s.Ctx()) + pool.OperatingCost

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert total pool funds
		Expect(pool.TotalFunds).To(Equal(200*i.KYVE - totalReward))
		Expect(pool.Funders).To(HaveLen(2))

		// assert individual funds
		fundersCharge := uint64(sdk.NewDec(int64(totalReward)).Quo(sdk.NewDec(2)).TruncateInt64())

		funderAlice, _ := pool.GetFunder(i.ALICE)
		Expect(funderAlice.Amount).To(Equal(100*i.KYVE - fundersCharge))

		funderBob, _ := pool.GetFunder(i.BOB)
		Expect(funderBob.Amount).To(Equal(100*i.KYVE - fundersCharge))

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)
		Expect(balanceAlice).To(Equal(initialBalanceAlice - 100*i.KYVE))

		balanceBob := s.GetBalanceFromAddress(i.BOB)
		Expect(balanceBob).To(Equal(initialBalanceBob - 100*i.KYVE))
	})

	It("Produce a valid bundle with multiple funders and different funding amounts", func() {
		// ARRANGE
		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.ALICE,
			Id:      0,
			Amount:  150 * i.KYVE,
		})

		s.RunTxPoolSuccess(&pooltypes.MsgFundPool{
			Creator: i.BOB,
			Id:      0,
			Amount:  50 * i.KYVE,
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
			StorageId:  "y62A3tfbSNcNYDGoL-eXwzyV-Zc9Q0OVtDvR1biJmNI",
			ByteSize:   100,
			FromHeight: 0,
			ToHeight:   100,
			FromKey:    "0",
			ToKey:      "99",
			ToValue:    "test_value",
			BundleHash: "test_hash",
		})

		s.CommitAfterSeconds(60)

		// ACT
		s.RunTxBundlesSuccess(&bundletypes.MsgSubmitBundleProposal{
			Creator:    i.VALADDRESS_0,
			Staker:     i.STAKER_0,
			PoolId:     0,
			StorageId:  "P9edn0bjEfMU_lecFDIPLvGO2v2ltpFNUMWp5kgPddg",
			ByteSize:   100,
			FromHeight: 100,
			ToHeight:   200,
			FromKey:    "99",
			ToKey:      "199",
			ToValue:    "test_value2",
			BundleHash: "test_hash2",
		})

		// ASSERT
		pool, _ := s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		totalReward := 100*s.App().BundlesKeeper.StorageCost(s.Ctx()) + pool.OperatingCost

		pool, _ = s.App().PoolKeeper.GetPool(s.Ctx(), 0)

		// assert total pool funds
		Expect(pool.TotalFunds).To(Equal(200*i.KYVE - totalReward))
		Expect(pool.Funders).To(HaveLen(2))

		// assert individual funds
		fundersCharge := uint64(sdk.NewDec(int64(totalReward)).Quo(sdk.NewDec(2)).TruncateInt64())

		funderAlice, _ := pool.GetFunder(i.ALICE)
		Expect(funderAlice.Amount).To(Equal(150*i.KYVE - fundersCharge))

		funderBob, _ := pool.GetFunder(i.BOB)
		Expect(funderBob.Amount).To(Equal(50*i.KYVE - fundersCharge))

		// assert individual balances
		balanceAlice := s.GetBalanceFromAddress(i.ALICE)
		Expect(balanceAlice).To(Equal(initialBalanceAlice - 150*i.KYVE))

		balanceBob := s.GetBalanceFromAddress(i.BOB)
		Expect(balanceBob).To(Equal(initialBalanceBob - 50*i.KYVE))
	})
})
