package keeper_test

import (
	"errors"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	deltamodulekeeper "github.com/delta-chain/delta/v2/x/delta/keeper"
	keepertest "github.com/delta-chain/delta/v2/x/delta/keeper/mock"
	"github.com/delta-chain/delta/v2/x/delta/types"
)

func (suite *KeeperTestSuite) TestGetSourceChannelID() {
	testCases := []struct {
		name          string
		ibcDenom      string
		expectedError error
		postCheck     func(channelID string)
	}{
		{
			"wrong ibc denom",
			"test",
			errors.New("test is invalid: ibc cro denom is invalid"),
			func(channelID string) {},
		},
		{
			"correct ibc denom",
			types.IbcCroDenomDefaultValue,
			nil,
			func(channelID string) {
				suite.Require().Equal(channelID, "channel-0")
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset
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

			channelID, err := suite.app.DeltaKeeper.GetSourceChannelID(suite.ctx, tc.ibcDenom)
			if tc.expectedError != nil {
				suite.Require().EqualError(err, tc.expectedError.Error())
			} else {
				suite.Require().NoError(err)
				tc.postCheck(channelID)
			}
		})
	}
}
