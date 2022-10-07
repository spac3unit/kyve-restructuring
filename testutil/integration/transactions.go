package integration

import (
	. "github.com/onsi/gomega"

	"github.com/KYVENetwork/chain/x/bundles"
	"github.com/KYVENetwork/chain/x/delegation"
	"github.com/KYVENetwork/chain/x/pool"
	"github.com/KYVENetwork/chain/x/stakers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
)

func (suite *KeeperTestSuite) RunTxGov(msg sdk.Msg) (*sdk.Result, error) {
	cachedCtx, commit := suite.ctx.CacheContext()
	resp, err := gov.NewHandler(suite.app.GovKeeper)(cachedCtx, msg)
	if err == nil {
		commit()
		return resp, nil
	}
	return nil, err
}

func (suite *KeeperTestSuite) RunTxPool(msg sdk.Msg) (*sdk.Result, error) {
	cachedCtx, commit := suite.ctx.CacheContext()
	resp, err := pool.NewHandler(suite.app.PoolKeeper)(cachedCtx, msg)
	if err == nil {
		commit()
		return resp, nil
	}
	return nil, err
}

func (suite *KeeperTestSuite) RunTxStakers(msg sdk.Msg) (*sdk.Result, error) {
	cachedCtx, commit := suite.ctx.CacheContext()
	resp, err := stakers.NewHandler(suite.app.StakersKeeper)(cachedCtx, msg)
	if err == nil {
		commit()
		return resp, nil
	}
	return nil, err
}

func (suite *KeeperTestSuite) RunTxDelegator(msg sdk.Msg) (*sdk.Result, error) {
	cachedCtx, commit := suite.ctx.CacheContext()
	resp, err := delegation.NewHandler(suite.app.DelegationKeeper)(cachedCtx, msg)
	if err == nil {
		commit()
		return resp, nil
	}
	return nil, err
}

func (suite *KeeperTestSuite) RunTxBundles(msg sdk.Msg) (*sdk.Result, error) {
	cachedCtx, commit := suite.ctx.CacheContext()
	resp, err := bundles.NewHandler(suite.app.BundlesKeeper)(cachedCtx, msg)
	if err == nil {
		commit()
		return resp, nil
	}
	return nil, err
}

func (suite *KeeperTestSuite) RunTxGovSuccess(msg sdk.Msg) {
	_, err := suite.RunTxGov(msg)
	Expect(err).To(BeNil())
}

func (suite *KeeperTestSuite) RunTxGovError(msg sdk.Msg) {
	_, err := suite.RunTxGov(msg)
	Expect(err).NotTo(BeNil())
}

func (suite *KeeperTestSuite) RunTxPoolSuccess(msg sdk.Msg) {
	_, err := suite.RunTxPool(msg)
	Expect(err).To(BeNil())
}

func (suite *KeeperTestSuite) RunTxPoolError(msg sdk.Msg) {
	_, err := suite.RunTxPool(msg)
	Expect(err).NotTo(BeNil())
}

func (suite *KeeperTestSuite) RunTxStakersSuccess(msg sdk.Msg) {
	_, err := suite.RunTxStakers(msg)
	Expect(err).To(BeNil())
}

func (suite *KeeperTestSuite) RunTxStakersError(msg sdk.Msg) {
	_, err := suite.RunTxStakers(msg)
	Expect(err).NotTo(BeNil())
}

func (suite *KeeperTestSuite) RunTxDelegatorSuccess(msg sdk.Msg) {
	_, err := suite.RunTxDelegator(msg)
	Expect(err).To(BeNil())
}

func (suite *KeeperTestSuite) RunTxDelegatorError(msg sdk.Msg) {
	_, err := suite.RunTxDelegator(msg)
	Expect(err).NotTo(BeNil())
}

func (suite *KeeperTestSuite) RunTxBundlesSuccess(msg sdk.Msg) {
	_, err := suite.RunTxBundles(msg)
	Expect(err).To(BeNil())
}

func (suite *KeeperTestSuite) RunTxBundlesError(msg sdk.Msg) {
	_, err := suite.RunTxBundles(msg)
	Expect(err).NotTo(BeNil())
}
