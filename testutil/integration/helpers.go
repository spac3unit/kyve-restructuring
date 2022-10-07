package integration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) GetBalanceFromAddress(address string) uint64 {
	accAddress, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return 0
	}

	balance := suite.App().BankKeeper.GetBalance(suite.Ctx(), accAddress, "tkyve")

	return uint64(balance.Amount.Int64())
}

func (suite *KeeperTestSuite) GetBalanceFromModule(moduleName string) uint64 {
	moduleAcc := suite.App().AccountKeeper.GetModuleAccount(suite.Ctx(), moduleName).GetAddress()
	return suite.App().BankKeeper.GetBalance(suite.Ctx(), moduleAcc, "tkyve").Amount.Uint64()
}

func (suite *KeeperTestSuite) GetNextUploader() (nextStaker string, nextValaddress string) {
	bundleProposal, _ := suite.App().BundlesKeeper.GetBundleProposal(suite.Ctx(), 0)

	switch bundleProposal.NextUploader {
	case STAKER_0:
		nextStaker = STAKER_0
		nextValaddress = VALADDRESS_0
	case STAKER_1:
		nextStaker = STAKER_1
		nextValaddress = VALADDRESS_1
	case STAKER_2:
		nextStaker = STAKER_2
		nextValaddress = VALADDRESS_2
	default:
		nextStaker = ""
		nextValaddress = ""
	}

	return
}
