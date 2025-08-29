package keeper_test

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	deltamodulekeeper "github.com/delta-chain/delta/v2/x/delta/keeper"
	handlers "github.com/delta-chain/delta/v2/x/delta/keeper/evmhandlers"
	keepertest "github.com/delta-chain/delta/v2/x/delta/keeper/mock"
	"github.com/delta-chain/delta/v2/x/delta/types"
)

func (suite *KeeperTestSuite) TestEvmHooks() {
	suite.SetupTest()

	contract := common.BigToAddress(big.NewInt(1))
	recipient := common.BigToAddress(big.NewInt(3))
	sender := common.BigToAddress(big.NewInt(4))

	testCases := []struct {
		msg      string
		malleate func()
	}{
		{
			"invalid log data, but still success",
			func() {
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{handlers.SendToAccountEvent.ID},
					},
				}
				receipt := &ethtypes.Receipt{
					Logs: logs,
				}
				err := suite.app.EvmKeeper.PostTxProcessing(suite.ctx, nil, receipt)
				suite.Require().NoError(err)
			},
		},
		{
			"not enough balance, expect fail",
			func() {
				data, err := handlers.SendToAccountEvent.Inputs.NonIndexed().Pack(
					recipient,
					big.NewInt(100),
				)
				suite.Require().NoError(err)
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{handlers.SendToAccountEvent.ID},
						Data:    data,
					},
				}
				receipt := &ethtypes.Receipt{
					Logs: logs,
				}
				err = suite.app.EvmKeeper.PostTxProcessing(suite.ctx, nil, receipt)
				suite.Require().Error(err)
			},
		},
		{
			"success send to account",
			func() {
				suite.app.DeltaKeeper.SetExternalContractForDenom(suite.ctx, denom, contract)
				coin := sdk.NewCoin(denom, sdkmath.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(contract.Bytes()), denom)
				suite.Require().Equal(coin, balance)

				data, err := handlers.SendToAccountEvent.Inputs.NonIndexed().Pack(
					recipient,
					coin.Amount.BigInt(),
				)
				suite.Require().NoError(err)
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{handlers.SendToAccountEvent.ID},
						Data:    data,
					},
				}
				receipt := &ethtypes.Receipt{
					Logs: logs,
				}
				err = suite.app.EvmKeeper.PostTxProcessing(suite.ctx, nil, receipt)
				suite.Require().NoError(err)

				balance = suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(contract.Bytes()), denom)
				suite.Require().Equal(sdk.NewCoin(denom, sdkmath.NewInt(0)), balance)
				balance = suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(recipient.Bytes()), denom)
				suite.Require().Equal(coin, balance)
			},
		},
		{
			"failed send to ibc, invalid ibc denom",
			func() {
				suite.SetupTest()
				// Create Delta Keeper with mock transfer keeper
				deltaKeeper := *deltamodulekeeper.NewKeeper(
					suite.app.EncodingConfig().Codec,
					suite.app.GetKey(types.StoreKey),
					suite.app.GetKey(types.MemStoreKey),
					suite.app.BankKeeper,
					keepertest.IbcKeeperMock{},
					suite.app.EvmKeeper,
					suite.app.AccountKeeper,
					authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				)
				suite.app.DeltaKeeper = deltaKeeper

				suite.app.DeltaKeeper.SetExternalContractForDenom(suite.ctx, denom, contract)
				coin := sdk.NewCoin(denom, sdkmath.NewInt(100))
				err := suite.MintCoins(sdk.AccAddress(contract.Bytes()), sdk.NewCoins(coin))
				suite.Require().NoError(err)

				balance := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.AccAddress(contract.Bytes()), denom)
				suite.Require().Equal(coin, balance)

				data, err := handlers.SendToIbcEvent.Inputs.NonIndexed().Pack(
					sender,
					"recipient",
					coin.Amount.BigInt(),
				)
				suite.Require().NoError(err)
				logs := []*ethtypes.Log{
					{
						Address: contract,
						Topics:  []common.Hash{handlers.SendToIbcEvent.ID},
						Data:    data,
					},
				}
				receipt := &ethtypes.Receipt{
					Logs: logs,
				}
				err = suite.app.EvmKeeper.PostTxProcessing(suite.ctx, nil, receipt)
				// should fail, because of not ibc denom name
				suite.Require().Error(err)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
		})
	}
}
